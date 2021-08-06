package pgmigrate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// const (
// 	infoColor    = "\033[0;36m%s\033[0m\n"
// 	warningColor = "\033[1;33m%s\033[0m\n"
// 	errorColor   = "\033[1;31m%s\033[0m\n"
// )

const (
	checkSchemaExistStmt = `SELECT EXISTS (
		SELECT schema_name FROM information_schema.schemata
		WHERE schema_name = (SELECT current_schema()));`

	currentSchemaStmt = `SELECT current_schema();`

	checkMigrateTableExistStmt = `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = (SELECT current_schema())
		AND    table_name   = 'pg_migrations');`

	createMigrateTableStmt = `CREATE TABLE IF NOT EXISTS pg_migrations (
			"version"    integer NOT NULL,
			"dirty"      boolean
		);
		INSERT INTO pg_migrations (version, dirty) VALUES (0, false);`

	updateMigrateTableStmt = `UPDATE pg_migrations SET version = $1, dirty = $2;`

	currentVersionStmt = `SELECT * FROM pg_migrations LIMIT 1;`
)

// DBWorker database interface
type DBWorker interface {
	CurrentSchema() string
	CheckSchemaExist() (bool, error)
	CheckMigrateTableExist() (bool, error)
	CurrentVersion() (int, bool, error)
	CreateMigrateTable() error
	UpdateMigrateTable(int, bool) error
	ExecMigration(string) error
}

// Migrate struct
type Migrate struct {
	Path              string
	DB                DBWorker
	step              int
	skip              []int
	gotov             int  // goto version
	version           int  // current version
	dirty             bool // dirty version
	migrateTableExist bool
}

// Files for migration
type Files struct {
	Version  int
	FileName string
}

// Step migrations
func (m *Migrate) Step(step int) *Migrate {
	m.step = step
	return m
}

// Up migrations
func (m *Migrate) Up() error {
	outputStarted()
	defer outputCompleted()
	fmt.Printf("Select current schema: %s\n", m.DB.CurrentSchema())
	return m.prepare().runUp()
}

// Down migrations
func (m *Migrate) Down() error {
	outputStarted()
	defer outputCompleted()
	fmt.Printf("Select current schema: %s\n", m.DB.CurrentSchema())
	return m.prepare().runDown()
}

// Goto migrate to version
func (m *Migrate) Goto(version int) error {
	outputStarted()
	defer outputCompleted()
	fmt.Printf("Select current schema: %s\n", m.DB.CurrentSchema())

	m.gotov = version
	tableExist, err := m.DB.CheckMigrateTableExist()
	if err != nil {
		return err
	}
	if tableExist {
		m.version, _, err = m.DB.CurrentVersion()
		if err != nil {
			return err
		}
	}

	if version == m.version {
		return err
	}

	if version > m.version {
		return m.prepare().runUp()
	}

	return m.prepare().runDown()
}

// Skip version (step)
func (m *Migrate) Skip(steps []int) *Migrate {
	m.skip = steps
	return m
}

// Version get current version
func (m *Migrate) Version() int {
	m.prepare()
	return m.version
}

func (m *Migrate) runUp() error {
	files, countFiles, err := m.getFilesUp()
	if err != nil {
		return err
	}
	if countFiles == 0 {
		fmt.Printf("warning: %s\n", "No change files")
		return nil
	}
	// maximum number of versions
	maxStep := maxStep(countFiles, m.step)
	// sort ascending
	sort.Slice(files[:], func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	for _, file := range files[0:maxStep] {
		if skipStep(file.Version, m.skip) {
			fmt.Printf("notice: %s marked as skipped\n", file.FileName)
			continue
		}
		err := m.migrateFromFile(m.Path + "/" + file.FileName)
		if err != nil {
			fmt.Printf("error: %s, %s\n", file.FileName, err)
			m.dirty = true
			break
		}
		m.version = file.Version
	}
	return m.complete()
}

func (m *Migrate) runDown() error {
	files, countFiles, err := m.getFilesDown()
	if err != nil {
		return err
	}
	if countFiles == 0 {
		fmt.Printf("warning: %s\n", "No change files")
		return nil
	}
	// maximum number of versions
	maxStep := maxStep(countFiles, m.step)
	// descending sort
	sort.Slice(files[:], func(i, j int) bool {
		return files[i].Version > files[j].Version
	})

	for _, file := range files[0:maxStep] {
		if skipStep(file.Version, m.skip) {
			fmt.Printf("notice: %s marked as skipped\n", file.FileName)
			continue
		}
		err := m.migrateFromFile(m.Path + "/" + file.FileName)
		if err != nil {
			fmt.Printf("error: %s, %s\n", file.FileName, err)
			m.dirty = true
			break
		}
		m.version = file.Version - 1
	}
	return m.complete()
}

func skipStep(step int, skip []int) bool {
	for _, v := range skip {
		if step == v {
			return true
		}
	}
	return false
}

func (m *Migrate) prepare() *Migrate {
	// checking for the existence of a schema and service table with migrations
	var err error
	m.migrateTableExist, err = m.DB.CheckMigrateTableExist()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return nil
	}
	// if the migrations table exists, get the current version of migrations
	if m.migrateTableExist {
		m.version, _, err = m.DB.CurrentVersion()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return nil
		}
	}
	return m
}

func (m *Migrate) complete() error {
	defer fmt.Printf("Now version: %d, dirty: %t", m.version, m.dirty)
	_, err := m.DB.CheckSchemaExist()
	if err != nil {
		return err
	}
	if !m.migrateTableExist {
		err := m.DB.CreateMigrateTable()
		if err != nil {
			return err
		}
	}
	err = m.DB.UpdateMigrateTable(m.version, m.dirty)
	if err != nil {
		return err
	}
	return nil
}

// Get the maximum number of steps from the current migration
func maxStep(countFiles int, step int) int {
	if step != 0 && step < countFiles {
		return step
	}
	return countFiles
}

// Retrieving file names from a directory with migrations
func (m *Migrate) getFilesUp() ([]Files, int, error) {
	files, err := ioutil.ReadDir(m.Path)
	if err != nil {
		return nil, 0, err
	}
	var migFiles []Files
	for _, f := range files {
		if strings.Contains(f.Name(), ".up.sql") {
			fileVersion, err := fileVersion(f.Name())
			if err != nil {
				fmt.Printf("notice: %s (skipped)\n", err.Error())
				continue
			}
			if fileVersion > m.version {
				// skip if 'goto version' is set
				if m.gotov != 0 && fileVersion > m.gotov {
					continue
				}
				migFiles = append(migFiles, Files{
					FileName: f.Name(),
					Version:  fileVersion,
				})
			}
		}
	}
	return migFiles, len(migFiles), nil
}

// Retrieving file names from a directory with migrations
func (m *Migrate) getFilesDown() ([]Files, int, error) {
	files, err := ioutil.ReadDir(m.Path)
	if err != nil {
		return nil, 0, err
	}
	var migFiles []Files
	for _, f := range files {
		if strings.Contains(f.Name(), ".down.sql") {
			fileVersion, err := fileVersion(f.Name())
			if err != nil {
				fmt.Printf("notice: %s (skipped)\n", err.Error())
				continue
			}
			if fileVersion <= m.version {
				// skip if 'goto version' is set
				if m.gotov != 0 && fileVersion <= m.gotov {
					continue
				}
				migFiles = append(migFiles, Files{
					FileName: f.Name(),
					Version:  fileVersion,
				})
			}
		}
	}
	return migFiles, len(migFiles), nil
}

// Get the version of a file from the name of a file
func fileVersion(fileName string) (int, error) {
	s := strings.Split(fileName, "_")
	if s[0] != "" {
		v, err := strconv.Atoi(s[0])
		if err != nil {
			return 0, errors.New("incorrect file name, it is not possible to read the version: " + err.Error())
		}
		return v, nil
	}
	return 0, nil
}

// Read the contents of the file and perform the migration
func (m *Migrate) migrateFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: file %s read error: %s (skipped)\n", err, filePath)
		return nil
	}
	if string(b) == "" {
		fmt.Printf("warning: file is empty: %s (skipped)\n", filePath)
		return nil
	}
	// execute a query on the server
	err = m.DB.ExecMigration(string(b))
	if err != nil {
		return err
	}
	fmt.Printf("Done: %s\n", filePath)
	return nil
}

func outputStarted() {
	fmt.Printf("\n%s\n\n", "****** Migration started *****")
}
func outputCompleted() {
	fmt.Printf("\n%s\n", "****** Migration completed *****")
}

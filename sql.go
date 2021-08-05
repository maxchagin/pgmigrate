package pgmigrate

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Sql structure for sql
type Sql struct {
	DB *sql.DB
}

// Config DB connection
type Config struct {
	DBname        string
	Host          string
	User          string
	Password      string
	Port          string
	SSLMode       string
	RuntimeParams map[string]string // ex: search_path, application_name
}

func NewWithConfig(sourcePath string, config *Config) (*Migrate, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.DBname, config.Password, config.SSLMode)
	for i, v := range config.RuntimeParams {
		if v != "" {
			connStr += fmt.Sprintf(" %s=%s", i, v)
		}
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return CompatibleWithSql(
		sourcePath,
		&Sql{
			DB: db,
		}), nil
}

func New(sourcePath string, host, port, user, dbname, password, sslmode string, runtimeParams ...map[string]string) (*Migrate, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	for _, params := range runtimeParams {
		for i, v := range params {
			if v != "" {
				connStr += fmt.Sprintf(" %s=%s", i, v)
			}
		}
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return CompatibleWithSql(
		sourcePath,
		&Sql{
			DB: db,
		}), nil
}

// CompatibleWithSql sql compatible
func CompatibleWithSql(sourcePath string, sql *Sql) *Migrate {
	return &Migrate{
		Path: sourcePath,
		DB:   sql,
	}
}

// CurrentSchema get the current schema
// Before the creation of the schema, it may not exist, in this case the value undefined is returned
func (s *Sql) CurrentSchema() string {
	var currentSchema string
	err := s.DB.QueryRow(currentSchemaStmt).Scan(&currentSchema)
	if err != nil {
		return "undefined"
	}
	if currentSchema == "" {
		return "undefined"
	}
	return currentSchema
}

// ExecMigration executing content from migration file
func (s *Sql) ExecMigration(content string) error {
	_, err := s.DB.Exec(content)
	if err != nil {
		return fmt.Errorf("%v: %v", errExecMigration, err)
	}
	return nil
}

// CheckSchemaExist checking for the existence of a schema
func (s *Sql) CheckSchemaExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(checkSchemaExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckSchemaExist, err)
	}
	return exists, nil
}

// CheckMigrateTableExist checking for the existence of the migration table
func (s *Sql) CheckMigrateTableExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(checkMigrateTableExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckMigrateTableExist, err)
	}
	return exists, nil
}

// CreateMigrateTable creating a migrations table with a zero version
func (s *Sql) CreateMigrateTable() error {
	_, err := s.DB.Exec(createMigrateTableStmt)
	if err != nil {
		return fmt.Errorf("%v: %v", errCreateMigrateTable, err)
	}
	return nil
}

// UpdateMigrateTable updating the pg migrations table
func (s *Sql) UpdateMigrateTable(version int) error {
	var dirty bool
	err := s.DB.QueryRow(updateMigrateTableStmt, version).Scan(&dirty)
	if err != nil {
		return fmt.Errorf("%v: %v", errUpdateMigrateTable, err)
	}
	return nil
}

// CurrentVersion getting the current version of the migration
func (s *Sql) CurrentVersion() (int, bool, error) {
	var version int
	var dirty bool
	err := s.DB.QueryRow(currentVersionStmt).Scan(&version, &dirty)
	if err != nil {
		return 0, false, fmt.Errorf("%v: %v", errCurrentVersion, err)
	}
	return version, dirty, nil
}

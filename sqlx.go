package pgmigrate

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Sqlx structure for sqlx
type Sqlx struct {
	DB *sqlx.DB
}

// CompatibleWithSqlx sqlx compatible
func CompatibleWithSqlx(sourcePath string, sqlx *Sqlx) *Migrate {
	return &Migrate{
		Path: sourcePath,
		DB:   sqlx,
	}
}

// CurrentSchema get the current schema
// Before the creation of the schema, it may not exist, in this case the value undefined is returned
func (s *Sqlx) CurrentSchema() string {
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
func (s *Sqlx) ExecMigration(content string) error {
	_, err := s.DB.Exec(content)
	if err != nil {
		return fmt.Errorf("%v: %v", errExecMigration, err)
	}
	return nil
}

// CheckSchemaExist checking for the existence of a schema
func (s *Sqlx) CheckSchemaExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(checkSchemaExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckSchemaExist, err)
	}
	return exists, nil
}

// CheckMigrateTableExist checking for the existence of the migration table
func (s *Sqlx) CheckMigrateTableExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(checkMigrateTableExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckMigrateTableExist, err)
	}
	return exists, nil
}

// CreateMigrateTable creating a migrations table with a zero version
func (s *Sqlx) CreateMigrateTable() error {
	_, err := s.DB.Exec(createMigrateTableStmt)
	if err != nil {
		return fmt.Errorf("%v: %v", errCreateMigrateTable, err)
	}
	return nil
}

// UpdateMigrateTable updating the pg migrations table
func (s *Sqlx) UpdateMigrateTable(version int, dirty bool) error {
	row := s.DB.QueryRowx(updateMigrateTableStmt, version, dirty)
	if row.Err() != nil {
		return fmt.Errorf("%v: %v", errUpdateMigrateTable, row.Err())
	}
	return nil
}

// CurrentVersion getting the current version of the migration
func (s *Sqlx) CurrentVersion() (int, bool, error) {
	var version int
	var dirty bool
	err := s.DB.QueryRow(currentVersionStmt).Scan(&version, &dirty)
	if err != nil {
		return 0, false, fmt.Errorf("%v: %v", errCurrentVersion, err)
	}
	return version, dirty, nil
}

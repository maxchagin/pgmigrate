package pgmigrate

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// Pgx structure for pgx
type Pgx struct {
	DB *pgx.Conn
}

// CompatibleWithPgx pgx compatible
func CompatibleWithPgx(sourcePath string, pgx *Pgx) *Migrate {
	return &Migrate{
		Path: sourcePath,
		DB:   pgx,
	}
}

// CurrentSchema get the current schema
// Before the creation of the schema, it may not exist, in this case the value undefined is returned
func (s *Pgx) CurrentSchema() string {
	var currentSchema string
	err := s.DB.QueryRow(context.Background(), currentSchemaStmt).Scan(&currentSchema)
	if err != nil {
		return "undefined"
	}
	if currentSchema == "" {
		return "undefined"
	}
	return currentSchema
}

// ExecMigration executing content from migration file
func (s *Pgx) ExecMigration(content string) error {
	_, err := s.DB.Exec(context.Background(), content)
	if err != nil {
		return fmt.Errorf("%v: %v", errExecMigration, err)
	}
	return nil
}

// CheckSchemaExist checking for the existence of a schema
func (s *Pgx) CheckSchemaExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(context.Background(), checkSchemaExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckSchemaExist, err)
	}
	return exists, nil
}

// CheckMigrateTableExist checking for the existence of the migration table
func (s *Pgx) CheckMigrateTableExist() (bool, error) {
	var exists bool
	err := s.DB.QueryRow(context.Background(), checkMigrateTableExistStmt).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%v: %v", errCheckMigrateTableExist, err)
	}
	return exists, nil
}

// CreateMigrateTable creating a migrations table with a zero version
func (s *Pgx) CreateMigrateTable() error {
	_, err := s.DB.Exec(context.Background(), createMigrateTableStmt)
	if err != nil {
		return fmt.Errorf("%v: %v", errCreateMigrateTable, err)
	}
	return nil
}

// UpdateMigrateTable updating the pg migrations table
func (s *Pgx) UpdateMigrateTable(version int, dirty bool) error {
	_, err := s.DB.Query(context.Background(), updateMigrateTableStmt, version, dirty)
	if err != nil {
		return fmt.Errorf("%v: %v", errUpdateMigrateTable, err)
	}
	return nil
}

// CurrentVersion getting the current version of the migration
func (s *Pgx) CurrentVersion() (int, bool, error) {
	var version int
	var dirty bool
	err := s.DB.QueryRow(context.Background(), currentVersionStmt).Scan(&version, &dirty)
	if err != nil {
		return 0, false, fmt.Errorf("%v: %v", errCurrentVersion, err)
	}
	return version, dirty, nil
}

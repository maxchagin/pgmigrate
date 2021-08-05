package pgmigrate

import "errors"

var (
	errCheckSchemaExist       = errors.New("failed to check exists schema")
	errExecMigration          = errors.New("failed to exec migration")
	errCheckMigrateTableExist = errors.New("failed to check exists table migrations")
	errCreateMigrateTable     = errors.New("failed to create migrations table")
	errUpdateMigrateTable     = errors.New("failed to update migrations table")
	errCurrentVersion         = errors.New("failed to select version from migrations")
)

package pgmigrate

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func OpenSqlxConn() (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s search_path=%s application_name=%s", "localhost", "5432", "root", "test", "root", "disable", "test", "test")

	return sqlx.Connect("postgres", connStr)
}

func TestSqlxUp(t *testing.T) {
	ConnSqlx, err := OpenSqlxConn()
	if err != nil {
		t.Error(err)
	}
	m := CompatibleWithSqlx(
		"./migrations",
		&Sqlx{
			DB: ConnSqlx,
		})
	// up to 4
	err = m.Up()
	if err != nil {
		t.Error(err)
	}
}

func TestSqlxDown(t *testing.T) {
	ConnSqlx, err := OpenSqlxConn()
	if err != nil {
		t.Error(err)
	}
	m := CompatibleWithSqlx(
		"./migrations",
		&Sqlx{
			DB: ConnSqlx,
		})
	// down to 0
	err = m.Down()
	if err != nil {
		t.Error(err)
	}
}

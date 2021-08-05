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

func TestSqlxUpAndDown(t *testing.T) {
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
	// go to 2
	// err = m.Goto(2)
	// if err != nil {
	// 	t.Error(err)
	// }
	// down to 0
	err = m.Down()
	if err != nil {
		t.Error(err)
	}
}

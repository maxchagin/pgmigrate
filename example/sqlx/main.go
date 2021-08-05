package main

import (
	"fmt"
	"log"

	"pgmigrate"

	"github.com/jmoiron/sqlx"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s search_path=%s application_name=%s", "localhost", "5432", "root", "test", "root", "disable", "test", "test")
	ConnSqlx, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	m := pgmigrate.CompatibleWithSqlx(
		"./migrations",
		&pgmigrate.Sqlx{
			DB: ConnSqlx,
		})
	// migrate up
	err = m.Up()
	if err != nil {
		log.Fatalln(err)
	}
}

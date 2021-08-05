package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maxchagin/pgmigrate"

	"github.com/jackc/pgx/v4"
)

func main() {
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "root", "root", "localhost", "5432", "test")
	config, err := pgx.ParseConfig(connURL)
	if err != nil {
		log.Fatalln(err)
	}
	config.RuntimeParams = map[string]string{
		"search_path":      "test",
		"application_name": "test",
	}

	connPgx, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalln(err)
	}
	m := pgmigrate.CompatibleWithPgx(
		"./migrations",
		&pgmigrate.Pgx{
			DB: connPgx,
		})
	// migrate up
	err = m.Up()
	if err != nil {
		log.Fatalln(err)
	}
}

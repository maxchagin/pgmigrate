package main

import (
	"log"

	"github.com/maxchagin/pgmigrate"
)

func main() {
	m, err := pgmigrate.New("./../../migrations", "localhost", "5432", "root", "test", "root", "disable", map[string]string{
		"search_path":      "test",
		"application_name": "test",
	})
	if err != nil {
		log.Fatalln(err)
	}
	// migrate up
	err = m.Up()
	if err != nil {
		log.Fatalln(err)
	}
}

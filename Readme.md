# Simple PostgreSQL migrations for Golang

This package allow run migrations on PostgreSQL database using Golang.   
The package is compatible with library [sqlx](https://github.com/jmoiron/sqlx) and driver [pgx](https://github.com/jackc/pgx)


## Installation
The package require a Go version with modules support
```
go mod init github.com/my/repo
go get github.com/maxchagin/pgmigrate
```

## Usage
To run the migration in your project, you need to follow the steps:
- define a list migration files
- init migration in your application
- run migration

## Migration files
The ordering and direction of the migration files is determined by the filenames used for them. The package expects the filenames of migrations to have the format:   
```
{version}_{title}.{action}.sql
```
Versioning with incrementing integers
```
1_create_table.up.sql - up migration
1_create_table.down.sql - down migration
```
Or timestamps resolution:
```
1627628025_create_table.up.sql - up migration
1627628025_create_table.down.sql - down migration
```
See more [example](https://github.com/maxchagin/pgmigrate/migrations)

## Run migrations
The following methods are supported:   
`Up()` - run all available migrations;   
`Down()` - down all migration;   
`Goto(version int)` - go to the specified migration;   
`Skip(steps []int)` - skip specified migrations;   
`Version()` - get the current version of the migration;

## Example Usage
### With [sql](https://pkg.go.dev/database/sql)
```go
package main

import (
	"log"

	"pgmigrate"
)

func main() {
	m, err := pgmigrate.New("./migrations", "localhost", "5432", "root", "test", "root", "disable", map[string]string{
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
```

### Compatible with [sqlx](https://github.com/jmoiron/sqlx)
```go
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

```

### Compatible with [pgx](https://github.com/jackc/pgx)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"pgmigrate"

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

```
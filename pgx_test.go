package pgmigrate

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jackc/pgx/v4"
)

// OpenPgxConn opening a connection with pgx
func OpenPgxConn() (*pgx.Conn, error) {
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "root", "root", "localhost", "5432", "test")
	config, err := pgx.ParseConfig(connURL)
	if err != nil {
		log.Printf("Unable to parse url: %v\n", err)
		return nil, err
	}
	config.RuntimeParams = map[string]string{
		"search_path":      "test",
		"application_name": "test",
	}

	db, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestPgxUp(t *testing.T) {
	ConnPgx, err := OpenPgxConn()
	if err != nil {
		t.Error(err)
	}
	m := CompatibleWithPgx(
		"./migrations",
		&Pgx{
			DB: ConnPgx,
		})
	// up to 4
	err = m.Up()
	if err != nil {
		t.Error(err)
	}
}

func TestPgxDown(t *testing.T) {
	ConnPgx, err := OpenPgxConn()
	if err != nil {
		t.Error(err)
	}
	m := CompatibleWithPgx(
		"./migrations",
		&Pgx{
			DB: ConnPgx,
		})
	// down to 0
	err = m.Down()
	if err != nil {
		t.Error(err)
	}
}

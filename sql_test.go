package pgmigrate

import (
	"testing"
)

func TestSqlUp(t *testing.T) {
	// m, err := New("./migrations", "localhost", "5432", "root", "test", "root", "disable")
	m, err := New("./migrations", "localhost", "5432", "root", "test", "root", "disable", map[string]string{
		"search_path":      "test",
		"application_name": "test",
	})

	if err != nil {
		t.Error(err)
	}
	// up to 0
	err = m.Up()
	if err != nil {
		t.Error(err)
	}
}
func TestSqlDown(t *testing.T) {
	// m, err := New("./migrations", "localhost", "5432", "root", "test", "root", "disable")
	m, err := New("./migrations", "localhost", "5432", "root", "test", "root", "disable", map[string]string{
		"search_path":      "test",
		"application_name": "test",
	})

	if err != nil {
		t.Error(err)
	}
	// down to 0
	err = m.Down()
	if err != nil {
		t.Error(err)
	}
}

func TestSqlUpAndDownWithConfig(t *testing.T) {
	m, err := NewWithConfig("./migrations", &Config{
		Host:     "localhost",
		Port:     "5432",
		DBname:   "test",
		User:     "root",
		Password: "root",
		SSLMode:  "disable",
		RuntimeParams: map[string]string{
			"search_path":      "test",
			"application_name": "test",
		},
	})
	if err != nil {
		t.Error(err)
	}
	// up to 4
	err = m.Up()
	if err != nil {
		t.Error(err)
	}
	// down to 0
	err = m.Down()
	if err != nil {
		t.Error(err)
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate"

	_ "github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to the database file")
	flag.StringVar(&migrationPath, "migrations-path", "", "path to the migration file")
	flag.StringVar(&migrationTable, "migrations-table", "", "name of the migration table")
	flag.Parse()

	if storagePath == "" || migrationPath == "" {
		panic("storage path or migration path cannot be empty")
	}

	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migration needed")
			return
		}
		panic(err)
	}
	fmt.Println("migration done")
}

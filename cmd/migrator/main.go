package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://migrations",
		"mysql://mysqladmin:mysqladmin@tcp(localhost:3306)/logging?multiStatements=true",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Применить все доступные миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Или для отката одной миграции
	// if err := m.Steps(-1); err != nil {
	//     log.Fatal(err)
	// }
}

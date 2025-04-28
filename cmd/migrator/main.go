package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func fetchMysqlConnStr() string {
	var res string

	flag.StringVar(&res, "mysql_conn_str", "", "")
	flag.Parse()

	if res == "" {
		res = os.Getenv("MYSQL_CONN_STR")
	}

	return res
}

func main() {
	dsn := "mysql://" + fetchMysqlConnStr() + "?multiStatements=true"
	log.Println(dsn)

	// Добавляем retry-логику
	var m *migrate.Migrate
	var err error

	for i := 0; i < 5; i++ {
		log.Printf("Attempt %d: %v", i+1, err)
		time.Sleep(5 * time.Second)
		m, err = migrate.New(
			"file://migrations",
			dsn,
		)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Fatal("Failed to connect after retries: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Migrations applied successfully!")
}

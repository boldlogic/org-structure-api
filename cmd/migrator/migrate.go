package main

import (
	"database/sql"
	"log"

	"github.com/boldlogic/org-structure-api/pkg/config"
	"github.com/boldlogic/packages/commonconfig"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	defaultConfigPath = "config.yaml"
)

func main() {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)

	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}
	dsn := conf.Database.GetDSN()
	if dsn == "" {
		log.Fatal("некорректный DSN")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

}

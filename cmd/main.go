package main

import (
	"database/sql"
	"github.com/AlecSmith96/dating-api/internal/adapters"
	"github.com/AlecSmith96/dating-api/internal/drivers"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

const (
	gooseDir = "./db/goose"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading in config", "err", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", conf.DatabaseConnectionString)
	if err != nil {
		slog.Error("connecting to postgres", "err", err)
		os.Exit(1)
	}

	postgresAdapter := adapters.NewPostgresAdapter(db, conf.JwtExpiryMillis, conf.JwtSecretKey)
	err = postgresAdapter.PerformDataMigration(gooseDir)
	if err != nil {
		slog.Error("performing data migration", "err", err)
		os.Exit(1)
	}

	router := drivers.NewRouter(postgresAdapter, postgresAdapter)

	router.Run(":8080")
}

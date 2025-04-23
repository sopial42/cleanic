package persistence

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/sopial42/cleanic/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func NewPGClient(cfg config.DBConfig) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.DBName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTimeout(5*time.Second)))

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	sqldb.SetConnMaxLifetime(30 * time.Minute)

	err := sqldb.Ping()
	if err != nil {
		panic(err)
	}

	client := bun.NewDB(sqldb, pgdialect.New())
	client.AddQueryHook(bundebug.NewQueryHook(
		// Ensure false by default
		bundebug.WithEnabled(false),
		bundebug.FromEnv("DB_LOG_LEVEL"),
	))

	return client
}

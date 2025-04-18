package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/kotai-tech/server/internal/config"
	domain "github.com/kotai-tech/server/internal/domain"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresPatientRepository struct {
	clientDB *bun.DB
}

func NewPGClient(cfg config.DBConfig) *PostgresPatientRepository {
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

	return &PostgresPatientRepository{clientDB: client}
}

func (r *PostgresPatientRepository) ListPatients(ctx context.Context) ([]domain.Patient, error) {
	var patientDAOs []patientDAO

	request := r.clientDB.NewSelect().Model(&patientDAOs)
	err := request.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	return ToDomainList(patientDAOs), nil
}

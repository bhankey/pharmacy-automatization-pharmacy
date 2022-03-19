package app

import (
	"fmt"

	"github.com/bhankey/pharmacy-automatization-pharmacy/internal/config"
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/postgresdb"
	"github.com/jmoiron/sqlx"
)

type dataSources struct {
	db *sqlx.DB
}

func newDataSource(config config.Config) (*dataSources, error) {
	postgresDB, err := postgresdb.NewClient(
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres connection error: %w", err)
	}

	return &dataSources{
		db: postgresDB,
	}, nil
}

package container

import (
	"github.com/bhankey/pharmacy-automatization-pharmacy/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	masterPostgresDB *sqlx.DB
	slavePostgresDB  *sqlx.DB
	logger           logger.Logger
	dependencies     map[string]interface{}
}

func NewContainer(
	log logger.Logger,
	masterPostgres, slavePostgres *sqlx.DB,
) *Container {
	return &Container{
		masterPostgresDB: masterPostgres,
		slavePostgresDB:  slavePostgres,
		logger:           log,
		dependencies:     make(map[string]interface{}),
	}
}

func (c *Container) CloseAllConnections() {
	if err := c.masterPostgresDB.Close(); err != nil {
		c.logger.Errorf("failed to close master postgres connection error: %v", err)
	}

	if err := c.slavePostgresDB.Close(); err != nil {
		c.logger.Errorf("failed to close slave postgres connection error: %v", err)
	}
}

package utils

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelvinator07/go-rest-template/internal/config"
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/kelvinator07/go-rest-template/internal/datasources/drivers"
)

func SetupPostgresConnection() (*sqlx.DB, error) {
	var dsn string
	switch config.AppConfig.Environment {
	case constants.EnvironmentDevelopment:
		dsn = config.AppConfig.DBPostgresDsn
	case constants.EnvironmentProduction:
		dsn = config.AppConfig.DBPostgresURL
	}

	// Setup sqlx config of postgreSQL
	config := drivers.SQLXConfig{
		DriverName:     config.AppConfig.DBPostgresDriver,
		DataSourceName: dsn,
		MaxOpenConns:   100,
		MaxIdleConns:   10,
		MaxLifetime:    15 * time.Minute,
	}

	// Initialize postgreSQL connection with sqlx
	conn, err := config.InitializeSQLXDatabase()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

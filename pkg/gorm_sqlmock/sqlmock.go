// Package gorm_sqlmock defines helpers for go orm mocking for test writing
// nolint
package gorm_sqlmock

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Config config
type Config struct {
	*gorm.Config
	DriverName              string
	DSN                     string
	DontSupportRenameColumn bool
	conn                    *sql.DB
	mock                    sqlmock.Sqlmock
}

// Open mock database with dsn
func Open(driverName string, dsn string) (db *gorm.DB, mock sqlmock.Sqlmock, err error) {
	config := &Config{
		DriverName: driverName,
	}

	config.conn, config.mock, err = sqlmock.NewWithDSN(dsn)
	if err != nil {
		return
	}

	return newMock(config)
}

// New mock database with config
func New(config Config) (db *gorm.DB, mock sqlmock.Sqlmock, err error) {

	config.conn, config.mock, err = sqlmock.New()
	if err != nil {
		return
	}

	return newMock(&config)
}

func newMock(config *Config) (db *gorm.DB, mock sqlmock.Sqlmock, err error) {

	var dialect gorm.Dialector

	switch config.DriverName {
	case "mysql":
	case "sqlserver", "pg":
		dialect = sqlserver.New(sqlserver.Config{
			DriverName: "sqlserver",
			DSN:        config.DSN,
			Conn:       config.conn,
		})
	case "postgres":
		// TODO
	case "sqlite":
		// TODO
	default:
		err = fmt.Errorf("the %s driver could not be matched", config.DriverName)
		return
	}

	db, err = gorm.Open(dialect, config.Config)

	return db, config.mock, err
}

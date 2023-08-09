// Package storage house all storage/database implementations that performs CRUD operations on our schema
package storage

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"employee-management-system/model/pagination"
	"employee-management-system/pkg/environment"
	"employee-management-system/pkg/gorm_sqlmock"
	"employee-management-system/pkg/helper"
)

const packageName = "storage"

// Storage object
type Storage struct {
	Logger zerolog.Logger
	Env    *environment.Env
	DB     *gorm.DB
}

// New Storage, however should panic if it can't be pinged. System should be able to connect to the database
func New(z zerolog.Logger, env *environment.Env) *Storage {
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()
	db, err := gorm.Open(
		sqlserver.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=%s",
			env.Get("PG_ADDRESS"),
			env.Get("PG_PORT"),
			env.Get("PG_USER"),
			env.Get("PG_DATABASE"),
			env.Get("PG_PASSWORD"),
			env.Get("TIMEZONE"),
		)),
		&gorm.Config{},
	)
	if err != nil {
		l.Fatal().Err(err)
		panic(err)
	}

	return &Storage{
		Logger: l,
		Env:    env,
		DB:     db,
	}
}

// GetStorage helper for tests/mock
// I expect our storage tests to use this helper going forward.
func GetStorage(t *testing.T) (sqlmock.Sqlmock, *Storage) {
	var (
		mock sqlmock.Sqlmock
		db   *gorm.DB
		err  error
	)

	db, mock, err = gorm_sqlmock.New(gorm_sqlmock.Config{
		Config:     &gorm.Config{},
		DriverName: "sqlserver",
		DSN:        "mock",
	})

	require.NoError(t, err)

	return mock, NewFromDB(db)
}

// NewFromDB created a new storage with just the database reference passed in
func NewFromDB(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

// Close securely closes the connection to the storage/database
func (d *Storage) Close() {
	sqlDD, _ := d.DB.DB()
	_ = sqlDD.Close()
}

func getPaging(page pagination.Page) pagination.Page {
	if page.Number == nil {
		tmpPageNumber := pagination.PageDefaultNumber
		page.Number = &tmpPageNumber
	}
	if page.Size == nil {
		tmpPageSize := pagination.PageDefaultSize
		page.Size = &tmpPageSize
	}
	if page.SortBy == nil {
		tmpPageSortBy := pagination.PageDefaultSortBy
		page.SortBy = &tmpPageSortBy
	}
	if page.SortDirectionDesc == nil {
		tmpPageSortDirectionDesc := pagination.PageDefaultSortDirectionDesc
		page.SortDirectionDesc = &tmpPageSortDirectionDesc
	}
	return page
}

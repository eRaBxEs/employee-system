package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	DB               *gorm.DB
	mock             sqlmock.Sqlmock
	employeeDatabase EmployeeDatabase
}

func (s *Suite) SetupSuite() {
	var store *Storage
	s.mock, store = GetStorage(s.T())

	s.employeeDatabase = *NewEmployee(store)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_Storage_Dummy() {
	// Ha ha, dont mind me with this test. SMH
	require.Equal(s.T(), true, true)
}

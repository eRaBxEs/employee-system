package storage

import (
	"context"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"employee-management-system/model"
)

var employeeTableColumns = []string{"id", "first_name", "last_name", "email",
	"dob", "department_id", "position", "updated_at"}

func (s *Suite) Test_GetEmployeeByID() {
	validUserID := 7
	id := 7
	departmentId := 27
	firstName := "James"
	lastName := "Blake"
	email := "jamie@yahoo.com"
	position := "lawyer"

	testEmployee := model.Employee{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Dob:          time.Now(),
		DepartmentID: departmentId,
		Position:     position,
		UpdatedAt:    time.Now(),
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "employee" WHERE id = $1`)).
		WithArgs(validUserID).
		WillReturnRows(sqlmock.NewRows(employeeTableColumns).
			AddRow(testEmployee.ID, testEmployee.FirstName, testEmployee.LastName, testEmployee.Email,
				testEmployee.Dob, testEmployee.DepartmentID, testEmployee.Position))

	retEmployee, err := s.employeeDatabase.GetEmployeeByID(context.Background(), validUserID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), retEmployee, testEmployee)
}

func (s *Suite) Test_AddEmployee() {
	id := 2
	departmentID := 30
	firstName := "Brown"
	lastName := "Lucid"
	email := "brown@yahoo.com"
	position := "recruiter"

	testEmployee := model.Employee{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Dob:          time.Now(),
		DepartmentID: departmentID,
		Position:     position,
		UpdatedAt:    time.Now(),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "employees" ("first_name","last_name","email","dob", "department_id", "position", "updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,) RETURNING "id"`)).
		WithArgs(testEmployee.FirstName, testEmployee.LastName, testEmployee.Email, testEmployee.Dob, testEmployee.Position).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(testEmployee.ID),
		)
	s.mock.ExpectCommit()

	newRecord, err := s.employeeDatabase.AddEmployee(context.Background(), testEmployee)

	require.NoError(s.T(), err)
	require.Equal(s.T(), newRecord.FirstName, testEmployee.FirstName)
}

func (s *Suite) Test_GetEmployeeByContext() {

}

func (s *Suite) Test_GetAllEmployees() {
	id := 5
	departmentID := 30
	firstName := "Brown"
	lastName := "Lucid"
	email := "brown@yahoo.com"
	position := "recruiter"
	testEmployee := model.Employee{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Dob:          time.Now(),
		DepartmentID: departmentID,
		Position:     position,
		UpdatedAt:    time.Now(),
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "employees"`)).
		WillReturnRows(sqlmock.NewRows(employeeTableColumns).
			AddRow(testEmployee.ID, testEmployee.FirstName, testEmployee.LastName, testEmployee.Email,
				testEmployee.Dob, testEmployee.DepartmentID, testEmployee.Position))

	retEmployees, err := s.employeeDatabase.GetAllEmployees(context.Background())

	require.NoError(s.T(), err)
	require.Equal(s.T(), retEmployees[0].ID, testEmployee.ID)
}

func (s *Suite) Test_UpdateEmployeeByID() {
	id := 5
	departmentID := 30
	firstName := "Brown"
	lastName := "Lucid"
	email := "brown@yahoo.com"
	position := "recruiter"
	testEmployee := model.Employee{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Dob:          time.Now(),
		DepartmentID: departmentID,
		Position:     position,
		UpdatedAt:    time.Now(),
	}

	updateUpdatedAt := time.Now()
	sqlmock.NewRows(employeeTableColumns).
		AddRow(testEmployee.ID, testEmployee.FirstName, testEmployee.LastName, testEmployee.Email,
			testEmployee.Dob, testEmployee.DepartmentID, testEmployee.Position)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "employees" SET "first_name"=$1,"last_name"=$2,"email"=$3,"dob"=$4, "department_id"=$5, "position"=$6 WHERE "id" = $7`)).
		WithArgs(testEmployee.FirstName, testEmployee.LastName, testEmployee.Email,
			testEmployee.Dob, testEmployee.DepartmentID, testEmployee.Position, testEmployee.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	retEmployee, err := s.employeeDatabase.UpdateEmployeeByID(context.Background(), testEmployee.ID, model.Employee{
		FirstName:    testEmployee.FirstName,
		LastName:     testEmployee.LastName,
		Email:        testEmployee.Email,
		DepartmentID: testEmployee.DepartmentID,
		Position:     testEmployee.Position,
		UpdatedAt:    updateUpdatedAt,
	})

	require.NoError(s.T(), err)
	require.Equal(s.T(), retEmployee.FirstName, testEmployee.FirstName)
}

func (s *Suite) Test_DeleteEmployeeByID() {
	validID := 6
	deletedAt := time.Now()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "employees" SET "deleted_at"=$1 WHERE "id" = $2`)).
		WithArgs(deletedAt, validID).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.employeeDatabase.DeleteEmployeeByID(context.Background(), validID)
	require.NoError(s.T(), err)
}

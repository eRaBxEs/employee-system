package storage

import (
	"context"

	"github.com/rs/zerolog"

	"employee-management-system/model"
	"employee-management-system/pkg/helper"
)

// EmployeeDatabase enlist all possible storage operations for Employee entity for User
//
//go:generate mockgen -source employee.go -destination ./mock/mock_employee.go -package mock EmployeeDatabase
type EmployeeDatabase interface {
	AddEmployee(ctx context.Context, employee model.Employee) (model.Employee, error)
	GetEmployeeByID(ctx context.Context, ID int) (model.Employee, error)
	GetEmployeeByContext(ctx context.Context, userID int) (model.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*model.Employee, error)
	UpdateEmployeeByID(ctx context.Context, id int, employee model.Employee) (model.Employee, error)
	DeleteEmployeeByID(ctx context.Context, id int) error
}

// Employee object
type Employee struct {
	logger  zerolog.Logger
	storage *Storage
}

// NewEmployee creates a new reference to the Employee storage entity
func NewEmployee(s *Storage) *EmployeeDatabase {
	l := s.Logger.With().Str(helper.LogStrKeyLevel, "employee").Logger()
	employee := &Employee{
		logger:  l,
		storage: s,
	}
	employeeDatabase := EmployeeDatabase(employee)
	return &employeeDatabase
}

// AddEmployee adds a new row into the employee table referencing users by user_id column
func (e *Employee) AddEmployee(ctx context.Context, employee model.Employee) (model.Employee, error) {
	db := e.storage.DB.WithContext(ctx).Create(&employee)
	if db.Error != nil {
		e.logger.Err(db.Error).Msgf("Employee::AddEmployee error: %v, (%v)", ErrRecordCreatingFailed, db.Error)
		return model.Employee{}, ErrRecordCreatingFailed
	}
	return employee, nil
}

// GetEmployeeByID retrieves a single row
func (e *Employee) GetEmployeeByID(ctx context.Context, ID int) (model.Employee, error) {
	var employee model.Employee
	db := e.storage.DB.WithContext(ctx).Where("id = ?", ID).Find(&employee)
	if db.Error != nil || employee.ID == 0 {
		e.logger.Err(db.Error).Msgf("Employee::GetEmployeeByID error: %v, (%v)", ErrRecordNotFound, db.Error)
		return employee, ErrRecordNotFound
	}

	return employee, nil
}

// GetEmployeeByContext retrieves a single row
func (e *Employee) GetEmployeeByContext(ctx context.Context, userID int) (model.Employee, error) {
	var employee model.Employee
	db := e.storage.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&employee)
	if db.Error != nil || employee.ID == 0 {
		e.logger.Err(db.Error).Msgf("Employee::GetEmployeeByID error: %v, (%v)", ErrRecordNotFound, db.Error)
		return employee, ErrRecordNotFound
	}

	return employee, nil
}

// GetAllEmployees retrieves all employees
func (e *Employee) GetAllEmployees(ctx context.Context) ([]*model.Employee, error) {
	var employees []*model.Employee
	db := e.storage.DB.WithContext(ctx).Find(&employees)
	if db.Error != nil {
		e.logger.Err(db.Error).Msgf("Employee::GetAllEmployees error: %v, (%v)", ErrRecordNotFound, db.Error)
		return nil, ErrRecordNotFound
	}

	return employees, nil
}

// UpdateEmployeeByID sets supported new values for a row accordingly
func (e *Employee) UpdateEmployeeByID(ctx context.Context, id int, employee model.Employee) (model.Employee, error) {
	db := e.storage.DB.WithContext(ctx).Model(&model.Employee{
		ID: id,
	}).UpdateColumns(model.Employee{
		FirstName:    employee.FirstName,
		LastName:     employee.LastName,
		Dob:          employee.Dob,
		DepartmentID: employee.DepartmentID,
		Position:     employee.Position,
	})
	if db.Error != nil {
		e.logger.Err(db.Error).Msgf("Employee::UpdateByID error: %v, (%v)", ErrRecordUpdateFailed, db.Error)
		return employee, ErrRecordUpdateFailed
	}
	return employee, nil
}

// DeleteEmployeeByID removes record completely from the storage
func (e *Employee) DeleteEmployeeByID(ctx context.Context, id int) error {
	db := e.storage.DB.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model.Employee{})
	if db.Error != nil {
		e.logger.Err(db.Error).Msgf("Employee::HardDeleteByID error: %v, (%v)", ErrDeleteFailed, db.Error)
		return ErrDeleteFailed
	}
	return nil
}

package controller

import (
	"context"

	"employee-management-system/model"
)

// AddEmployee returns an Employee
func (c *Controller) AddEmployee(ctx context.Context, employee model.Employee) (model.Employee, error) {
	return c.employeeStorage.AddEmployee(ctx, employee)
}

// GetEmployeeByID returns an Employee by id supplied
func (c *Controller) GetEmployeeByID(ctx context.Context, ID int) (model.Employee, error) {
	return c.employeeStorage.GetEmployeeByID(ctx, ID)
}

// GetEmployeeByContext for getting employee
func (c *Controller) GetEmployeeByContext(ctx context.Context, userID int) (model.Employee, error) {
	return c.employeeStorage.GetEmployeeByContext(ctx, userID)
}

// GetAllEmployees returns all Employees
func (c *Controller) GetAllEmployees(ctx context.Context) ([]*model.Employee, error) {
	return c.employeeStorage.GetAllEmployees(ctx)
}

// UpdateEmployeeByID for update
func (c *Controller) UpdateEmployeeByID(ctx context.Context, id int, employee model.Employee) (model.Employee, error) {
	return c.employeeStorage.UpdateEmployeeByID(ctx, id, employee)
}

// DeleteEmployeeByID for delete
func (c *Controller) DeleteEmployeeByID(ctx context.Context, id int) error {
	return c.employeeStorage.DeleteEmployeeByID(ctx, id)
}

// Package controller defines implementation that exposes logics of the app
package controller

import (
	"context"

	"github.com/rs/zerolog"

	"employee-management-system/model"
	"employee-management-system/pkg/environment"
	"employee-management-system/pkg/helper"
	"employee-management-system/pkg/middleware"
	"employee-management-system/storage"
)

const packageName = "controller"

// Operations enlist all possible operations for this controller across all modules
//
//go:generate mockgen -source controller.go -destination ./mock/mock_controller.go -package mock Operations
type Operations interface {
	Middleware() *middleware.Middleware

	AddEmployee(ctx context.Context, employee model.Employee) (model.Employee, error)
	GetEmployeeByID(ctx context.Context, ID int) (model.Employee, error)
	GetEmployeeByContext(ctx context.Context, userID int) (model.Employee, error)
	GetAllEmployees(ctx context.Context) ([]*model.Employee, error)
	UpdateEmployeeByID(ctx context.Context, id int, employee model.Employee) (model.Employee, error)
	DeleteEmployeeByID(ctx context.Context, id int) error
}

// Controller object to hold necessary reference to other dependencies
type Controller struct {
	storage         storage.Storage
	logger          zerolog.Logger
	employeeStorage storage.EmployeeDatabase
	env             *environment.Env
	middleware      *middleware.Middleware
}

// New creates a new instance of Controller
func New(z zerolog.Logger, s *storage.Storage, m *middleware.Middleware) *Operations {
	l := z.With().Str(helper.LogStrKeyModule, packageName).Logger()
	// init all storage layer here
	employee := storage.NewEmployee(s)

	ctrl := &Controller{
		storage:         *s,
		logger:          l,
		employeeStorage: *employee,
		env:             s.Env,
		middleware:      m,
	}

	op := Operations(ctrl)
	return &op
}

// Middleware returns the middleware Middleware object exposed by this app
func (c *Controller) Middleware() *middleware.Middleware {
	return c.middleware
}

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"employee-management-system/graph/model"
	"fmt"
)

// CreateEmployee is the resolver for the createEmployee field.
func (r *mutationResolver) CreateEmployee(ctx context.Context, input model.CreateEmployeeInput) (*model.Employee, error) {
	employee := &model.Employee{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  input.Password,
		Email:     input.DepartmentID,
		Dob:       input.Dob,
		Position:  input.Position,
	}

	return employee, nil
}

// UpdateEmployee is the resolver for the updateEmployee field.
func (r *mutationResolver) UpdateEmployee(ctx context.Context, id string, input model.UpdateEmployeeInput) (*model.Employee, error) {
	panic(fmt.Errorf("not implemented: UpdateEmployee - updateEmployee"))
}

// DeleteEmployee is the resolver for the deleteEmployee field.
func (r *mutationResolver) DeleteEmployee(ctx context.Context, id string) (*model.DeleteEmployeeResponse, error) {
	panic(fmt.Errorf("not implemented: DeleteEmployee - deleteEmployee"))
}

// GetAllEmployees is the resolver for the getAllEmployees field.
func (r *queryResolver) GetAllEmployees(ctx context.Context) ([]*model.Employee, error) {
	employees := []*model.Employee{
		{
			ID:        "1",
			FirstName: "Henry",
			LastName:  "Erabor",
			Email:     "henry@yahoo.com",
			Position:  "Software Engineer",
		},
		{
			ID:        "2",
			FirstName: "James",
			LastName:  "Erabor",
			Email:     "james@yahoo.com",
			Position:  "Mechanical Engineer ",
		},
	}

	return employees, nil
}

// GetEmployee is the resolver for the getEmployee field.
func (r *queryResolver) GetEmployee(ctx context.Context, id string) (*model.Employee, error) {
	employee := &model.Employee{
		ID:        "1",
		FirstName: "Henry",
		LastName:  "Erabor",
		Email:     "henry@yahoo.com",
		Position:  "Software Engineer",
	}

	return employee, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

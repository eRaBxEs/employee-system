package graph

import "database/sql"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db *sql.DB
}

// New created a new instance of Resolver
func New(db *sql.DB) *Resolver {
	return &Resolver{
		db: db,
	}
}

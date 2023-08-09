package model

import "time"

type Employee struct {
	ID           int `gorm:"column:id;PRIMARY_KEY;type:int;"`
	UserID       int
	FirstName    string
	LastName     string
	Email        string
	Dob          time.Time
	DepartmentID int `gorm:"column:departmentID;FOREIGNKEY"`
	Position     string
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

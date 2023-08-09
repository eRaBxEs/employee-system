package model

import "time"

type Department struct {
	ID             string `gorm:"column:id;PRIMARY_KEY;type:int;"`
	DepartmentName string
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

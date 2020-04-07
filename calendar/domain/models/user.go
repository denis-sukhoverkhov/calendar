package models

import "time"

type User struct {
	Id        int64
	FirstName string
	LastName  string
	Active    bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

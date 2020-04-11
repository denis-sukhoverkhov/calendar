package models

import "time"

type Event struct {
	Id        int64
	Name      string
	From      time.Time
	To        time.Time
	UserId    int64
	Active    bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

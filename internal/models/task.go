package models

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	Created_at  time.Time
	Updated_at  time.Time
}

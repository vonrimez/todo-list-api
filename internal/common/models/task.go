package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskCreateInput struct {
	Title       *string `json:"title" validate:"required,min=2,max=50"`
	Description *string `json:"description" validate:"required,min=2,max=255"`
	Status      *string `json:"status"`
}

type TaskUpdateInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

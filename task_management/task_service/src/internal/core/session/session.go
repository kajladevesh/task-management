package session

import (
	"time"
)

type Task struct {
	TaskID      int64     `json:"task_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	AssignedTo  *int64    `json:"assigned_to,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Task) AssignedToOrZero() int64 {
	if t.AssignedTo == nil {
		return 0
	}
	return *t.AssignedTo
}

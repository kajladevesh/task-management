package task

import (
	"context"
	"task_management/task_service/src/internal/core/session"
)

type Repository interface {
	CreateTask(ctx context.Context, task *session.Task) error
	UpdateTask(ctx context.Context, taskID int, task *session.Task) error
	ListTasks(ctx context.Context, assignedTo *int, status *string) ([]*session.Task, error)
	MarkTaskCompleted(ctx context.Context, taskID int64) error
}

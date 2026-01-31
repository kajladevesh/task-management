
package usecase

import (
	"context"
	"errors"
	"fmt"
	"task_management/task_service/src/internal/adaptors/persistence/redis"
	"task_management/task_service/src/internal/core/session"
	"task_management/task_service/src/internal/core/task"
)

type TaskUsecase struct {
	repo      task.Repository
	publisher *redis.RedisPublisher
}

func NewTaskUsecase(repo task.Repository, pub *redis.RedisPublisher) *TaskUsecase {
	return &TaskUsecase{repo: repo, publisher: pub}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, title, description, priority string, assignedTo *int64) (*session.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if priority == "" {
		priority = "medium"
	}
	t := &session.Task{
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		AssignedTo:  assignedTo,
	}
	err := u.repo.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}

	// Publish notification during task creation
	if t.AssignedTo != nil {
		_ = u.publisher.Publish(ctx, "task_notifications", session.TaskNotification{
			Event:      "assigned",
			TaskID:     t.TaskID,
			AssignedTo: *t.AssignedTo,
			Message:    fmt.Sprintf("Task %d assigned to user %d", t.TaskID, *t.AssignedTo),
		})
	}

	return t, nil
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, taskID int, updatedTask *session.Task) error {
	err := u.repo.UpdateTask(ctx, taskID, updatedTask)
	if err != nil {
		return err
	}

	// Publish update notification
	_ = u.publisher.Publish(ctx, "task_notifications", session.TaskNotification{
		Event:      "updated",
		TaskID:     int64(taskID),
		AssignedTo: updatedTask.AssignedToOrZero(),
		Message:    fmt.Sprintf("Task %d updated", taskID),
	})

	return nil
}

func (u *TaskUsecase) ListTasks(ctx context.Context, assignedTo *int, status *string) ([]*session.Task, error) {
	return u.repo.ListTasks(ctx, assignedTo, status)
}

func (s *TaskUsecase) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	return s.repo.MarkTaskCompleted(ctx, taskID)
}

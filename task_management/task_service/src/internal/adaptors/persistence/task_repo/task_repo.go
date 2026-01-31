package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"task_management/task_service/src/internal/core/session"
	"time"
)

type TaskRepo struct {
	DB *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{DB: db}
}

func (r *TaskRepo) CreateTask(ctx context.Context, t *session.Task) error {
	query := `
	INSERT INTO tasks (title, description, status, priority, assigned_to, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING task_id, created_at, updated_at
    `
	assignedTo := sql.NullInt64{} //NullInt64 represents an int64 that may be null. NullInt64 implements the [Scanner] interface so it can be used as a scan destination, similar to [NullString].
	if t.AssignedTo != nil {
		assignedTo = sql.NullInt64{Int64: *t.AssignedTo, Valid: true}
	}

	return r.DB.QueryRowContext(ctx, query,
		t.Title, t.Description, t.Status, t.Priority, assignedTo, time.Now(), time.Now(),
	).Scan(&t.TaskID, &t.CreatedAt, &t.UpdatedAt)
}

//----------------------------------------------------------------------------

func (r *TaskRepo) UpdateTask(ctx context.Context, taskID int, task *session.Task) error {
	_, err := r.DB.ExecContext(ctx, `
	UPDATE tasks
	SET status = $1, description = $2, priority = $3, updated_at = CURRENT_TIMESTAMP
	WHERE task_id = $4
	`, task.Status, task.Description, task.Priority, taskID)

	return err
}

//--------------------------------------------------------------------------------

func (r *TaskRepo) ListTasks(ctx context.Context, assignedTo *int, status *string) ([]*session.Task, error) {
	query := `SELECT task_id, title, description, status, priority, assigned_to, created_at, updated_at FROM tasks WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if assignedTo != nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", argIndex)
		args = append(args, *assignedTo)
		argIndex++
	}

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*session.Task
	for rows.Next() {
		var t session.Task
		if err := rows.Scan(&t.TaskID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

//-------------------------------------------------------------------------------

func (r *TaskRepo) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	result, err := r.DB.ExecContext(ctx, `
		UPDATE tasks SET status = 'completed', updated_at = CURRENT_TIMESTAMP WHERE task_id = $1
	`, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %d", taskID)
	}

	return nil
}

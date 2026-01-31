package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"task_management/task_service/src/internal/core/errors"
	"task_management/task_service/src/internal/core/session"
	"task_management/task_service/src/internal/usecase"
	pkg "task_management/task_service/src/pkg/jsonResponse"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	usecase *usecase.TaskUsecase
}

func NewTaskHandler(u *usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase: u}
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority,omitempty"`
	AssignedTo  *int64 `json:"assigned_to,omitempty"`
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	
	log.Println("🧿 [Handler] CreateTask hit by instance")

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.usecase.CreateTask(r.Context(), req.Title, req.Description, req.Priority, req.AssignedTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

type UpdateTaskInput struct {
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var input UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatedTask := &session.Task{
		Status:      input.Status,
		Description: input.Description,
		Priority:    input.Priority,
	}

	err = h.usecase.UpdateTask(r.Context(), taskID, updatedTask)
	if err != nil {
		http.Error(w, "failed to update task ", http.StatusInternalServerError)
		return
	}

	pkg.Created(w, nil, "task updated successfully")

}

//-----------------------------------------------------------------------------------------------

func (h *TaskHandler) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var assignedTo *int
	if val := query.Get("assigned_to"); val != "" {
		parsed, err := strconv.Atoi(val)
		if err == nil {
			assignedTo = &parsed
		}
	}

	var status *string
	if val := query.Get("status"); val != "" {
		status = &val
	}

	tasks, err := h.usecase.ListTasks(r.Context(), assignedTo, status)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//------------------------------------------------------------------------------------

func (h *TaskHandler) CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.usecase.MarkTaskCompleted(r.Context(), taskID)
	if err != nil {
		switch err {
		case errors.ErrTaskNotFound:
			pkg.Error(w, http.StatusNotFound, "Task not found")
		default:
			pkg.Error(w, http.StatusInternalServerError, "Failed to complete task")
		}
		return
	}

	pkg.Success(w, nil, "Task marked as completed")
}

package routes

import (
	"net/http"
	"task_management/task_service/src/internal/interfaces/input/api/rest/handler"
	"task_management/task_service/src/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
)

// func InitRoutes(userHandler *handler.TaskHandler) http.Handler {
func InitRoutes(userHandler *handler.TaskHandler, auth *middleware.AuthMiddleware) http.Handler {

	r := chi.NewRouter()

	// r.Post("/tasks", userHandler.CreateTaskHandler)
	// r.Put("/tasks/{id}", userHandler.UpdateTaskHandler)
	// r.Get("/tasks", userHandler.ListTasksHandler)
	// r.Patch("/tasks/{id}/complete", userHandler.CompleteTaskHandler)

	r.Route("/tasks", func(r chi.Router) {
		r.Use(auth.Validate) // Apply JWT validation to all /tasks routes

		r.Post("/", userHandler.CreateTaskHandler)
		r.Put("/{id}", userHandler.UpdateTaskHandler)
		r.Get("/", userHandler.ListTasksHandler)
		r.Patch("/{id}/complete", userHandler.CompleteTaskHandler)
	})

	return r
}

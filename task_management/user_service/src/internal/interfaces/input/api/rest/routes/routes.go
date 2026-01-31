package routes

import (
	"net/http"
	"task-management/user-service/src/internal/interfaces/input/api/rest/handler"
	"task-management/user-service/src/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *handler.UserHandler, jwtSecret string) http.Handler {

	r := chi.NewRouter()

	r.Post("/register", userHandler.RegisterHandler)
	r.Post("/login", userHandler.LoginHandler)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Validate)

		r.Get("/profile", userHandler.GetUserProfileHandler)
		r.Put("/update", userHandler.UpdateUserHandler)
	})

	return r
}

package main

import (
	"log"
	"net/http"
	"os"
	"task_management/task_service/src/internal/adaptors/external"
	"task_management/task_service/src/internal/config"

	persistance "task_management/task_service/src/internal/adaptors/persistence/db"
	"task_management/task_service/src/internal/adaptors/persistence/redis"
	persistence "task_management/task_service/src/internal/adaptors/persistence/task_repo"
	"task_management/task_service/src/internal/interfaces/input/api/rest/handler"
	"task_management/task_service/src/internal/interfaces/input/api/rest/middleware"
	"task_management/task_service/src/internal/interfaces/input/api/rest/routes"
	"task_management/task_service/src/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	database, err := persistance.ConnectToDatabase(config)
	if err != nil {
		log.Fatalf("could not connect to database :- %v", err)
	}

	// defer db.Close()

	taskRepo := persistence.NewTaskRepo(database.DB)
	publisher := redis.NewRedisPublisher("localhost:6379", "")
	taskUC := usecase.NewTaskUsecase(taskRepo, publisher)
	taskHandler := handler.NewTaskHandler(taskUC)

	userClient := external.NewUserServiceClient("localhost:50051")
	authMiddleware := middleware.NewAuthMiddleware(userClient)

	r := routes.InitRoutes(taskHandler, authMiddleware)

	// log.Println("Task service running on :8081")
	// if err := http.ListenAndServe(":8081", r); err != nil {
	// 	log.Fatalf("failed to start server: %v", err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // fallback if not provided
	}

	log.Printf("Task service running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

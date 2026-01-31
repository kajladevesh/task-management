package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"task-management/user-service/src/internal/adaptors/persistance"
	"task-management/user-service/src/internal/config"
	"task-management/user-service/src/internal/interfaces/input/api/rest/handler"
	"task-management/user-service/src/internal/interfaces/input/api/rest/routes"
	"task-management/user-service/src/internal/interfaces/input/grpc/grpcserver"
	"task-management/user-service/src/internal/interfaces/input/grpc/user"
	"task-management/user-service/src/internal/usecase"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
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

	userRepository := persistance.NewUserRepo(database)
	userUsecase := usecase.NewUserService(userRepository, config.JWT_SECRET)
	userHandler := handler.NewUserHandler(userUsecase)

	router := routes.InitRoutes(userHandler, config.JWT_SECRET)

	grpcServer := grpc.NewServer()
	grpcHandler := grpcserver.NewUserGRPCServer(userUsecase)

	user.RegisterUserServiceServer(grpcServer, grpcHandler)
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("user_service gRPC server running at :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	fmt.Printf("Starting server on port 8080 \n")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Failed to connect to server : %v", err)
	}
}

// user_grpc_server.go
package grpcserver

import (
	"context"
	userpb "task-management/user-service/src/internal/interfaces/input/grpc/user"
	"task-management/user-service/src/internal/usecase"
)

type UserGRPCServer struct {
	userpb.UnimplementedUserServiceServer
	usecase *usecase.UserService
}

func NewUserGRPCServer(u *usecase.UserService) *UserGRPCServer {
	return &UserGRPCServer{usecase: u}
}

func (s *UserGRPCServer) ValidateToken(ctx context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	// fmt.Println("grpc req token : ", req.Token)
	userID, err := s.usecase.ValidateToken(req.Token)

	// fmt.Println("grpc user id : ", userID)
	// fmt.Println("error : ", err)
	if err != nil {
		return &userpb.ValidateTokenResponse{Valid: false}, nil
	}
	return &userpb.ValidateTokenResponse{Valid: true, UserId: userID}, nil
}

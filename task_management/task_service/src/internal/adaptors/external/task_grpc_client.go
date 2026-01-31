package external

import (
	"context"
	"log"
	"time"

	userpb "task_management/task_service/src/internal/interfaces/output/grpc"

	"google.golang.org/grpc"
)

type UserServiceClient struct {
	client userpb.UserServiceClient
}

func NewUserServiceClient(grpcAddr string) *UserServiceClient {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure()) // for now, without TLS
	if err != nil {
		log.Fatalf("failed to connect to user_service: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserServiceClient{client: client}
}

func (u *UserServiceClient) ValidateToken(token string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// log.Println("Sending token to user_service via gRPC:", token)

	req := &userpb.ValidateTokenRequest{Token: token}
	res, err := u.client.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("gRPC error calling ValidateToken: %v", err)
		return "", false
	}

	//log.Printf("gRPC response from user_service: Valid=%v, UserID=%v", res.Valid, res.UserId)

	if res.Valid {
		return res.UserId, true
	}
	return "", false
}

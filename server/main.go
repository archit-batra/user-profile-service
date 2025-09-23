package main

import (
	"context"
	"log"
	"net"

	pb "github.com/archit-batra/user-profile-service/proto"

	"google.golang.org/grpc"
)

// Implement the server
type userServer struct {
	pb.UnimplementedUserServiceServer
}

// Example: return dummy user data
func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user := &pb.User{
		Id:    req.GetId(),
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	return &pb.GetUserResponse{User: user}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServer{})

	log.Println("Server listening on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

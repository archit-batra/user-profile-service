package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/archit-batra/user-profile-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server implements pb.UserServiceServer
type server struct {
	pb.UnimplementedUserServiceServer
}

// Sample in-memory users
var users = map[string]*pb.User{
	"1": {Id: "1", Name: "Alice", Email: "alice@example.com"},
	"2": {Id: "2", Name: "Bob", Email: "bob@example.com"},
}

// GetUser handles unary RPC requests
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received GetUser request for ID: %s", req.Id)

	// Handle context cancellation
	select {
	case <-ctx.Done():
		log.Println("Request canceled by client")
		return nil, status.Errorf(codes.Canceled, "request canceled")
	default:
		user, ok := users[req.Id]
		if !ok {
			log.Printf("User with ID %s not found", req.Id)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		log.Printf("Returning user: %v", user)
		return &pb.GetUserResponse{User: user}, nil
	}
}

// ListUsers streams multiple user records to the client
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
	log.Println("Streaming all users to client...")
	for _, user := range users {
		// Simulate delay for streaming effect
		time.Sleep(time.Millisecond * 500)
		if err := stream.Send(&pb.ListUsersResponse{User: user}); err != nil {
			return err
		}
	}
	log.Println("Finished streaming users")
	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

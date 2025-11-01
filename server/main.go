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

type server struct {
	pb.UnimplementedUserServiceServer
}

var users = map[string]*pb.User{
	"1": {Id: "1", Name: "Alice", Email: "alice@example.com"},
	"2": {Id: "2", Name: "Bob", Email: "bob@example.com"},
	"3": {Id: "3", Name: "Charlie", Email: "charlie@example.com"},
}

// Unary RPC: Get user by ID
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received GetUser request for ID: %s", req.Id)
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "request canceled")
	default:
		user, ok := users[req.Id]
		if !ok {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return &pb.GetUserResponse{User: user}, nil
	}
}

// Server streaming RPC: Stream all users
func (s *server) ListUsersStream(req *pb.ListUsersRequest, stream pb.UserService_ListUsersStreamServer) error {
	log.Println("Streaming all users to client...")
	for _, user := range users {
		// Simulate delay (for demo)
		time.Sleep(500 * time.Millisecond)
		log.Printf("Sending user: %v", user)
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	log.Println("Finished streaming users.")
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

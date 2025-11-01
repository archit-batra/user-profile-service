package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/archit-batra/user-profile-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a client connection using the modern API.
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	// ClientConn implements io.Closer; close when done
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Unary call: GetUser with RPC timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}
	log.Printf("Unary Response: %v", res.GetUser())

	// Server streaming: ListUsersStream with its own timeout
	streamCtx, streamCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer streamCancel()

	stream, err := client.ListUsersStream(streamCtx, &pb.ListUsersRequest{})
	if err != nil {
		log.Fatalf("Error calling ListUsersStream: %v", err)
	}

	log.Println("Streaming all users:")
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			// stream finished normally
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		log.Printf("Received user: %v", user)
	}
	log.Println("Stream ended.")
}

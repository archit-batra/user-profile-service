package main

import (
	"context"
	"log"
	"time"

	pb "github.com/archit-batra/user-profile-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Call GetUser API
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}

	log.Printf("User: %v", res.GetUser())

	// Call ListUsers API
	stream, err := client.ListUsers(context.Background(), &pb.ListUsersRequest{})
	if err != nil {
		log.Fatalf("Error calling ListUsers: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Received user: %v", resp.User)
	}
}

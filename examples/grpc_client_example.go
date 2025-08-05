package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"xis-data-aggregator/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)

	// Example 1: Get data by ID
	fmt.Println("=== GetDataById Example ===")
	getDataByIDExample(client)

	// Example 2: List data by time range
	fmt.Println("\n=== ListDataByTimeRange Example ===")
	listDataByTimeRangeExample(client)
}

func getDataByIDExample(client pb.DataServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create bidirectional stream
	stream, err := client.GetDataById(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Send request
	request := &pb.GetDataByIDRequest{
		Id: "123e4567-e89b-12d3-a456-426614174000", // Example UUID
	}

	if err := stream.Send(request); err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Receive response
	response, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	fmt.Printf("Received data: ID=%s, Timestamp=%d, Max=%d\n",
		response.Id, response.Timestamp, response.Max)

	// Close the stream
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close stream: %v", err)
	}
}

func listDataByTimeRangeExample(client pb.DataServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create bidirectional stream
	stream, err := client.ListDataByTimeRange(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Send request
	request := &pb.ListDataByTimeRangeRequest{
		From: "1640995200", // 2022-01-01 00:00:00 UTC
		To:   "1641081600", // 2022-01-02 00:00:00 UTC
	}

	if err := stream.Send(request); err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	// Receive response
	response, err := stream.Recv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	fmt.Printf("Received %d data items:\n", len(response.DataItems))
	for i, data := range response.DataItems {
		fmt.Printf("  [%d] ID=%s, Timestamp=%d, Max=%d\n",
			i+1, data.Id, data.Timestamp, data.Max)
	}

	// Close the stream
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close stream: %v", err)
	}
}

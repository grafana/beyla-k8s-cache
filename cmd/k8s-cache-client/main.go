package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grafana/beyla-k8s-cache/pkg/informer"
)

const (
	address = "localhost:8999"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := informer.NewEventStreamServiceClient(conn)

	// Subscribe to the event stream.
	stream, err := client.Subscribe(context.TODO(), &informer.SubscribeMessage{})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	// Receive and print messages.
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Printf("error receiving message: %v", err)
			break
		}
		fmt.Printf("Received event: %v\n", event)
	}
}

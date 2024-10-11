package main

import (
	context "context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/grafana/beyla-k8s-cache/pkg/informer"
	"github.com/grafana/beyla-k8s-cache/pkg/meta"
)

const port = 8999

type observer struct {
	msg    *informer.SubscribeMessage
	server grpc.ServerStreamingServer[informer.Event]
}

func (o *observer) ID() string {
	return o.msg.String()
}

func (o *observer) On(event *informer.Event) {
	if err := o.server.Send(event); err != nil {
		slog.Error("sending message", "clientID", o.ID(), "error", err)
	}
}

// server is used to implement informer.EventStreamServiceServer.
type server struct {
	informer.UnimplementedEventStreamServiceServer
	informers *meta.Informers
}

func (s *server) Subscribe(msg *informer.SubscribeMessage, server grpc.ServerStreamingServer[informer.Event]) error {
	o := &observer{msg: msg, server: server}
	slog.Info("subscribed component", "id", o.ID())
	s.informers.Subscribe(o)
	// Keep the connection open
	for {
		select {
		case <-server.Context().Done():
			log.Printf("Client %s disconnected", o.ID())
			s.informers.Unsubscribe(o)
			return nil
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	infors, err := meta.InitInformers(context.Background(), "", time.Minute)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	panicOnErr(err)

	s := grpc.NewServer()
	informer.RegisterEventStreamServiceServer(s, &server{informers: infors})

	log.Printf("Server listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

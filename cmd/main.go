package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"notifier/internal/twitter"
)

const (
	port = ":6565"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	twitter.RegisterTwitterServer(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

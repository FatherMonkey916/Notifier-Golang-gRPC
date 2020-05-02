package main

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"notifier/internal/twitter"
)

const (
	port     = ":6565"
	promAddr = ":9100"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s)
	twitter.RegisterTwitterServer(s)
	runPrometheus()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runPrometheus() {
	mux := http.NewServeMux()
	// Enable histogram
	grpc_prometheus.EnableHandlingTimeHistogram()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Prometheus metrics bind address", promAddr)
		log.Fatal(http.ListenAndServe(promAddr, mux))
	}()
}

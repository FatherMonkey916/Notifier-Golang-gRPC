package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
	"notifier/gateway/wsproxy"
	"notifier/grpc/notifier/twitter"
)

var (
	grpcAddr = flag.String("grpc-server", "localhost:6565", "listen grpc addr")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := twitter.RegisterTweetServiceHandlerFromEndpoint(ctx, mux, *grpcAddr, opts)
	if err != nil {
		return err
	}
	go http.ListenAndServe(":8080", nil)
	fmt.Println("listening")
	http.ListenAndServe(":8081", wsproxy.WebsocketProxy(mux))
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"math/rand"
	"net/http"
	"notifier/gateway/wsproxy"
	"notifier/grpc/notifier/twitter"
	"time"
)

var (
	grpcAddr = flag.String("grpc-server", "localhost:6565", "listen grpc addr")
	inFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "code", "method"},
	)

	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(inFlight, counter, duration, responseSize)
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	httpMux := http.NewServeMux()
	grpc_prometheus.EnableHandlingTimeHistogram()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := twitter.RegisterTweetServiceHandlerFromEndpoint(ctx, mux, *grpcAddr, opts)
	if err != nil {
		return err
	}

	aliveChain := genInstrumentChain("alive", alive)

	fmt.Println("listening")
	httpMux.Handle("/", mux)
	httpMux.Handle("/alive", aliveChain)
	httpMux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", wsproxy.WebsocketProxy(httpMux))
	return nil
}

func genInstrumentChain(name string, handler http.HandlerFunc) http.Handler {
	return promhttp.InstrumentHandlerInFlight(inFlight,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": name}),
			promhttp.InstrumentHandlerCounter(counter.MustCurryWith(prometheus.Labels{"handler": name}),
				promhttp.InstrumentHandlerResponseSize(responseSize, handler),
			),
		),
	)
}

func alive(w http.ResponseWriter, _ *http.Request) {
	dur := rand.Intn(1000)
	time.Sleep(time.Duration(dur) * time.Millisecond) // 処理を表現するためのsleep
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, "OK")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

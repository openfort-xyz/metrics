package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/openfort-xyz/jsonrpc"

	"github.com/openfort-xyz/pubsub"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// HTTPMiddleware is a HTTP middleware that records the request count and duration of the request.
func HTTPMiddleware(next http.Handler) http.Handler {
	once.Do(func() {
		prometheus.MustRegister(requestCount, requestDuration)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(srw, r)
		duration := time.Since(start)

		registerNormalized(r.Method, r.URL.Path, srw.status, duration)
	})
}

// GRPCUnaryMiddleware is a GRPC middleware that records the request count and duration of the request.
func GRPCUnaryMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	once.Do(func() {
		prometheus.MustRegister(requestCount, requestDuration)
	})

	start := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(start)
	status := 200
	if err != nil {
		status = 500
	}

	register("grpc", info.FullMethod, status, duration)
	return res, err
}

// GRPCStreamMiddleware is a GRPC middleware that records the request count and duration of the request.
func GRPCStreamMiddleware(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	once.Do(func() {
		prometheus.MustRegister(requestCount, requestDuration)
	})

	start := time.Now()
	err := handler(srv, ss)
	duration := time.Since(start)
	status := 200
	if err != nil {
		status = 500
	}

	register("grpc", info.FullMethod, status, duration)
	return err
}

// RabbitMQMiddleware is a PubSub middleware that records the request count and duration of the request.
func RabbitMQMiddleware(next pubsub.Handler) pubsub.Handler {
	once.Do(func() {
		prometheus.MustRegister(requestCount, requestDuration)
	})

	return func(ctx context.Context, event *pubsub.Event) error {
		start := time.Now()
		err := next(ctx, event)
		duration := time.Since(start)
		status := 200
		if err != nil {
			status = 500
		}

		register("rabbitmq", event.Topic.String(), status, duration)
		return err
	}
}

func JSONRPCMiddleware(next jsonrpc.Handler) jsonrpc.Handler {
	once.Do(func() {
		prometheus.MustRegister(requestCount, requestDuration)
	})

	return func(ctx context.Context, r *jsonrpc.Request) *jsonrpc.Response {
		start := time.Now()

		n := next(ctx, r)
		duration := time.Since(start)

		register("JSONRPC "+r.JSONRPC, r.Method, n.StatusCode, duration)
		return n
	}
}

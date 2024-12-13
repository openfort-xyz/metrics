package metrics

import (
	"context"
	"net/http"
	"time"

	"go.openfort.xyz/jsonrpc"

	"github.com/prometheus/client_golang/prometheus"
	"go.openfort.xyz/pubsub"
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

		register(r.Method, r.URL.Path, srw.status, duration)
	})
}

// GRPCUnaryMiddleware is a GRPC middleware that records the request count and duration of the request.
func GRPCUnaryMiddleware() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		register("GRPC", method, 200, duration)
		return err
	}
}

// RabbitMQMiddleware is a PubSub middleware that records the request count and duration of the request.
func RabbitMQMiddleware(next pubsub.Handler) pubsub.Handler {
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
	return func(ctx context.Context, r *jsonrpc.Request) *jsonrpc.Response {
		start := time.Now()

		n := next(ctx, r)
		duration := time.Since(start)

		register("JSONRPC "+r.JSONRPC, r.Method, n.StatusCode, duration)
		return n
	}
}

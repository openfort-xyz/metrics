package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ExposeHTTP returns a http.Handler that exposes the metrics
func ExposeHTTP() http.Handler {
	return promhttp.Handler()
}

// Server is a metrics server
type Server struct {
	server *http.Server
}

// NewServer creates a new metrics server
func NewServer(port int) *Server {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}
}

// Start starts the metrics server
func (s *Server) Start(ctx context.Context) error {
	http.Handle("/metrics", ExposeHTTP())
	return s.server.ListenAndServe()
}

// Stop stops the metrics server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

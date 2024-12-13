package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ExposeHTTP returns a http.Handler that exposes the metrics
func ExposeHTTP() http.Handler {
	return promhttp.Handler()
}

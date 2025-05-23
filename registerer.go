package metrics

import (
	"strconv"
	"time"
)

func register(method, path string, status int, duration time.Duration) {
	requestCount.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
	requestDuration.WithLabelValues(method, path, strconv.Itoa(status)).Observe(duration.Seconds())
}

func registerNormalized(method, path string, status int, duration time.Duration) {
	register(method, normalizer.normalize(path), status, duration)
}

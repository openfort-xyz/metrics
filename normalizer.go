package metrics

import (
	"regexp"
)

// Some of our measurements contains label with high-cardinality values,
// saturating Prometheus. These values need to be normalized.
type normalization struct {
	pattern *regexp.Regexp
	exp     string
}

var normalizations = [2]normalization{
	// UUIDs.
	{
		pattern: regexp.MustCompile(`/\/[\w]*[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}/g`),
		exp:     "/:uuid",
	},
	// Tokens.
	{
		pattern: regexp.MustCompile(`/[\w-]{24,}`),
		exp:     "/:token",
	},
}

func normalize(p string) string {
	if p == "" || p == "/" {
		return p
	}

	for _, r := range normalizations {
		p = r.pattern.ReplaceAllString(p, r.exp)
	}
	return p
}

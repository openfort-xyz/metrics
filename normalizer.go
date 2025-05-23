package metrics

import (
	"regexp"
)

// Some of our measurements contains label with high-cardinality values,
// saturating Prometheus. These values need to be normalized.

var normalizer = pathNormalizer{
	uuids: normalization{
		pattern: regexp.MustCompile(`\/[\w]*[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}`),
		exp:     "/:uuid",
	},
	tokens: normalization{
		pattern: regexp.MustCompile(`\/[\w-]{24,}`),
		exp:     "/:token",
	},
}

type normalization struct {
	pattern *regexp.Regexp
	exp     string
}

type pathNormalizer struct {
	uuids  normalization
	tokens normalization
}

func (pn pathNormalizer) normalize(p string) string {
	if p == "" || p == "/" {
		return p
	}

	p = pn.uuids.pattern.ReplaceAllString(p, pn.uuids.exp)
	return pn.tokens.pattern.ReplaceAllString(p, pn.tokens.exp)
}

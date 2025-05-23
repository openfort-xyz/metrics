package metrics

import (
	"regexp"
)

var normalizer = newPathNormalizer()

type pathNormalizer struct {
	uuid  replacement
	token replacement
}

type replacement struct {
	pattern *regexp.Regexp
	exp     string
}

func newPathNormalizer() pathNormalizer {
	return pathNormalizer{
		uuid: replacement{
			pattern: regexp.MustCompile(`/\/[\w]*[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}/g`),
			exp:     "/:uuid",
		},
		token: replacement{
			pattern: regexp.MustCompile(`/[\w-]{24,}`),
			exp:     "/:token",
		},
	}
}

func (pn pathNormalizer) path(p string) string {
	if p == "" || p == "/" {
		return p
	}

	p = pn.uuid.pattern.ReplaceAllString(p, pn.uuid.exp)
	p = pn.token.pattern.ReplaceAllString(p, pn.token.exp)

	return p
}

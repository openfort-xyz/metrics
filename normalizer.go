package metrics

import (
	"regexp"
)

var normalizer = newPathNormalizer()

type pathNormalizer struct {
	ids   replacement
	uuid  replacement
	token replacement
}

type replacement struct {
	pattern *regexp.Regexp
	exp     string
}

func newPathNormalizer() pathNormalizer {
	return pathNormalizer{
		ids: replacement{
			pattern: regexp.MustCompile(`/\/\d+/g`),
			exp:     "/:id",
		},
		uuid: replacement{
			pattern: regexp.MustCompile(`/\/[\w]*[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}/g`),
			exp:     "/:uuid",
		},
		token: replacement{
			pattern: regexp.MustCompile(`/[A-Z0-9_-]{24,}`),
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
	p = pn.ids.pattern.ReplaceAllString(p, pn.ids.exp)

	return p
}

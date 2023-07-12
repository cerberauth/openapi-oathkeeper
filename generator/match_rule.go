package generator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
	"github.com/ory/oathkeeper/rule"
)

var (
	argre = regexp.MustCompile(`(?m)({(.*)})`)

	separatorToken = rex.Chars.Single('/')
	stringToken    = rex.Group.NonCaptured(
		rex.Chars.Alphanumeric().Repeat().ZeroOrOne(),
		rex.Chars.Single('-').Repeat().ZeroOrOne(),
		rex.Chars.Single('=').Repeat().ZeroOrOne(),
		rex.Chars.Single('?').Repeat().ZeroOrOne(),
		rex.Chars.Single('&').Repeat().ZeroOrOne(),
		rex.Chars.Single('_').Repeat().ZeroOrOne(),
	)
)

func _hasQueryParam(p *openapi3.Parameters) bool {
	if p == nil {
		return false
	}

	for _, param := range *p {
		if param.Value.In == "query" {
			return true
		}
	}

	return false
}

func createMatchRule(serverUrls []string, v string, p string, params *openapi3.Parameters) (*rule.Match, error) {
	if len(serverUrls) == 0 {
		return nil, errors.New("a matching rule must has at least one server url")
	}

	var serverUrlsTokens []dialect.Token
	for _, serverUrl := range serverUrls {
		serverUrlsTokens = append(serverUrlsTokens, rex.Common.Text(serverUrl))
	}

	var pathTokens []dialect.Token
	pathTokens = append(pathTokens, separatorToken)
	for _, dir := range strings.Split(p, "/") {
		if dir == "" {
			continue
		}

		dirToken := rex.Common.Text(dir)
		if argre.Match([]byte(dir)) {
			dirToken = stringToken.Repeat().OneOrMore()
		}

		pathTokens = append(pathTokens, dirToken)
		pathTokens = append(pathTokens, separatorToken)
	}

	if len(pathTokens) > 1 {
		pathTokens[len(pathTokens)-1] = separatorToken.Repeat().ZeroOrOne()
	}

	if _hasQueryParam(params) {
		pathTokens = append(pathTokens, rex.Group.Define(
			rex.Chars.Single('?'),
			rex.Chars.Any().Repeat().OneOrMore(),
		).Repeat().ZeroOrOne())
	}

	url := rex.New(
		rex.Chars.Begin(),
		rex.Group.Composite(serverUrlsTokens...),
		rex.Group.Define(pathTokens...),
		rex.Chars.End(),
	).MustCompile()

	match := &rule.Match{
		// URL Regexp is encapsulated in brackets: https://www.ory.sh/docs/oathkeeper/api-access-rules#access-rule-format
		URL:     "<" + url.String() + ">",
		Methods: []string{v},
	}

	return match, nil
}

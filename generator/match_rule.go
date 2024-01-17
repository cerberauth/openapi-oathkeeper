package generator

import (
	"log"
	"regexp"
	"strings"

	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
)

var (
	argre = regexp.MustCompile(`(?m)({(.*)})`)

	numberToken = rex.Group.Define(
		rex.Group.Composite(
			rex.Chars.Single('-'),
			rex.Chars.Single('+'),
		).Repeat().ZeroOrOne(),
		rex.Chars.Digits().Repeat().OneOrMore(),
		rex.Group.NonCaptured(
			rex.Chars.Single('.'),
			rex.Chars.Digits().Repeat().OneOrMore(),
		).Repeat().ZeroOrOne(),
	)
	integerToken = rex.Chars.Digits().Repeat().OneOrMore()
	stringToken  = rex.Chars.Any().Repeat().OneOrMore()
	defaultToken = stringToken
)

func encapsulateRegex(r *regexp.Regexp) string {
	// URL Regexp is encapsulated in brackets: https://www.ory.sh/docs/oathkeeper/api-access-rules#access-rule-format
	return "<" + r.String() + ">"
}

func encapsulateRegexToken(t dialect.Token) string {
	return encapsulateRegex(rex.New(t).MustCompile())
}

func createServerUrlMatchingGroup(serverUrls []string) string {
	if len(serverUrls) == 0 {
		return ""
	}

	if len(serverUrls) == 1 {
		return strings.TrimSuffix(serverUrls[0], "/")
	}

	var serverUrlsTokens []dialect.Token
	for _, serverUrl := range serverUrls {
		serverUrlsTokens = append(serverUrlsTokens, rex.Common.Text(serverUrl))
	}

	return encapsulateRegexToken(rex.Group.Composite(serverUrlsTokens...))
}

func getPathParamType(name string, params *openapi3.Parameters) string {
	if params == nil {
		log.Default().Print("no path parameters has been defined")
		return ""
	}

	p := params.GetByInAndName(openapi3.ParameterInPath, name)
	if p == nil {
		log.Default().Printf("path param %s is not defined", name)
		return ""
	}

	return p.Schema.Value.Type
}

func createParamsMatchingGroup(name string, params *openapi3.Parameters) string {
	var t dialect.Token
	switch getPathParamType(name, params) {
	case "string":
		t = stringToken
	case "number":
		t = numberToken
	case "integer":
		t = integerToken
	default:
		t = defaultToken
	}

	return encapsulateRegexToken(t)
}

func createMatchRule(serverUrls []string, v string, p string, params *openapi3.Parameters) (*oathkeeper.RuleMatch, error) {
	pathTokens := []string{createServerUrlMatchingGroup(serverUrls)}
	for _, dir := range strings.Split(p, "/") {
		if dir == "" {
			continue
		}

		if matches := argre.FindStringSubmatch(dir); len(matches) > 0 {
			pathTokens = append(pathTokens, createParamsMatchingGroup(string(matches[2]), params))
		} else {
			pathTokens = append(pathTokens, dir)
		}
	}

	u := strings.Join(pathTokens, "/")
	match := &oathkeeper.RuleMatch{
		URL:     u,
		Methods: []string{v},
	}

	return match, nil
}

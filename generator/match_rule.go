package generator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/hedhyw/rex/pkg/dialect"
	"github.com/hedhyw/rex/pkg/rex"
	"github.com/ory/oathkeeper/rule"
)

var argre = regexp.MustCompile(`(?m)({(.*)})`)

func createMatchRule(serverUrls []string, v string, p string) (*rule.Match, error) {
	if len(serverUrls) == 0 {
		return nil, errors.New("a matching rule must has at least one server url")
	}

	var serverUrlsTokens []dialect.Token
	for _, serverUrl := range serverUrls {
		serverUrlsTokens = append(serverUrlsTokens, rex.Common.Text(serverUrl))
	}

	var pathTokens []dialect.Token
	pathTokens = append(pathTokens, rex.Chars.Single('/'))
	for _, dir := range strings.Split(p, "/") {
		if dir == "" {
			continue
		}

		dirToken := rex.Common.Text(dir)
		if argre.Match([]byte(dir)) {
			dirToken = rex.Group.Define(
				rex.Chars.Any().Repeat().OneOrMore(),
			)
		}

		pathTokens = append(pathTokens, dirToken)
		pathTokens = append(pathTokens, rex.Chars.Single('/'))
	}

	if len(pathTokens) > 1 {
		pathTokens[len(pathTokens)-1] = rex.Chars.Single('/').Repeat().ZeroOrOne()
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

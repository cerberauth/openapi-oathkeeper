package generator

import (
	"context"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Generator struct {
	doc *openapi3.T
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Document(d *openapi3.T) *Generator {
	g.doc = d

	return g
}

func (g *Generator) Generate(ctx context.Context) ([]rule.Rule, error) {
	err := g.doc.Validate(ctx)
	if err != nil {
		return nil, err
	}

	rules := []rule.Rule{}
	for _, s := range g.doc.Servers {
		for path, p := range g.doc.Paths {
			for verb, o := range p.Operations() {
				matchUrl, err := url.JoinPath(s.URL, path)
				if err != nil {
					return nil, err
				}

				authenticator := &AuthenticatorNoop{}

				rules = append(rules, rule.Rule{
					ID:          o.OperationID,
					Description: o.Description,
					Match: &rule.Match{
						URL: matchUrl,
						Methods: []string{
							verb,
						},
					},
					Authenticators: []rule.Handler{
						*authenticator.CreateAuthenticator(o),
					},
					Authorizer: rule.Handler{
						Handler: "allow",
					},
					Mutators: []rule.Handler{
						{
							Handler: "noop",
						},
					},
				})
			}
		}
	}

	return rules, nil
}

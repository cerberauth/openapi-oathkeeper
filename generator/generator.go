package generator

import (
	"context"
	"encoding/json"

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

func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	err := g.doc.Validate(ctx)
	if err != nil {
		return nil, err
	}

	rules := []rule.Rule{}
	// for _, p := range g.doc.Paths {
	// 	for _, o := range p.Operations() {
	// 		basePath, _ := o.Servers.BasePath()
	// 		matchUrl, _ := url.JoinPath(basePath, o.ExternalDocs.URL)
	// 		rules = append(rules, rule.Rule{
	// 			ID:          o.OperationID,
	// 			Description: o.Description,
	// 			Match:       &rule.Match{URL: matchUrl},
	// 		})
	// 	}
	// }

	return json.Marshal(rules)
}

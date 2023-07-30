package generator

import (
	"context"
	"sort"

	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Generator struct {
	doc *openapi3.T
	cfg *config.Config

	authenticators map[string]authenticator.Authenticator
}

type RulesById []rule.Rule

func (r RulesById) Len() int           { return len(r) }
func (r RulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RulesById) Less(i, j int) bool { return r[i].GetID() < r[j].GetID() }

func (g *Generator) computeId(operationId string) string {
	if g.cfg.Prefix == "" {
		return operationId
	}

	return g.cfg.Prefix + ":" + operationId
}

func (g *Generator) createRule(verb string, path string, o *openapi3.Operation) (*rule.Rule, error) {
	match, matchRuleErr := createMatchRule(g.cfg.ServerUrls, verb, path, &o.Parameters)
	if matchRuleErr != nil {
		return nil, matchRuleErr
	}

	var authenticators = []rule.Handler{}
	appendAuthenticator := func(sr *openapi3.SecurityRequirements) error {
		for _, s := range *sr {
			for k := range s {
				if a, ok := g.authenticators[k]; ok {
					ar, arerror := a.CreateAuthenticator(&s)
					if arerror != nil {
						return arerror
					}
					authenticators = append(authenticators, *ar)
				}
			}
		}

		return nil
	}

	if o.Security != nil && len(*o.Security) > 0 {
		appendAuthenticator(o.Security)
	} else if g.doc.Security != nil && len(g.doc.Security) > 0 {
		appendAuthenticator(&g.doc.Security)
	} else {
		ar, arerror := g.authenticators[string(authenticator.AuthenticatorTypeNoop)].CreateAuthenticator(nil)
		if arerror != nil {
			return nil, arerror
		}

		authenticators = append(authenticators, *ar)
	}

	return &rule.Rule{
		ID:             g.computeId(o.OperationID),
		Description:    o.Description,
		Match:          match,
		Upstream:       g.cfg.Upstream,
		Authenticators: authenticators,
		Authorizer: rule.Handler{
			Handler: "allow",
		},
		Mutators: g.cfg.Mutators,
		Errors:   g.cfg.Errors,
	}, nil
}

func createAuthenticators(d *openapi3.T, cfg *config.Config) (map[string]authenticator.Authenticator, error) {
	authenticators := make(map[string]authenticator.Authenticator)

	// Create a first authenticator for operations without security configured
	authenticators[string(authenticator.AuthenticatorTypeNoop)] = &authenticator.AuthenticatorNoop{}

	newAuthenticator := func(name string) (authenticator.Authenticator, error) {
		s := d.Components.SecuritySchemes[name]
		v, ok := cfg.Authenticators[name]
		if !ok {
			return authenticator.NewAuthenticatorFromSecurityScheme(s, nil)
		}

		return authenticator.NewAuthenticatorFromSecurityScheme(s, &v)
	}

	if d.Components != nil && d.Components.SecuritySchemes != nil {
		for name := range d.Components.SecuritySchemes {
			a, err := newAuthenticator(name)
			if err != nil {
				return nil, err
			}
			authenticators[name] = a
		}
	}

	return authenticators, nil
}

func NewGenerator(ctx context.Context, d *openapi3.T, cfg *config.Config) (*Generator, error) {
	if validateErr := d.Validate(ctx, openapi3.DisableExamplesValidation(), openapi3.DisableSchemaDefaultsValidation()); validateErr != nil {
		return nil, validateErr
	}

	if cfg.ServerUrls == nil {
		for _, s := range d.Servers {
			cfg.ServerUrls = append(cfg.ServerUrls, s.URL)
		}
	}

	authenticators, err := createAuthenticators(d, cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		doc: d,
		cfg: cfg,

		authenticators: authenticators,
	}, nil
}

func (g *Generator) Generate() ([]rule.Rule, error) {
	rules := []rule.Rule{}
	for path, p := range g.doc.Paths {
		for verb, o := range p.Operations() {
			rule, createRuleErr := g.createRule(verb, path, o)
			if createRuleErr != nil {
				return nil, createRuleErr
			}

			rules = append(rules, *rule)
		}
	}

	sort.Sort(RulesById(rules))
	return rules, nil
}

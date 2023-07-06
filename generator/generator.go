package generator

import (
	"context"
	"sort"

	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Generator struct {
	doc *openapi3.T

	authenticators map[string]authenticator.Authenticator
	PrefixId       string
	serverUrls     []string
	upstream       *rule.Upstream
}

type RulesById []rule.Rule

func (r RulesById) Len() int           { return len(r) }
func (r RulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RulesById) Less(i, j int) bool { return r[i].ID < r[j].ID }

func (g *Generator) computeId(operationId string) string {
	if g.PrefixId == "" {
		return operationId
	}

	return g.PrefixId + ":" + operationId
}

func (g *Generator) createRule(verb string, path string, o *openapi3.Operation) (*rule.Rule, error) {
	match, matchRuleErr := createMatchRule(g.serverUrls, verb, path, &o.Parameters)
	if matchRuleErr != nil {
		return nil, matchRuleErr
	}

	rule := rule.Rule{
		ID:             g.computeId(o.OperationID),
		Description:    o.Description,
		Match:          match,
		Upstream:       *g.upstream,
		Authenticators: []rule.Handler{},
		Authorizer: rule.Handler{
			Handler: "allow",
		},
		Mutators: []rule.Handler{
			{
				Handler: "noop",
			},
		},
		Errors: []rule.ErrorHandler{
			{
				Handler: "json",
			},
		},
	}

	appendAuthenticator := func(sr *openapi3.SecurityRequirements) error {
		for _, s := range *sr {
			for k := range s {
				if a, ok := g.authenticators[k]; ok {
					ar, arerror := a.CreateAuthenticator(&s)
					if arerror != nil {
						return arerror
					}
					rule.Authenticators = append(rule.Authenticators, *ar)
				}
			}
		}

		return nil
	}

	if o.Security != nil {
		appendAuthenticator(o.Security)
	} else if g.doc.Security != nil {
		appendAuthenticator(&g.doc.Security)
	} else {
		ar, arerror := g.authenticators[string(authenticator.AuthenticatorTypeNoop)].CreateAuthenticator(nil)
		if arerror != nil {
			return nil, arerror
		}

		rule.Authenticators = append(rule.Authenticators, *ar)
	}

	return &rule, nil
}

func NewGenerator(ctx context.Context, d *openapi3.T, prefixId string, jwksUris map[string]string, allowedIssuers map[string]string, allowedAudiences map[string]string, serverUrls []string, upstreamUrl string, upstreamStripPath string) (*Generator, error) {
	var upstream = rule.Upstream{}
	if upstreamUrl != "" {
		upstream.URL = upstreamUrl
	}

	if upstreamStripPath != "" {
		upstream.StripPath = upstreamStripPath
	}

	if validateErr := d.Validate(ctx, openapi3.DisableExamplesValidation(), openapi3.DisableSchemaDefaultsValidation()); validateErr != nil {
		return nil, validateErr
	}

	if serverUrls == nil {
		for _, s := range d.Servers {
			serverUrls = append(serverUrls, s.URL)
		}
	}

	authenticators := map[string]authenticator.Authenticator{}
	authenticators[string(authenticator.AuthenticatorTypeNoop)] = &authenticator.AuthenticatorNoop{}
	if d.Components.SecuritySchemes != nil {
		for ssn, ss := range d.Components.SecuritySchemes {
			var jwksUri, allowedIssuer, allowedAudience *string = nil, nil, nil

			if uri, ok := jwksUris[ssn]; ok {
				jwksUri = &uri
			}

			if iss, ok := allowedIssuers[ssn]; ok {
				allowedIssuer = &iss
			}

			if aud, ok := allowedAudiences[ssn]; ok {
				allowedAudience = &aud
			}

			a, err := NewAuthenticatorFromSecurityScheme(ss, jwksUri, allowedIssuer, allowedAudience)
			if err != nil {
				return nil, err
			}
			authenticators[ssn] = a
		}
	}

	return &Generator{
		doc: d,

		authenticators: authenticators,
		PrefixId:       prefixId,
		serverUrls:     serverUrls,
		upstream:       &upstream,
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

package generator

import (
	"context"
	"errors"
	"sort"

	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Generator struct {
	doc            *openapi3.T
	authenticators map[string]authenticator.Authenticator
	PrefixId       string

	serverUrls       []string
	jwksUris         map[string]string
	allowedIssuers   map[string]string
	allowedAudiences map[string]string

	upstream *rule.Upstream
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

func NewGenerator(prefixId string, jwksUris map[string]string, allowedIssuers map[string]string, allowedAudiences map[string]string, serverUrls []string, upstreamUrl string, upstreamStripPath string) *Generator {
	var upstream = rule.Upstream{}
	if upstreamUrl != "" {
		upstream.URL = upstreamUrl
	}

	if upstreamStripPath != "" {
		upstream.StripPath = upstreamStripPath
	}

	return &Generator{
		PrefixId: prefixId,

		jwksUris:         jwksUris,
		allowedIssuers:   allowedIssuers,
		allowedAudiences: allowedAudiences,
		serverUrls:       serverUrls,

		upstream: &upstream,
	}
}

func (g *Generator) LoadOpenAPI3Doc(ctx context.Context, d *openapi3.T) error {
	g.doc = d

	if validateErr := g.doc.Validate(ctx, openapi3.DisableExamplesValidation(), openapi3.DisableSchemaDefaultsValidation()); validateErr != nil {
		return validateErr
	}

	if g.serverUrls == nil {
		for _, s := range g.doc.Servers {
			g.serverUrls = append(g.serverUrls, s.URL)
		}
	}

	authenticators, createAuthErr := g.createAuthenticators(g.doc)
	if createAuthErr != nil {
		return createAuthErr
	}

	g.authenticators = authenticators
	return nil
}

func (g *Generator) getSSJwksUri(ssn string) (string, error) {
	jwksUri, jwksUriExists := (g.jwksUris)[ssn]
	if !jwksUriExists {
		return "", errors.New("no jwksUris found for a given security scheme")
	}

	return jwksUri, nil
}

func (g *Generator) getSSIssuer(ssn string) (string, error) {
	issuer, issuerExists := (g.allowedIssuers)[ssn]
	if !issuerExists {
		return "", errors.New("no issuer found for a given security scheme")
	}

	return issuer, nil
}

func (g *Generator) getSSAudience(ssn string) (string, error) {
	audience, audienceExist := (g.allowedAudiences)[ssn]
	if !audienceExist || audience == "" {
		return "", errors.New("no audience found for a given security scheme")
	}

	return audience, nil
}

func (g *Generator) createAuthenticators(doc *openapi3.T) (map[string]authenticator.Authenticator, error) {
	authenticators := map[string]authenticator.Authenticator{}
	authenticators[string(authenticator.AuthenticatorTypeNoop)] = &authenticator.AuthenticatorNoop{}
	var err error
	for ssn, ss := range doc.Components.SecuritySchemes {
		sstype := ss.Value.Type
		switch sstype {
		case string(authenticator.AuthenticatorTypeOpenIdConnect):
			audience, _ := g.getSSAudience(ssn)

			authenticators[ssn], err = authenticator.NewAuthenticatorOpenIdConnect(ss, audience)
		case string(authenticator.AuthenticatorTypeOAuth2):
			jwksUri, jwksUriErr := g.getSSJwksUri(ssn)
			if jwksUriErr != nil {
				return nil, jwksUriErr
			}

			issuer, issuerErr := g.getSSIssuer(ssn)
			if issuerErr != nil {
				return nil, issuerErr
			}

			audience, _ := g.getSSAudience(ssn)

			authenticators[ssn], err = authenticator.NewAuthenticatorOAuth2(ss, jwksUri, issuer, audience)
		case string(authenticator.AuthenticatorTypeHttp):
			if ss.Value.Scheme != "bearer" {
				return nil, errors.New("http security scheme must be bearer")
			}

			jwksUri, jwksUriErr := g.getSSJwksUri(ssn)
			if jwksUriErr != nil {
				return nil, jwksUriErr
			}

			issuer, issuerErr := g.getSSIssuer(ssn)
			if issuerErr != nil {
				return nil, issuerErr
			}

			audience, audienceErr := g.getSSAudience(ssn)
			if audienceErr != nil {
				return nil, audienceErr
			}

			authenticators[ssn], err = authenticator.NewAuthenticatorHttpBearer(ss, jwksUri, issuer, audience)

		default:
			return nil, errors.New("unknown security scheme")
		}
	}

	return authenticators, err
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

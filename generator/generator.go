package generator

import (
	"context"
	"errors"
	"net/url"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Generator struct {
	doc            *openapi3.T
	authenticators map[string]Authenticator
	PrefixId       string

	JwksUris       *map[string]string
	AllowedIssuers *map[string]string
}

type AuthenticatorType string

const (
	AuthenticatorTypeNoop          AuthenticatorType = "noop"
	AuthenticatorTypeOpenIdConnect AuthenticatorType = "openIdConnect"
	AuthenticatorTypeOAuth2        AuthenticatorType = "oauth2"
	AuthenticatorTypeHttp          AuthenticatorType = "http"
)

var argre = regexp.MustCompile(`(?m)({(.*)})`)

func (g *Generator) computeId(operationId string) string {
	if g.PrefixId == "" {
		return operationId
	}

	return g.PrefixId + ":" + operationId
}

func (g *Generator) createRule(verb string, path string, s *openapi3.Server, o *openapi3.Operation) (*rule.Rule, error) {
	joinUrl, joinErr := url.JoinPath(s.URL, argre.ReplaceAllString(path, string("<.*>")))
	if joinErr != nil {
		return nil, joinErr
	}

	globalUrl, unescapedErr := url.PathUnescape(joinUrl)
	if unescapedErr != nil {
		return nil, unescapedErr
	}

	rule := rule.Rule{
		ID:          g.computeId(o.OperationID),
		Description: o.Description,
		Match: &rule.Match{
			URL:     globalUrl,
			Methods: []string{verb},
		},
		Authenticators: []rule.Handler{},
		Authorizer: rule.Handler{
			Handler: "allow",
		},
		Mutators: []rule.Handler{
			{
				Handler: "noop",
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
		ar, arerror := g.authenticators[string(AuthenticatorTypeNoop)].CreateAuthenticator(nil)
		if arerror != nil {
			return nil, arerror
		}

		rule.Authenticators = append(rule.Authenticators, *ar)
	}

	return &rule, nil
}

func NewGenerator(prefixId string, jwksUris *map[string]string, allowedIssuers *map[string]string) *Generator {
	return &Generator{
		PrefixId: prefixId,

		JwksUris:       jwksUris,
		AllowedIssuers: allowedIssuers,
	}
}

func (g *Generator) LoadOpenAPI3Doc(ctx context.Context, d *openapi3.T) error {
	g.doc = d

	if validateErr := g.doc.Validate(ctx); validateErr != nil {
		return validateErr
	}

	authenticators, createAuthErr := g.createAuthenticators(g.doc)
	if createAuthErr != nil {
		return createAuthErr
	}

	g.authenticators = authenticators
	return nil
}

func (g *Generator) getSSJwksUri(ssn string) (string, error) {
	jwksUri, jwksUriExists := (*g.JwksUris)[ssn]
	if !jwksUriExists {
		return "", errors.New("no jwksUris for a given security scheme")
	}

	return jwksUri, nil
}

func (g *Generator) getSSIssuer(ssn string) (string, error) {
	issuer, issuerExists := (*g.AllowedIssuers)[ssn]
	if !issuerExists {
		return "", errors.New("no issuer for a given security scheme")
	}

	return issuer, nil
}

func (g *Generator) createAuthenticators(doc *openapi3.T) (map[string]Authenticator, error) {
	authenticators := map[string]Authenticator{}
	authenticators[string(AuthenticatorTypeNoop)] = &AuthenticatorNoop{}
	var err error
	for ssn, ss := range doc.Components.SecuritySchemes {
		sstype := ss.Value.Type
		switch sstype {
		case string(AuthenticatorTypeOpenIdConnect):
			authenticators[ssn], err = NewAuthenticatorOpenIdConnect(ss)
		case string(AuthenticatorTypeOAuth2):
			jwksUri, jwksUriErr := g.getSSJwksUri(ssn)
			if jwksUriErr != nil {
				return nil, jwksUriErr
			}

			issuer, issuerErr := g.getSSIssuer(ssn)
			if issuerErr != nil {
				return nil, issuerErr
			}

			authenticators[ssn], err = NewAuthenticatorOAuth2(ss, jwksUri, issuer)
		case string(AuthenticatorTypeHttp):
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

			authenticators[ssn], err = NewAuthenticatorHttpBearer(ss, jwksUri, issuer)

		default:
			return nil, errors.New("unknown security scheme")
		}
	}

	return authenticators, err
}

func (g *Generator) Generate() ([]rule.Rule, error) {
	rules := []rule.Rule{}
	for _, s := range g.doc.Servers {
		for path, p := range g.doc.Paths {
			for verb, o := range p.Operations() {
				rule, createRuleErr := g.createRule(verb, path, s, o)
				if createRuleErr != nil {
					return nil, createRuleErr
				}

				rules = append(rules, *rule)
			}
		}
	}

	return rules, nil
}

package authenticator

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorOAuth2)(nil)

type AuthenticatorOAuth2 struct {
	jwksUri  string
	issuer   string
	audience *string
}

func createRulesFromOAuth2SecurityRequirement(s *openapi3.SecurityRequirement, jwksUri string, issuer string, audience *string) (*rule.Handler, error) {
	ta := []string{}
	if audience != nil {
		ta = append(ta, *audience)
	}

	c := JWTAuthenticatorConfig{
		JwksUrls:       []string{jwksUri},
		TrustedIssuers: []string{issuer},
		RequiredScope:  []string{},
		TargetAudience: ta,
	}

	for _, scope := range *s {
		c.RequiredScope = append(c.RequiredScope, scope...)
	}

	jsonConfig, jsonErr := json.Marshal(c)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}, nil
}

func NewAuthenticatorOAuth2(s *openapi3.SecuritySchemeRef, jwksUri string, issuer string, audience *string) (*AuthenticatorOAuth2, error) {
	return &AuthenticatorOAuth2{
		jwksUri:  jwksUri,
		issuer:   issuer,
		audience: audience,
	}, nil
}

func (a *AuthenticatorOAuth2) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	return createRulesFromOAuth2SecurityRequirement(s, a.jwksUri, a.issuer, a.audience)
}

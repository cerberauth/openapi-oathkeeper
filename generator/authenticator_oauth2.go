package generator

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorOAuth2)(nil)

type AuthenticatorOAuth2 struct {
	JwksUri string
	Issuer  string
}

func createRulesFromOAuth2SecurityRequirement(s *openapi3.SecurityRequirement, jwksUri string, issuer string) (*rule.Handler, error) {
	c := JWTAuthenticatorConfig{
		JwksUrls:       []string{jwksUri},
		TrustedIssuers: []string{issuer},
		RequiredScope:  []string{},
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

func NewAuthenticatorOAuth2(s *openapi3.SecuritySchemeRef, jwksUri string, issuer string) (*AuthenticatorOAuth2, error) {
	return &AuthenticatorOAuth2{
		JwksUri: jwksUri,
		Issuer:  issuer,
	}, nil
}

func (a *AuthenticatorOAuth2) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	return createRulesFromOAuth2SecurityRequirement(s, a.JwksUri, a.Issuer)
}

package generator

import (
	"encoding/json"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorHttpBearer)(nil)

type AuthenticatorHttpBearer struct {
	JwksUri string
	Issuer  string
}

func NewAuthenticatorHttpBearer(s *openapi3.SecuritySchemeRef, jwksUri string, issuer string) (*AuthenticatorHttpBearer, error) {
	if s.Value.BearerFormat != "JWT" {
		return nil, errors.New("bearer format must be jwt")
	}

	return &AuthenticatorHttpBearer{
		JwksUri: jwksUri,
		Issuer:  issuer,
	}, nil
}

func (a *AuthenticatorHttpBearer) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	c := JWTAuthenticatorConfig{
		JwksUrls:       []string{a.JwksUri},
		TrustedIssuers: []string{a.Issuer},
		RequiredScope:  []string{},
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

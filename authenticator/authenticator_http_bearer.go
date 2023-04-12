package authenticator

import (
	"encoding/json"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorHttpBearer)(nil)

type AuthenticatorHttpBearer struct {
	JwksUri  string
	Issuer   string
	Audience string
}

func NewAuthenticatorHttpBearer(s *openapi3.SecuritySchemeRef, jwksUri string, issuer string, audience string) (*AuthenticatorHttpBearer, error) {
	if s.Value.BearerFormat != "JWT" {
		return nil, errors.New("bearer format must be jwt")
	}

	return &AuthenticatorHttpBearer{
		JwksUri:  jwksUri,
		Issuer:   issuer,
		Audience: audience,
	}, nil
}

func (a *AuthenticatorHttpBearer) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	ta := []string{}
	if a.Audience != "" {
		ta = append(ta, a.Audience)
	}

	c := JWTAuthenticatorConfig{
		JwksUrls:       []string{a.JwksUri},
		TrustedIssuers: []string{a.Issuer},
		RequiredScope:  []string{},
		TargetAudience: ta,
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

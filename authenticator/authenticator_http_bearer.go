package authenticator

import (
	"encoding/json"
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorHttpBearer)(nil)

type AuthenticatorHttpBearer struct {
	jwksUri  string
	issuer   string
	audience string
}

func NewAuthenticatorHttpBearer(s *openapi3.SecuritySchemeRef, jwksUri string, issuer string, audience string) (*AuthenticatorHttpBearer, error) {
	if s.Value.BearerFormat != "JWT" {
		return nil, errors.New("bearer format must be jwt")
	}

	return &AuthenticatorHttpBearer{
		jwksUri:  jwksUri,
		issuer:   issuer,
		audience: audience,
	}, nil
}

func (a *AuthenticatorHttpBearer) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	ta := []string{}
	if a.audience != "" {
		ta = append(ta, a.audience)
	}

	c := JWTAuthenticatorConfig{
		JwksUrls:       []string{a.jwksUri},
		TrustedIssuers: []string{a.issuer},
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

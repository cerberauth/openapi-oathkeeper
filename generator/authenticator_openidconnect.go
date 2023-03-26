package generator

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorOpenIdConnect)(nil)

var httpClient = http.Client{
	Timeout: time.Second * 5,
}

type JWTAuthenticatorConfig struct {
	JwksUrls       []string `json:"jwks_urls"`
	TrustedIssuers []string `json:"trusted_issuers"`
	RequiredScope  []string `json:"required_scope"`
}

type OpenIdConfiguration struct {
	JwksUri string `json:"jwks_uri"`
	Issuer  string `json:"issuer"`
}

type AuthenticatorOpenIdConnect struct {
	C *OpenIdConfiguration
}

func NewAuthenticatorOpenIdConnect(s *openapi3.SecuritySchemeRef) (*AuthenticatorOpenIdConnect, error) {
	res, err := httpClient.Get(s.Value.OpenIdConnectUrl)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	c := OpenIdConfiguration{}
	jsonErr := json.Unmarshal(body, &c)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &AuthenticatorOpenIdConnect{C: &c}, nil
}

func (a *AuthenticatorOpenIdConnect) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	c := JWTAuthenticatorConfig{
		JwksUrls: []string{
			a.C.JwksUri,
		},
		TrustedIssuers: []string{
			a.C.Issuer,
		},
		RequiredScope: []string{},
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

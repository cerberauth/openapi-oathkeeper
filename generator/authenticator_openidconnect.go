package generator

import (
	"encoding/json"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorOpenIdConnect)(nil)

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
	return createRulesFromOAuth2SecurityRequirement(s, a.C.JwksUri, a.C.Issuer)
}

package authenticator

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

func TestNewHttpBearerAuthenticatorWhenBearerFormatIsNotJWT(t *testing.T) {
	_, newAuthenticatorErr := NewAuthenticatorHttpBearer(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewSecurityScheme(),
	}, "", "", "")

	require.Error(t, newAuthenticatorErr)
}

func TestHttpBearerAuthenticatorCreateAuthenticator(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://oauth.cerberauth.com"},
		RequiredScope:  []string{},
		TargetAudience: []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorHttpBearer(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, "https://oauth.cerberauth.com/.well-known/jwks.json", "https://oauth.cerberauth.com", "")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestHttpBearerAuthenticatorCreateAuthenticatorWithNonEmptyAudience(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://oauth.cerberauth.com"},
		RequiredScope:  []string{},
		TargetAudience: []string{"https://api.cerberauth.com"},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorHttpBearer(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, "https://oauth.cerberauth.com/.well-known/jwks.json", "https://oauth.cerberauth.com", "https://api.cerberauth.com")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

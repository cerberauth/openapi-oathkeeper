package authenticator

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestAuthenticatorDefaultCreateAuthenticator(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://ory.projects.oryapis.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
		"target_audience": []string{},
	})
	expectedAuthenticator := oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorDefault(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, &config.AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []string{"https://ory.projects.oryapis.com/.well-known/jwks.json"},
			"trusted_issuers": []string{"https://oauth.cerberauth.com"},
			"target_audience": []string{},
		},
	})
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, expectedAuthenticator, *auth)
}

func TestAuthenticatorDefaultCreateAuthenticatorWithScopes(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://ory.projects.oryapis.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{"resource:read", "resource:write"},
		"target_audience": []string{"https://api.cerberauth.com"},
	})

	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorDefault(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, &config.AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []string{"https://ory.projects.oryapis.com/.well-known/jwks.json"},
			"trusted_issuers": []string{"https://oauth.cerberauth.com"},
			"target_audience": []string{"https://api.cerberauth.com"},
		},
	})
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{
		"": []string{"resource:read", "resource:write"},
	})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, expectedAuthenticator, auth)
}

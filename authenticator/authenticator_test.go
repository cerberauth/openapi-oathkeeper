package authenticator

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

func TestNewAuthenticatorFromSecurityScheme(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	s := openapi3.NewJWTSecurityScheme()
	s.Extensions = make(map[string]interface{})
	s.Extensions["x-authenticator-jwks-uri"] = "https://oauth.cerberauth.com/.well-known/jwks.json"
	s.Extensions["x-authenticator-issuer"] = "https://oauth.cerberauth.com"
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(&openapi3.SecuritySchemeRef{
		Value: s,
	}, nil)
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestNewAuthenticatorFromSecuritySchemeWhenTypeIsOpenIDConnect(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://console.ory.sh/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://console.ory.sh"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, nil)
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestNewAuthenticatorFromSecuritySchemeWhenTypeIsOpenIDConnectWithConfig(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, &config.AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
			"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		},
	})
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestNewAuthenticatorFromSecuritySchemeWithConfiguration(t *testing.T) {
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, &config.AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
			"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		},
	})
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

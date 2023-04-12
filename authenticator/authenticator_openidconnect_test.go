package authenticator

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

func TestNewAuthenticatorOpenIdConnect(t *testing.T) {
	expectedAuthenticator := &AuthenticatorOpenIdConnect{
		C: &OpenIdConfiguration{
			JwksUri: "https://console.ory.sh/.well-known/jwks.json",
			Issuer:  "https://console.ory.sh",
		},
		Audience: "",
	}
	a, newAuthenticatorErr := NewAuthenticatorOpenIdConnect(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, "")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	assert.Equal(t, a, expectedAuthenticator)
}

func TestNewOpenIdConnectAuthenticatorCreateAuthenticator(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://console.ory.sh/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://console.ory.sh"},
		RequiredScope:  []string{},
		TargetAudience: []string{},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorOpenIdConnect(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, "")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestNewOpenIdConnectAuthenticatorCreateAuthenticatorWithNonEmptyAudience(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://console.ory.sh/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://console.ory.sh"},
		RequiredScope:  []string{},
		TargetAudience: []string{"https://api.cerberauth.com"},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorOpenIdConnect(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, "https://api.cerberauth.com")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

func TestNewOpenIdConnectAuthenticatorCreateAuthenticatorWithScopes(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://console.ory.sh/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://console.ory.sh"},
		RequiredScope:  []string{"resource:read", "resource:write"},
		TargetAudience: []string{"https://api.cerberauth.com"},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorOpenIdConnect(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme("https://project.console.ory.sh/.well-known/openid-configuration"),
	}, "https://api.cerberauth.com")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{
		"": []string{"resource:read", "resource:write"},
	})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

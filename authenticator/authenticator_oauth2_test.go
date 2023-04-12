package authenticator

import (
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

func TestNewOauth2Authenticator(t *testing.T) {
	expectedAuthenticator := &AuthenticatorOAuth2{
		JwksUri:  "",
		Issuer:   "",
		Audience: "",
	}
	a, newAuthenticatorErr := NewAuthenticatorOAuth2(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, "", "", "")
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	assert.Equal(t, a, expectedAuthenticator)
}

func TestNewOauth2AuthenticatorCreateAuthenticator(t *testing.T) {
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
	a, newAuthenticatorErr := NewAuthenticatorOAuth2(&openapi3.SecuritySchemeRef{
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

func TestNewOauth2AuthenticatorCreateAuthenticatorWithNonEmptyAudience(t *testing.T) {
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
	a, newAuthenticatorErr := NewAuthenticatorOAuth2(&openapi3.SecuritySchemeRef{
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

func TestNewOauth2AuthenticatorCreateAuthenticatorWithScopes(t *testing.T) {
	jsonConfig, _ := json.Marshal(JWTAuthenticatorConfig{
		JwksUrls:       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		TrustedIssuers: []string{"https://oauth.cerberauth.com"},
		RequiredScope:  []string{"resource:read", "resource:write"},
		TargetAudience: []string{"https://api.cerberauth.com"},
	})
	expectedAuthenticator := &rule.Handler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorOAuth2(&openapi3.SecuritySchemeRef{
		Value: openapi3.NewJWTSecurityScheme(),
	}, "https://oauth.cerberauth.com/.well-known/jwks.json", "https://oauth.cerberauth.com", "https://api.cerberauth.com")
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

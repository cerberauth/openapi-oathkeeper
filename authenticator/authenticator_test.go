package authenticator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jarcoal/httpmock"
)

var (
	oidcConfigurationUrl = "https://oauth.cerberauth.com/.well-known/openid-configuration"
	oidcConfiguration    = OpenIdConfiguration{
		Issuer:  "https://oauth.cerberauth.com",
		JwksUri: "https://oauth.cerberauth.com/.well-known/jwks.json",
	}
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	httpmock.Activate()
	resp, err := httpmock.NewJsonResponder(200, oidcConfiguration)
	if err != nil {
		tb.Fatal(err)
	}
	httpmock.RegisterResponder("GET", oidcConfigurationUrl, resp)

	return func(tb testing.TB) {
		defer httpmock.DeactivateAndReset()
	}
}

func TestNewAuthenticatorFromSecurityScheme(t *testing.T) {
	ctx := context.Background()
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	s := openapi3.NewJWTSecurityScheme()
	s.Extensions = make(map[string]interface{})
	s.Extensions["x-authenticator-jwks-uri"] = "https://oauth.cerberauth.com/.well-known/jwks.json"
	s.Extensions["x-authenticator-issuer"] = "https://oauth.cerberauth.com"
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(ctx, &openapi3.SecuritySchemeRef{
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(ctx, &openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme(oidcConfigurationUrl),
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

func TestNewAuthenticatorFromSecuritySchemeWhenTypeIsOpenIDConnectWithLowercaseType(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(ctx, &openapi3.SecuritySchemeRef{
		Value: &openapi3.SecurityScheme{
			Type:             "openidconnect",
			OpenIdConnectUrl: "https://oauth.cerberauth.com/.well-known/openid-configuration",
		},
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
	ctx := context.Background()
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(ctx, &openapi3.SecuritySchemeRef{
		Value: openapi3.NewOIDCSecurityScheme(oidcConfigurationUrl),
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
	ctx := context.Background()
	jsonConfig, _ := json.Marshal(map[string]interface{}{
		"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
		"trusted_issuers": []string{"https://oauth.cerberauth.com"},
		"required_scope":  []string{},
	})
	expectedAuthenticator := &oathkeeper.RuleHandler{
		Handler: "jwt",
		Config:  jsonConfig,
	}
	a, newAuthenticatorErr := NewAuthenticatorFromSecurityScheme(ctx, &openapi3.SecuritySchemeRef{
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

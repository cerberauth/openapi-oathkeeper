package authenticator

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

func TestNewNoopAuthenticator(t *testing.T) {
	expectedAuthenticator := &AuthenticatorNoop{}
	a, newAuthenticatorErr := NewAuthenticatorNoop()
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	assert.Equal(t, a, expectedAuthenticator)
}

func TestNoopAuthenticatorCreateAuthenticator(t *testing.T) {
	expectedAuthenticator := &rule.Handler{
		Handler: "noop",
	}
	a, newAuthenticatorErr := NewAuthenticatorNoop()
	if newAuthenticatorErr != nil {
		t.Fatal(newAuthenticatorErr)
	}

	auth, createAuthenticatorErr := a.CreateAuthenticator(&openapi3.SecurityRequirement{})
	if createAuthenticatorErr != nil {
		t.Fatal(createAuthenticatorErr)
	}

	assert.Equal(t, auth, expectedAuthenticator)
}

package authenticator

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestNoopAuthenticatorCreateAuthenticator(t *testing.T) {
	expectedAuthenticator := &oathkeeper.RuleHandler{
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

	assert.Equal(t, expectedAuthenticator, auth)
}

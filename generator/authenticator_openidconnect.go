package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorOpenIdConnect)(nil)

type AuthenticatorOpenIdConnect struct{}

func (*AuthenticatorOpenIdConnect) CreateAuthenticator(o *openapi3.Operation) *rule.Handler {
	return &rule.Handler{
		Handler: "jwt",
	}
}

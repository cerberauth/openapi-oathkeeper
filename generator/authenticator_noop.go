package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorNoop)(nil)

type AuthenticatorNoop struct{}

func (*AuthenticatorNoop) CreateAuthenticator(o *openapi3.Operation) *rule.Handler {
	return &rule.Handler{
		Handler: "noop",
	}
}

package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorNoop)(nil)

type AuthenticatorNoop struct{}

func (*AuthenticatorNoop) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	return &rule.Handler{
		Handler: "noop",
	}, nil
}

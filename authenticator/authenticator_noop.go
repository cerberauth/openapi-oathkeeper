package authenticator

import (
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
)

var _ Authenticator = (*AuthenticatorNoop)(nil)

type AuthenticatorNoop struct{}

func NewAuthenticatorNoop() (*AuthenticatorNoop, error) {
	return &AuthenticatorNoop{}, nil
}

func (*AuthenticatorNoop) CreateAuthenticator(s *openapi3.SecurityRequirement) (*oathkeeper.RuleHandler, error) {
	return &oathkeeper.RuleHandler{
		Handler: "noop",
	}, nil
}

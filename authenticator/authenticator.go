package authenticator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Authenticator interface {
	CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error)
}

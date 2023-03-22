package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
)

type Authenticator interface {
	CreateAuthenticator(o *openapi3.Operation) *rule.Handler
}

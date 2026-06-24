package config

import (
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
)

var defaultConfig = map[string]interface{}{
	"authorizers": []oathkeeper.RuleHandler{
		{
			Handler: "allow",
		},
	},
	"mutators": []oathkeeper.RuleHandler{
		{
			Handler: "noop",
		},
	},
	"errors": []oathkeeper.RuleErrorHandler{
		{
			Handler: handlerJSON,
		},
	},
	"upstream": oathkeeper.RuleUpstream{},
}

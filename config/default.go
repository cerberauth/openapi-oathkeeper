package config

import "github.com/ory/oathkeeper/rule"

var defaultConfig = map[string]interface{}{
	"authorizers": []rule.Handler{
		rule.Handler{
			Handler: "allow",
		},
	},
	"mutators": []rule.Handler{
		rule.Handler{
			Handler: "noop",
		},
	},
	"errors": []rule.ErrorHandler{
		rule.ErrorHandler{
			Handler: "json",
		},
	},
	"upstream": rule.Upstream{},
}

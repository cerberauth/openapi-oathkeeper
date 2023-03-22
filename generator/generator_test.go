package generator

import (
	"context"
	"encoding/json"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func newJWTConfig() json.RawMessage {
	config := JWTAuthenticatorConfig{
		jwks_urls: []string{
			"https://console.ory.sh/.well-known/jwks.json",
		},
		trusted_issuers: []string{
			"https://console.ory.sh",
		},
		required_scope: []string{
			"write:pets",
			"read:pets",
		},
	}
	jsonConfig, _ := json.Marshal(config)

	return jsonConfig
}

func TestGenerateFromSimpleOpenAPI(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []rule.Handler{
				{
					Handler: "noop",
				},
			},
			Authorizer: rule.Handler{
				Handler: "allow",
			},
			Mutators: []rule.Handler{
				{
					Handler: "noop",
				},
			},
		},
	}

	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, "../test/stub/simple.openapi.json"))
	require.NoError(t, err)

	ctx := context.Background()
	rules, err := New().Document(doc).Generate(ctx)
	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

// func TestGenerateFromSimpleOpenAPIWithOpenIdConnect(t *testing.T) {
// 	expectedRules := []rule.Rule{
// 		{
// 			ID:          "findPetsByStatus",
// 			Description: "Multiple status values can be provided with comma separated strings",
// 			Match: &rule.Match{
// 				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
// 				Methods: []string{"GET"},
// 			},
// 			Authenticators: []rule.Handler{
// 				{
// 					Handler: "jwt",
// 					Config:  newJWTConfig(),
// 				},
// 			},
// 			Authorizer: rule.Handler{
// 				Handler: "allow",
// 			},
// 			Mutators: []rule.Handler{
// 				{
// 					Handler: "noop",
// 				},
// 			},
// 		},
// 	}

// 	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, "../test/stub/simple_openidconnect.openapi.json"))
// 	require.NoError(t, err)

// 	ctx := context.Background()
// 	rules, err := New().Document(doc).Generate(ctx)
// 	require.NoError(t, err)
// 	assert.Equal(t, rules, expectedRules)
// }

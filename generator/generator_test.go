package generator

import (
	"context"
	"encoding/json"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func newJWTConfig(jwksUrls []string, issuers []string, scopes []string, audiences []string) json.RawMessage {
	c := authenticator.JWTAuthenticatorConfig{
		JwksUrls:       jwksUrls,
		TrustedIssuers: issuers,
		RequiredScope:  scopes,
		TargetAudience: audiences,
	}
	jsonConfig, _ := json.Marshal(c)

	return jsonConfig
}

func getRuleById(rules []rule.Rule, id string) *rule.Rule {
	for _, r := range rules {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

func newGenerator(docpath string, prefixId string, jwksUris map[string]string, allowedIssuers map[string]string, allowedAudiences map[string]string, serverUrls []string, upstreamUrl string, upstreamStripPath string) (*Generator, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, docpath))
	if err != nil {
		return nil, err
	}

	g := NewGenerator(prefixId, jwksUris, allowedIssuers, allowedAudiences, serverUrls, upstreamUrl, upstreamStripPath)

	ctx := context.Background()
	if loadErr := g.LoadOpenAPI3Doc(ctx, doc); loadErr != nil {
		return nil, loadErr
	}

	return g, nil
}

func TestGenerateFromSimpleOpenAPI(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "", nil, nil, nil, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithPrefixId(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "prefix:findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "prefix", nil, nil, nil, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithOneServerUrl(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://www\\.cerberauth\\.com/api)(/pet/findByStatus/?)$>",
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	serverUrls := []string{"https://www.cerberauth.com/api"}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "", nil, nil, nil, serverUrls, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithSeveralServerUrls(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://www\\.cerberauth\\.com/api|https://api\\.cerberauth\\.com/api)(/pet/findByStatus/?)$>",
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	serverUrls := []string{
		"https://www.cerberauth.com/api",
		"https://api.cerberauth.com/api",
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "", nil, nil, nil, serverUrls, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnect(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
				Methods: []string{"GET"},
			},
			Authenticators: []rule.Handler{
				{
					Handler: "jwt",
					Config: newJWTConfig([]string{
						"https://console.ory.sh/.well-known/jwks.json",
					}, []string{
						"https://console.ory.sh",
					}, []string{
						"write:pets",
						"read:pets",
					}, []string{}),
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect.openapi.json", "", nil, nil, nil, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithOAuth2(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
				Methods: []string{"GET"},
			},
			Authenticators: []rule.Handler{
				{
					Handler: "jwt",
					Config: newJWTConfig([]string{
						"https://oauth.cerberauth.com/.well-known/jwks.json",
					}, []string{
						"https://cerberauth.com",
					}, []string{
						"write:pets",
						"read:pets",
					}, []string{
						"https://api.cerberauth.com",
					}),
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_oauth2.openapi.json", "", map[string]string{
		"petstore_auth": "https://oauth.cerberauth.com/.well-known/jwks.json",
	}, map[string]string{
		"petstore_auth": "https://cerberauth.com",
	}, map[string]string{
		"petstore_auth": "https://api.cerberauth.com",
	}, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithHttpBearer(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
				Methods: []string{"GET"},
			},
			Authenticators: []rule.Handler{
				{
					Handler: "jwt",
					Config: newJWTConfig([]string{
						"https://oauth.cerberauth.com/.well-known/jwks.json",
					}, []string{
						"https://cerberauth.com",
					}, []string{}, []string{
						"https://api.cerberauth.com",
					}),
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_http_bearer_jwt.openapi.json", "", map[string]string{
		"petstore_auth": "https://oauth.cerberauth.com/.well-known/jwks.json",
	}, map[string]string{
		"petstore_auth": "https://cerberauth.com",
	}, map[string]string{
		"petstore_auth": "https://api.cerberauth.com",
	}, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnectWithGlobalSecurityScheme(t *testing.T) {
	expectedRule := rule.Rule{
		ID:          "updatePet",
		Description: "Update an existing pet by Id",
		Match: &rule.Match{
			URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/?)$>",
			Methods: []string{"PUT"},
		},
		Authenticators: []rule.Handler{
			{
				Handler: "jwt",
				Config: newJWTConfig([]string{
					"https://console.ory.sh/.well-known/jwks.json",
				}, []string{
					"https://console.ory.sh",
				}, []string{
					"write:pets",
					"read:pets",
				}, []string{}),
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
		Errors: []rule.ErrorHandler{
			{
				Handler: "json",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", "", nil, nil, nil, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, *getRuleById(rules, "updatePet"), expectedRule)
}

func TestGenerateFromSimpleOpenAPIWithUpstreamUrlAndPath(t *testing.T) {
	expectedRules := []rule.Rule{
		{
			ID:          "prefix:findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &rule.Match{
				URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?)$>",
				Methods: []string{"GET"},
			},
			Upstream: rule.Upstream{
				URL:       "https://petstore.com",
				StripPath: "/api",
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
			Errors: []rule.ErrorHandler{
				{
					Handler: "json",
				},
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "prefix", nil, nil, nil, nil, "https://petstore.com", "/api")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnectWithGlobalAndLocalOverrideSecurityScheme(t *testing.T) {
	expectedRule := rule.Rule{
		ID:          "findPetsByStatus",
		Description: "Multiple status values can be provided with comma separated strings",
		Match: &rule.Match{
			URL:     "<^(https://petstore\\.swagger\\.io/api/v3)(/pet/findByStatus/?(\\?.+)?)$>",
			Methods: []string{"GET"},
		},
		Authenticators: []rule.Handler{
			{
				Handler: "jwt",
				Config: newJWTConfig([]string{
					"https://console.ory.sh/.well-known/jwks.json",
				}, []string{
					"https://console.ory.sh",
				}, []string{
					"read:pets",
				}, []string{}),
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
		Errors: []rule.ErrorHandler{
			{
				Handler: "json",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", "", nil, nil, nil, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, *getRuleById(rules, "findPetsByStatus"), expectedRule)
}

func TestGenerateFromPetstoreWithOpenIdConnect(t *testing.T) {
	g, newGeneratorErr := newGenerator("../test/stub/petstore.openapi.json", "", map[string]string{
		"petstore_auth": "https://oauth.cerberauth.com/.well-known/jwks.json",
	}, map[string]string{
		"petstore_auth": "https://cerberauth.com",
	}, map[string]string{
		"petstore_auth": "https://api.cerberauth.com",
	}, nil, "", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].GetID() < rules[j].GetID() })

	require.NoError(t, err)
	cupaloy.SnapshotT(t, rules)
}

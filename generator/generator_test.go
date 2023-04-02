package generator

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func newJWTConfig(scopes []string) json.RawMessage {
	c := JWTAuthenticatorConfig{
		JwksUrls: []string{
			"https://console.ory.sh/.well-known/jwks.json",
		},
		TrustedIssuers: []string{
			"https://console.ory.sh",
		},
		RequiredScope: scopes,
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

func newGenerator(docpath string, prefixId string) (*Generator, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, docpath))
	if err != nil {
		return nil, err
	}

	g := NewGenerator(prefixId)

	ctx := context.Background()
	if loadErr := g.LoadOpenAPI3Doc(ctx, doc); loadErr != nil {
		return nil, errors.New("an error occurred loading the openapi doc")
	}

	return g, nil
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
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "")
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
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", "prefix")
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
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []rule.Handler{
				{
					Handler: "jwt",
					Config: newJWTConfig([]string{
						"write:pets",
						"read:pets",
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
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect.openapi.json", "")
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
			URL:     "https://petstore.swagger.io/api/v3/pet",
			Methods: []string{"PUT"},
		},
		Authenticators: []rule.Handler{
			{
				Handler: "jwt",
				Config: newJWTConfig([]string{
					"write:pets",
					"read:pets",
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
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, *getRuleById(rules, "updatePet"), expectedRule)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnectWithGlobalAndLocalOverrideSecurityScheme(t *testing.T) {
	expectedRule := rule.Rule{
		ID:          "findPetsByStatus",
		Description: "Multiple status values can be provided with comma separated strings",
		Match: &rule.Match{
			URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
			Methods: []string{"GET"},
		},
		Authenticators: []rule.Handler{
			{
				Handler: "jwt",
				Config: newJWTConfig([]string{
					"read:pets",
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
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, *getRuleById(rules, "findPetsByStatus"), expectedRule)
}

func TestGenerateFromPetstoreWithOpenIdConnect(t *testing.T) {
	g, newGeneratorErr := newGenerator("../test/stub/petstore_openidconnect.openapi.json", "")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].GetID() < rules[j].GetID() })

	require.NoError(t, err)
	cupaloy.SnapshotT(t, rules)
}

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

func newJWTConfig() json.RawMessage {
	c := JWTAuthenticatorConfig{
		JwksUrls: []string{
			"https://console.ory.sh/.well-known/jwks.json",
		},
		TrustedIssuers: []string{
			"https://console.ory.sh",
		},
		RequiredScope: []string{
			"write:pets",
			"read:pets",
		},
	}
	jsonConfig, _ := json.Marshal(c)

	return jsonConfig
}

func newGenerator(docpath string) (*Generator, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, docpath))
	if err != nil {
		return nil, err
	}

	g := NewGenerator()

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
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json")
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
					Config:  newJWTConfig(),
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
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect.openapi.json")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()

	require.NoError(t, err)
	assert.Equal(t, rules, expectedRules)
}

func TestGenerateFromPetstoreWithOpenIdConnect(t *testing.T) {
	g, newGeneratorErr := newGenerator("../test/stub/petstore_openidconnect.openapi.json")
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate()
	sort.SliceStable(rules, func(i, j int) bool { return rules[i].GetID() < rules[j].GetID() })

	require.NoError(t, err)
	cupaloy.SnapshotT(t, rules)
}

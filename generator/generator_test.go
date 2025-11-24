package generator

import (
	"context"
	"encoding/json"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _           = runtime.Caller(0)
	basepath             = filepath.Dir(b)
	oidcConfigurationUrl = "https://oauth.cerberauth.com/.well-known/openid-configuration"
	oidcConfiguration    = authenticator.OpenIdConfiguration{
		Issuer:  "https://oauth.cerberauth.com",
		JwksUri: "https://oauth.cerberauth.com/.well-known/jwks.json",
	}
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	httpmock.Activate()
	resp, err := httpmock.NewJsonResponder(200, oidcConfiguration)
	if err != nil {
		tb.Fatal(err)
	}
	httpmock.RegisterResponder("GET", oidcConfigurationUrl, resp)

	return func(tb testing.TB) {
		defer httpmock.DeactivateAndReset()
	}
}

func getRuleById(rules []oathkeeper.Rule, id string) *oathkeeper.Rule {
	for _, r := range rules {
		if r.ID == id {
			return &r
		}
	}

	return nil
}

func newGenerator(docpath string, cfg *config.Config) (*Generator, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(path.Join(basepath, docpath))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return NewGenerator(ctx, doc, cfg)
}

func TestGenerateFromSimpleOpenAPI(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithPrefixId(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "prefix:findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", &config.Config{
		Prefix: "prefix",
	})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithOneServerUrl(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://www.cerberauth.com/api/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", &config.Config{
		ServerUrls: []string{"https://www.cerberauth.com/api"},
	})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithSeveralServerUrls(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "<(https://www\\.cerberauth\\.com/api|https://api\\.cerberauth\\.com/api)>/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", &config.Config{
		ServerUrls: []string{
			"https://www.cerberauth.com/api",
			"https://api.cerberauth.com/api",
		},
	})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateOpenAPIWithoutSecurity(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "withEmptySecurity",
			Description: "",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://api.cerberauth.com/v1/withEmptySecurity",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
					Config:  nil,
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},

		{
			ID:          "withSecurity",
			Description: "",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://api.cerberauth.com/v1/withSecurity",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
					Config:  nil,
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, err := newGenerator("../test/stub/simple_no_security.openapi.json", &config.Config{})
	if err != nil {
		t.Fatal(err)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnect(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	c, _ := json.Marshal(map[string]interface{}{
		"jwks_urls": []string{
			"https://oauth.cerberauth.com/.well-known/jwks.json",
		},
		"trusted_issuers": []string{
			"https://oauth.cerberauth.com",
		},
		"required_scope": []string{
			"write:pets",
			"read:pets",
		},
	})
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "jwt",
					Config:  c,
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithOAuth2(t *testing.T) {
	ctx := context.Background()
	c, _ := json.Marshal(map[string]interface{}{
		"jwks_urls": []string{
			"https://oauth.cerberauth.com/.well-known/jwks.json",
		},
		"trusted_issuers": []string{
			"https://oauth.cerberauth.com",
		},
		"required_scope": []string{
			"write:pets",
			"read:pets",
		},
		"target_audience": []string{
			"https://api.cerberauth.com",
		},
	})
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "jwt",
					Config:  c,
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_oauth2.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithHttpBearer(t *testing.T) {
	ctx := context.Background()
	c, _ := json.Marshal(map[string]interface{}{
		"jwks_urls": []string{
			"https://oauth.cerberauth.com/.well-known/jwks.json",
		},
		"trusted_issuers": []string{
			"https://oauth.cerberauth.com",
		},
		"required_scope": []string{},
		"target_audience": []string{
			"https://api.cerberauth.com",
		},
	})
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "jwt",
					Config:  c,
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_http_bearer_jwt.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnectWithGlobalSecurityScheme(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	c, _ := json.Marshal(map[string]interface{}{
		"jwks_urls": []string{
			"https://oauth.cerberauth.com/.well-known/jwks.json",
		},
		"trusted_issuers": []string{
			"https://oauth.cerberauth.com",
		},
		"required_scope": []string{
			"write:pets",
			"read:pets",
		},
	})
	expectedRule := oathkeeper.Rule{
		ID:          "updatePet",
		Description: "Update an existing pet by Id",
		Match: &oathkeeper.RuleMatch{
			URL:     "https://petstore.swagger.io/api/v3/pet",
			Methods: []string{"PUT"},
		},
		Authenticators: []oathkeeper.RuleHandler{
			{
				Handler: "jwt",
				Config:  c,
			},
		},
		Authorizer: oathkeeper.RuleHandler{
			Handler: "allow",
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRule, *getRuleById(rules, "updatePet"))
}

func TestGenerateFromSimpleOpenAPIWithUpstreamUrlAndPath(t *testing.T) {
	ctx := context.Background()
	expectedRules := []oathkeeper.Rule{
		{
			ID:          "prefix:findPetsByStatus",
			Description: "Multiple status values can be provided with comma separated strings",
			Match: &oathkeeper.RuleMatch{
				URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
				Methods: []string{"GET"},
			},
			Upstream: oathkeeper.RuleUpstream{
				URL:       "https://petstore.com",
				StripPath: "/api",
			},
			Authenticators: []oathkeeper.RuleHandler{
				{
					Handler: "noop",
				},
			},
			Authorizer: oathkeeper.RuleHandler{
				Handler: "allow",
			},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple.openapi.json", &config.Config{
		Prefix: "prefix",
		Upstream: oathkeeper.RuleUpstream{
			URL:       "https://petstore.com",
			StripPath: "/api",
		},
	})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func TestGenerateFromSimpleOpenAPIWithOpenIdConnectWithGlobalAndLocalOverrideSecurityScheme(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	c, _ := json.Marshal(map[string]interface{}{
		"jwks_urls": []string{
			"https://oauth.cerberauth.com/.well-known/jwks.json",
		},
		"trusted_issuers": []string{
			"https://oauth.cerberauth.com",
		},
		"required_scope": []string{
			"read:pets",
		},
	})
	expectedRule := oathkeeper.Rule{
		ID:          "findPetsByStatus",
		Description: "Multiple status values can be provided with comma separated strings",
		Match: &oathkeeper.RuleMatch{
			URL:     "https://petstore.swagger.io/api/v3/pet/findByStatus",
			Methods: []string{"GET"},
		},
		Authenticators: []oathkeeper.RuleHandler{
			{
				Handler: "jwt",
				Config:  c,
			},
		},
		Authorizer: oathkeeper.RuleHandler{
			Handler: "allow",
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/simple_openidconnect_global.openapi.json", &config.Config{})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	got := *getRuleById(rules, "findPetsByStatus")
	assert.Equal(t, expectedRule, got)
}

func TestGenerateFromPetstoreWithOpenIdConnect(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	ctx := context.Background()
	var authenticators = make(map[string]config.AuthenticatorRuleConfig)
	authenticators["petstore_auth"] = config.AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []string{"https://oauth.cerberauth.com/.well-known/jwks.json"},
			"trusted_issuers": []string{"https://oauth.cerberauth.com"},
			"target_audience": []string{"https://api.cerberauth.com"},
		},
	}
	g, newGeneratorErr := newGenerator("../test/stub/petstore.openapi.json", &config.Config{
		Authenticators: authenticators,
	})
	if newGeneratorErr != nil {
		t.Fatal(newGeneratorErr)
	}

	rules, err := g.Generate(ctx)

	require.NoError(t, err)
	cupaloy.SnapshotT(t, rules)
}

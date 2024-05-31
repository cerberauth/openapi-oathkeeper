package generator

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestGenerateMatchRule(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/api/v3",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3"}, "GET", "/", nil)

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithServerUrlEndingSlash(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/api/v3",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3/"}, "GET", "/", nil)

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithMultipleServerUrls(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "<(/api/v3|https://cerberauth\\.com/api/v3)>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3", "https://cerberauth.com/api/v3"}, "GET", "/", nil)

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithNoPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/<.+>/resource/<.+>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{}, "GET", "/{param}/resource/{otherParam}", nil)

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithUnknownPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/<.+>/resource/<.+>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{}, "GET", "/{param}/resource/{otherParam}", &openapi3.Parameters{})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithStringPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<.+>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"string"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithNumberPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<((\\x2D|\\+)?\\d+(?:\\.\\d+)?)>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"number"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithIntegerPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<\\d+>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"integer"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithBooleanPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<.+>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"boolean"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithArrayPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<.+>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"array"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

func TestGenerateMatchRuleWithObjectPathParams(t *testing.T) {
	expectedMatchingRule := oathkeeper.RuleMatch{
		URL:     "/resource/<.+>",
		Methods: []string{"GET"},
	}
	types := &openapi3.Types{"object"}
	param := openapi3.NewPathParameter("param").WithSchema(&openapi3.Schema{Type: types})
	matchRule, err := createMatchRule([]string{}, "GET", "/resource/{param}", &openapi3.Parameters{&openapi3.ParameterRef{
		Ref:   "param",
		Value: param,
	}})

	require.NoError(t, err)
	assert.Equal(t, &expectedMatchingRule, matchRule)
}

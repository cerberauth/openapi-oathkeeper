package generator

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

func TestGenerateMatchRule(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3)(/)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3"}, "GET", "/", nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWhenThereIsNoServerUrl(t *testing.T) {
	_, err := createMatchRule([]string{}, "GET", "/", nil)

	require.Error(t, err)
}

func TestGenerateMatchRuleWithMultipleServerUrls(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3|https://cerberauth\\.com/api/v3)(/)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3", "https://cerberauth.com/api/v3"}, "GET", "/", nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithPathParams(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(https://cerberauth\\.com/api/v3)(/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/resource/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"https://cerberauth.com/api/v3"}, "GET", "/{param}/resource/{otherParam}", nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithQueryParams(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(https://cerberauth\\.com/api/v3)(/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/resource/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?(\\?.+)?)$>",
		Methods: []string{"GET"},
	}
	params := openapi3.NewParameters()
	params = append(params, &openapi3.ParameterRef{
		Value: openapi3.NewQueryParameter("test"),
	})
	matchRule, err := createMatchRule([]string{"https://cerberauth.com/api/v3"}, "GET", "/{param}/resource/{otherParam}", &params)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

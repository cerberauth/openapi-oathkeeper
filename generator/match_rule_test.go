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
	matchRule, err := createMatchRule([]string{"/api/v3"}, []string{"GET"}, []string{"/"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWhenThereIsNoServerUrl(t *testing.T) {
	_, err := createMatchRule([]string{}, []string{"GET"}, []string{"/"}, nil)

	require.Error(t, err)
}

func TestGenerateMatchRuleWithMultipleServerUrls(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3|https://cerberauth\\.com/api/v3)(/)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3", "https://cerberauth.com/api/v3"}, []string{"GET"}, []string{"/"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithPathParams(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(https://cerberauth\\.com/api/v3)(/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/resource/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"https://cerberauth.com/api/v3"}, []string{"GET"}, []string{"/{param}/resource/{otherParam}"}, nil)

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
	matchRule, err := createMatchRule([]string{"https://cerberauth.com/api/v3"}, []string{"GET"}, []string{"/{param}/resource/{otherParam}"}, &params)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithMultipleMethods(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3)(/)$>",
		Methods: []string{"GET", "POST", "PUT"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3"}, []string{"GET", "POST", "PUT"}, []string{"/"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithMultipleTimesTheSamePath(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3)(/)$>",
		Methods: []string{"GET", "POST", "PUT"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3", "/api/v3"}, []string{"GET", "POST", "PUT"}, []string{"/"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithMultiplePaths(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/[v2|v3])(/)$>",
		Methods: []string{"GET", "POST", "PUT"},
	}
	matchRule, err := createMatchRule([]string{"/api/v2", "/api/v3"}, []string{"GET", "POST", "PUT"}, []string{"/"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithMultiplePathsAndSamePathParams(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(https://cerberauth\\.com/api/[v2|v3])(/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/resource/(?:[[:alnum:]]?\\x2D?=?\\??&?)+/?(\\?.+)?)$>",
		Methods: []string{"GET", "POST", "PUT"},
	}
	matchRule, err := createMatchRule([]string{"/api/v2", "/api/v3"}, []string{"GET", "POST", "PUT"}, []string{"/{param}/resource/{otherParam}"}, nil)

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

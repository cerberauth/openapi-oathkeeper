package generator

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/require"
)

func TestGenerateMatchRule(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3)(/)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3"}, "GET", "/")

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWhenThereIsNoServerUrl(t *testing.T) {
	_, err := createMatchRule([]string{}, "GET", "/")

	require.Error(t, err)
}

func TestGenerateMatchRuleWithMultipleServerUrls(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(/api/v3|https://cerberauth\\.com/api/v3)(/)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"/api/v3", "https://cerberauth.com/api/v3"}, "GET", "/")

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

func TestGenerateMatchRuleWithParams(t *testing.T) {
	expectedMatchingRule := rule.Match{
		URL:     "<^(https://cerberauth\\.com/api/v3)(/(.+)/resource/(.+)/?)$>",
		Methods: []string{"GET"},
	}
	matchRule, err := createMatchRule([]string{"https://cerberauth.com/api/v3"}, "GET", "/{param}/resource/{otherParam}")

	require.NoError(t, err)
	assert.Equal(t, matchRule, &expectedMatchingRule)
}

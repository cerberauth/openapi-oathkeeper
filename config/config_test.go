package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ory/oathkeeper/rule"
	"github.com/stretchr/testify/assert"
)

func TestLoadWithConfigFile(t *testing.T) {
	config := `
prefix: "prefix:"

server_urls:
- "https://www.cerberauth.com/api"
- "https://api.cerberauth.com/api"

upstream:
  url: "https://api.cerberauth.com"
`
	tempFile, err := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, err)
	fmt.Println("Create temp file::", tempFile.Name())
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(config)
	assert.NoError(t, err)

	// when
	cfg, err := New(tempFile.Name())

	// then
	assert.NoError(t, err)
	assert.Equal(t, "prefix:", cfg.Prefix)
	assert.Equal(t, []string{"https://www.cerberauth.com/api", "https://api.cerberauth.com/api"}, cfg.ServerUrls)
	assert.Equal(t, "https://api.cerberauth.com", cfg.Upstream.URL)
}

func TestLoadWithConfigFileWithDefaultConfig(t *testing.T) {
	config := `
prefix: "prefix:"

server_urls:
- "https://www.cerberauth.com/api"
- "https://api.cerberauth.com/api"
`
	tempFile, err := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, err)
	fmt.Println("Create temp file::", tempFile.Name())
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(config)
	assert.NoError(t, err)

	// when
	cfg, err := New(tempFile.Name())

	// then
	assert.NoError(t, err)
	assert.Equal(t, []rule.Handler{
		rule.Handler{
			Handler: "noop",
		},
	}, cfg.Mutators)
	assert.Equal(t, []rule.ErrorHandler{
		rule.ErrorHandler{
			Handler: "json",
		},
	}, cfg.Errors)
	assert.Equal(t, rule.Upstream{}, cfg.Upstream)
}

func TestLoadWithConfigFileWithRules(t *testing.T) {
	config := `
prefix: "prefix:"

server_urls:
- "https://www.cerberauth.com/api"
- "https://api.cerberauth.com/api"

authenticators:
  oauth:
    handler: "jwt"
    config:
      jwks_urls:
      - https://console.ory.sh/.well-known/jwks.json
      trusted_issuers:
      - https://console.ory.sh
      target_audience:
      - https://api.cerberauth.com
`
	tempFile, err := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, err)
	fmt.Println("Create temp file:", tempFile.Name())
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(config)
	assert.NoError(t, err)

	// when
	cfg, err := New(tempFile.Name())

	// then
	assert.NoError(t, err)
	assert.Equal(t, AuthenticatorRuleConfig{
		Handler: "jwt",
		Config: map[string]interface{}{
			"jwks_urls":       []interface{}{"https://console.ory.sh/.well-known/jwks.json"},
			"trusted_issuers": []interface{}{"https://console.ory.sh"},
			"target_audience": []interface{}{"https://api.cerberauth.com"},
		},
	}, cfg.Authenticators["oauth"])
}

func TestLoadWithConfigFileWithInvalidFile(t *testing.T) {
	config := `foo`
	tempFile, err := ioutil.TempFile(os.TempDir(), "config-test")
	assert.NoError(t, err)
	fmt.Println("Create temp file:", tempFile.Name())
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(config)
	assert.NoError(t, err)

	// when
	_, err = New(tempFile.Name())

	// then
	assert.Error(t, err)
}

func TestLoadWithConfigFileWithInvalidFilePath(t *testing.T) {
	tempFilePath := filepath.Join(os.TempDir(), "config-test")

	// when
	_, err := New(tempFilePath)

	// then
	assert.Error(t, err)
}

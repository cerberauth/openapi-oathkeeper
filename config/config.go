package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Prefix     string                  `json:"prefix" yaml:"prefix"`
	ServerUrls []string                `json:"server_urls" yaml:"server_urls"`
	Upstream   oathkeeper.RuleUpstream `json:"upstream" yaml:"upstream"`

	Authenticators map[string]AuthenticatorRuleConfig `json:"authenticators" yaml:"authenticators"`
	Authorizer     oathkeeper.RuleHandler             `json:"authorizer" yaml:"authorizer"`
	Mutators       []oathkeeper.RuleHandler           `json:"mutators" yaml:"mutators"`
	Errors         []oathkeeper.RuleErrorHandler      `json:"errors" yaml:"errors"`
}

type AuthenticatorRuleConfig struct {
	Handler string                 `json:"handler" yaml:"handler"`
	Config  map[string]interface{} `json:"config" yaml:"handler"`
}

var k = koanf.New(".")

func New(configFilePath string) (*Config, error) {
	if _, err := os.Stat(configFilePath); err != nil {
		log.Printf("the configuration file has not been found on %s", configFilePath)

		return nil, err
	}

	// load from default config
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		log.Printf("error loading default config: %v", err)
	}

	// load from config file if exist
	if configFilePath != "" {
		path, err := filepath.Abs(configFilePath)
		if err != nil {
			log.Printf("failed to get absolute config path for %s: %v", configFilePath, err)
			return nil, err
		}
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("error loading config: %v", err)
			return nil, err
		}
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf: %v", err)
		return nil, err
	}
	return &cfg, err
}

package authenticator

import (
	"encoding/json"

	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/knadh/koanf/maps"
	"github.com/ory/oathkeeper/rule"
)

var _ Authenticator = (*AuthenticatorDefault)(nil)

type AuthenticatorDefault struct {
	cfg *config.AuthenticatorRuleConfig
}

func NewAuthenticatorDefault(s *openapi3.SecuritySchemeRef, cfg *config.AuthenticatorRuleConfig) (*AuthenticatorDefault, error) {
	return &AuthenticatorDefault{
		cfg: cfg,
	}, nil
}

func (a *AuthenticatorDefault) CreateAuthenticator(s *openapi3.SecurityRequirement) (*rule.Handler, error) {
	required_scope := make([]string, 0)
	for _, scope := range *s {
		required_scope = append(required_scope, scope...)
	}

	cfg := maps.Copy(a.cfg.Config)
	cfg["required_scope"] = required_scope

	jsonConfig, jsonErr := json.Marshal(cfg)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &rule.Handler{
		Handler: a.cfg.Handler,
		Config:  jsonConfig,
	}, nil
}

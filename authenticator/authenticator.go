package authenticator

import (
	"context"
	"errors"
	"strings"

	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/cerberauth/x/telemetryx"
	"github.com/getkin/kin-openapi/openapi3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Authenticator interface {
	CreateAuthenticator(s *openapi3.SecurityRequirement) (*oathkeeper.RuleHandler, error)
}

var otelName = "github.com/cerberauth/openapi-oathkeeper/authenticator"

var (
	JWKSUriExtensionName  = "x-authenticator-jwks-uri"
	IssuerExtensionName   = "x-authenticator-issuer"
	AudienceExtensionName = "x-authenticator-audience"
)

func getFromExtension(s *openapi3.SecuritySchemeRef, name string) *string {
	v, e := s.Value.Extensions[name]
	if !e {
		return nil
	}

	ext := v.(string)
	return &ext
}

func createConfigFromSecurityScheme(s *openapi3.SecuritySchemeRef) (*config.AuthenticatorRuleConfig, error) {
	cfg := config.AuthenticatorRuleConfig{
		Config: make(map[string]interface{}),
	}
	switch strings.ToLower(s.Value.Type) {
	case string(AuthenticatorTypeOpenIdConnect),
		string(AuthenticatorTypeOAuth2):
		cfg.Handler = "jwt"

	case string(AuthenticatorTypeHttp):
		if s.Value.Scheme == "bearer" {
			cfg.Handler = "jwt"
		}
	}

	if cfg.Handler == "" {
		return nil, errors.New("can not detect rule handler for the security scheme")
	}

	return &cfg, nil
}

func NewAuthenticatorFromSecurityScheme(ctx context.Context, s *openapi3.SecuritySchemeRef, cfg *config.AuthenticatorRuleConfig) (Authenticator, error) {
	telemetryMeter := telemetryx.GetMeterProvider().Meter(otelName)
	telemetryAuthenticatorCreatedCounter, _ := telemetryMeter.Int64Counter(
		"authenticator.created.counter",
		metric.WithDescription("Number of authenticator created"),
		metric.WithUnit("{operation}"),
	)

	if cfg == nil {
		defaultCfg, defaultCfgErr := createConfigFromSecurityScheme(s)
		if defaultCfgErr != nil {
			return nil, defaultCfgErr
		}

		cfg = defaultCfg
	}

	if jwksUri := getFromExtension(s, JWKSUriExtensionName); jwksUri != nil {
		cfg.Config["jwks_urls"] = []string{*jwksUri}
	}

	if trusted_issuer := getFromExtension(s, IssuerExtensionName); trusted_issuer != nil {
		cfg.Config["trusted_issuers"] = []string{*trusted_issuer}
	}

	if allowedAudience := getFromExtension(s, AudienceExtensionName); allowedAudience != nil {
		cfg.Config["target_audience"] = []string{*allowedAudience}
	}

	if strings.ToLower(s.Value.Type) == string(AuthenticatorTypeOpenIdConnect) && (cfg.Config["jwks_urls"] == nil || cfg.Config["trusted_issuers"] == nil) {
		c, err := fetchOpenIDConfiguration(s.Value.OpenIdConnectUrl)
		if err != nil {
			return nil, err
		}

		cfg.Config["jwks_urls"] = []string{c.JwksUri}
		cfg.Config["trusted_issuers"] = []string{c.Issuer}
	}

	telemetryAuthenticatorCreatedCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("type", s.Value.Type), attribute.String("handler", s.Value.Scheme)))

	return NewAuthenticatorDefault(s, cfg)
}

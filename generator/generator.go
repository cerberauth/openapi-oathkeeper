package generator

import (
	"context"
	"sort"

	"github.com/cerberauth/openapi-oathkeeper/authenticator"
	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/cerberauth/x/telemetryx"
	"github.com/getkin/kin-openapi/openapi3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var otelName = "github.com/cerberauth/openapi-oathkeeper/generator"

type Generator struct {
	doc *openapi3.T
	cfg *config.Config

	authenticators map[string]authenticator.Authenticator
}

type RulesById []oathkeeper.Rule

func (r RulesById) Len() int           { return len(r) }
func (r RulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RulesById) Less(i, j int) bool { return r[i].GetID() < r[j].GetID() }

func (g *Generator) computeId(operationId string) string {
	if g.cfg.Prefix == "" {
		return operationId
	}

	return g.cfg.Prefix + ":" + operationId
}

func (g *Generator) createRule(verb string, path string, o *openapi3.Operation) (*oathkeeper.Rule, error) {
	match, matchRuleErr := createMatchRule(g.cfg.ServerUrls, verb, path, &o.Parameters)
	if matchRuleErr != nil {
		return nil, matchRuleErr
	}

	var authenticators = []oathkeeper.RuleHandler{}
	appendAuthenticator := func(sr *openapi3.SecurityRequirements) error {
		for _, s := range *sr {
			for k := range s {
				if a, ok := g.authenticators[k]; ok {
					ar, arerror := a.CreateAuthenticator(&s)
					if arerror != nil {
						return arerror
					}
					authenticators = append(authenticators, *ar)
				}
			}
		}

		return nil
	}

	var err error
	// nolint: gocritic
	if o.Security != nil && len(*o.Security) > 0 {
		err = appendAuthenticator(o.Security)
	} else if len(g.doc.Security) > 0 {
		err = appendAuthenticator(&g.doc.Security)
	} else {
		ar, arerror := g.authenticators[string(authenticator.AuthenticatorTypeNoop)].CreateAuthenticator(nil)
		if arerror != nil {
			return nil, arerror
		}

		authenticators = append(authenticators, *ar)
	}

	if err != nil {
		return nil, err
	}

	return &oathkeeper.Rule{
		ID:             g.computeId(o.OperationID),
		Description:    o.Description,
		Match:          match,
		Upstream:       g.cfg.Upstream,
		Authenticators: authenticators,
		Authorizer: oathkeeper.RuleHandler{
			Handler: "allow",
		},
		Mutators: g.cfg.Mutators,
		Errors:   g.cfg.Errors,
	}, nil
}

func createAuthenticators(ctx context.Context, d *openapi3.T, cfg *config.Config) (map[string]authenticator.Authenticator, error) {
	authenticators := make(map[string]authenticator.Authenticator)

	// Create a first authenticator for operations without security configured
	authenticators[string(authenticator.AuthenticatorTypeNoop)] = &authenticator.AuthenticatorNoop{}

	newAuthenticator := func(name string) (authenticator.Authenticator, error) {
		s := d.Components.SecuritySchemes[name]
		v, ok := cfg.Authenticators[name]
		if !ok {
			return authenticator.NewAuthenticatorFromSecurityScheme(ctx, s, nil)
		}

		return authenticator.NewAuthenticatorFromSecurityScheme(ctx, s, &v)
	}

	if d.Components != nil && d.Components.SecuritySchemes != nil {
		for name := range d.Components.SecuritySchemes {
			a, err := newAuthenticator(name)
			if err != nil {
				return nil, err
			}
			authenticators[name] = a
		}
	}

	return authenticators, nil
}

func NewGenerator(ctx context.Context, d *openapi3.T, cfg *config.Config) (*Generator, error) {
	if validateErr := d.Validate(ctx, openapi3.DisableExamplesValidation(), openapi3.DisableSchemaDefaultsValidation()); validateErr != nil {
		return nil, validateErr
	}

	if cfg.ServerUrls == nil {
		for _, s := range d.Servers {
			cfg.ServerUrls = append(cfg.ServerUrls, s.URL)
		}
	}

	authenticators, err := createAuthenticators(ctx, d, cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{
		doc: d,
		cfg: cfg,

		authenticators: authenticators,
	}, nil
}

func (g *Generator) Generate(ctx context.Context) ([]oathkeeper.Rule, error) {
	telemetryMeter := telemetryx.GetMeterProvider().Meter(otelName)
	telemetryRuleGeneratedSuccessfullyCounter, _ := telemetryMeter.Int64Counter(
		"generator.rule_generated_successfully.counter",
		metric.WithDescription("Number of operations"),
		metric.WithUnit("{operation}"),
	)
	telemetryRuleGenerationFailedCounter, _ := telemetryMeter.Int64Counter(
		"generator.rule_generation_failed.counter",
		metric.WithDescription("Number of operations"),
		metric.WithUnit("{operation}"),
	)

	rules := []oathkeeper.Rule{}
	for path, p := range g.doc.Paths.Map() {
		for verb, o := range p.Operations() {
			rule, createRuleErr := g.createRule(verb, path, o)
			if createRuleErr != nil {
				telemetryRuleGenerationFailedCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verb", verb)))
				return nil, createRuleErr
			}

			telemetryRuleGeneratedSuccessfullyCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("verb", verb)))
			rules = append(rules, *rule)
		}
	}

	sort.Sort(RulesById(rules))
	return rules, nil
}

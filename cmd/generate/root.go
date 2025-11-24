package generate

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/generator"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/cerberauth/x/telemetryx"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"gopkg.in/yaml.v3"
)

var (
	configFilePath string
	filepath       string
	fileurl        string

	prefix      string
	outputpath  string
	upstreamUrl string
	serverUrls  []string

	jsonOutput bool
	yamlOutput bool
)

var (
	otelName = "github.com/cerberauth/openapi-oathkeeper/cmd/generate"

	errorReasonAttributeKey = attribute.Key("error_reason")
	encodingAttributeKey    = attribute.Key("encoding")
)

func encodeJSON(rules []oathkeeper.Rule) (*bytes.Buffer, error) {
	outputBuf := new(bytes.Buffer)
	enc := json.NewEncoder(outputBuf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")

	if encodeErr := enc.Encode(rules); encodeErr != nil {
		return nil, encodeErr
	}

	return outputBuf, nil
}

func encodeYAML(rules []oathkeeper.Rule) (*bytes.Buffer, error) {
	outputBuf := new(bytes.Buffer)
	enc := yaml.NewEncoder(outputBuf)

	if encodeErr := enc.Encode(rules); encodeErr != nil {
		return nil, encodeErr
	}

	return outputBuf, nil
}

func NewGenerateCmd() (generateCmd *cobra.Command) {
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate Ory Oathkeeper rules from an OpenAPI 3 to file or Std output",
		Run: func(cmd *cobra.Command, args []string) {
			telemetryMeter := telemetryx.GetMeterProvider().Meter(otelName)
			telemetryGenerateSuccessCounter, _ := telemetryMeter.Int64Counter("generate.success.counter")
			telemetryGenerateErrorCounter, _ := telemetryMeter.Int64Counter("generate.error.counter")

			ctx := cmd.Context()
			var cfg *config.Config
			var doc *openapi3.T
			var err error

			if configFilePath != "" {
				cfg, err = config.New(configFilePath)
				if err != nil {
					telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("failed to load config file")))
					log.Fatal(err)
				}
			} else {
				cfg = &config.Config{
					Prefix:     prefix,
					ServerUrls: serverUrls,
					Upstream: oathkeeper.RuleUpstream{
						URL: upstreamUrl,
					},
				}
			}

			if fileurl != "" {
				uri, urlerr := url.Parse(fileurl)
				if urlerr != nil {
					telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("failed to parse url")))
					log.Fatal(urlerr)
				}

				doc, err = openapi3.NewLoader().LoadFromURI(uri)
			}

			if filepath != "" {
				if _, err := os.Stat(filepath); err != nil {
					telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("the openapi file has not been found")))
					log.Fatalf("the openapi file has not been found on %s", filepath)
				}

				doc, err = openapi3.NewLoader().LoadFromFile(filepath)
			}

			if err != nil {
				telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("failed to load openapi file")))
				log.Fatal(err)
			}

			g, err := generator.NewGenerator(ctx, doc, cfg)
			if err != nil {
				telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("failed to create generator")))
				log.Fatal(err)
			}

			rules, err := g.Generate(ctx)
			if err != nil {
				telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(errorReasonAttributeKey.String("failed to generate rules")))
				log.Fatal(err)
			}

			var outputBuf *bytes.Buffer
			var encodeErr error
			var otelEncodingAttributeValue attribute.KeyValue
			if yamlOutput && !jsonOutput {
				otelEncodingAttributeValue = encodingAttributeKey.String("yaml")
				outputBuf, encodeErr = encodeYAML(rules)
			} else {
				otelEncodingAttributeValue = encodingAttributeKey.String("json")
				outputBuf, encodeErr = encodeJSON(rules)
			}

			if encodeErr != nil {
				telemetryGenerateErrorCounter.Add(ctx, 1, metric.WithAttributes(otelEncodingAttributeValue, errorReasonAttributeKey.String("failed to encode rules")))
				log.Fatal(err)
			}

			telemetryGenerateSuccessCounter.Add(ctx, 1, metric.WithAttributes(otelEncodingAttributeValue))

			if outputpath != "" {
				// nolint:errcheck
				os.WriteFile(outputpath, outputBuf.Bytes(), 0600)
				return
			}

			// nolint:errcheck
			os.Stdout.Write(outputBuf.Bytes())
		},
	}

	generateCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "", false, "Use JSON as output format")
	generateCmd.PersistentFlags().BoolVarP(&yamlOutput, "yaml", "", false, "Use YAML as output format")

	generateCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "", "Path to one .yaml, .yml, config file.")
	generateCmd.PersistentFlags().StringVarP(&fileurl, "url", "u", "", "OpenAPI URL")
	generateCmd.PersistentFlags().StringVarP(&filepath, "file", "f", "", "OpenAPI File Path")
	generateCmd.PersistentFlags().StringVarP(&outputpath, "output", "o", "", "Oathkeeper Rules output path")

	generateCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "OpenAPI Prefix Id")
	generateCmd.PersistentFlags().StringArrayVarP(&serverUrls, "server-url", "", nil, "API Server Urls")
	generateCmd.PersistentFlags().StringVarP(&upstreamUrl, "upstream-url", "", "", "The Upstream URL the request will be forwarded to")

	return generateCmd
}

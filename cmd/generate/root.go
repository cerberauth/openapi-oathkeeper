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
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

var tracer = otel.Tracer("cmd/generate")

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
			ctx, span := tracer.Start(cmd.Context(), "Generate")
			defer span.End()

			var cfg *config.Config
			var doc *openapi3.T
			var err error

			if configFilePath != "" {
				cfg, err = config.New(configFilePath)
				if err != nil {
					span.RecordError(err)
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
					span.RecordError(urlerr)
					log.Fatal(urlerr)
				}

				doc, err = openapi3.NewLoader().LoadFromURI(uri)
			}

			if filepath != "" {
				if _, err := os.Stat(filepath); err != nil {
					span.RecordError(err)
					log.Fatalf("the openapi file has not been found on %s", filepath)
				}

				doc, err = openapi3.NewLoader().LoadFromFile(filepath)
			}

			if err != nil {
				span.RecordError(err)
				log.Fatal(err)
			}

			g, err := generator.NewGenerator(ctx, doc, cfg)
			if err != nil {
				span.RecordError(err)
				log.Fatal(err)
			}

			rules, err := g.Generate()
			if err != nil {
				span.RecordError(err)
				log.Fatal(err)
			}

			var outputBuf *bytes.Buffer
			var encodeErr error
			if yamlOutput && !jsonOutput {
				span.SetAttributes(attribute.String("encode", "yaml"))
				outputBuf, encodeErr = encodeYAML(rules)
			} else {
				span.SetAttributes(attribute.String("encode", "json"))
				outputBuf, encodeErr = encodeJSON(rules)
			}

			if encodeErr != nil {
				span.RecordError(encodeErr)
				log.Fatal(err)
			}

			if outputpath != "" {
				os.WriteFile(outputpath, outputBuf.Bytes(), 0644)
				return
			}

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

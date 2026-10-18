package generate

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	fpath "path/filepath"

	"github.com/cerberauth/openapi-oathkeeper/config"
	"github.com/cerberauth/openapi-oathkeeper/generator"
	"github.com/cerberauth/openapi-oathkeeper/oathkeeper"
	"github.com/cerberauth/x/fsx"
	"github.com/cerberauth/x/telemetryx"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
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

			ctx := cmd.Context()
			var cfg *config.Config
			var doc *openapi3.T
			var err error

			if configFilePath != "" {
				if _, err = fsx.FileExists(configFilePath); err != nil {
					telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-read-config-file"))
				}
				cfg, err = config.New(configFilePath)
				if err != nil {
					telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-load-config-file"))
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
					telemetryx.Fatal(cmd, telemetryx.Reason(urlerr, "failed-to-parse-url"))
				}

				doc, err = openapi3.NewLoader().LoadFromURI(uri)
			}

			if filepath != "" {
				var data []byte
				data, err = fsx.ReadFile(filepath)
				if err != nil {
					telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-read-openapi-file"))
				}

				doc, err = openapi3.NewLoader().LoadFromDataWithPath(data, &url.URL{Path: filepath})
			}

			if err != nil {
				telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-load-openapi-file"))
			}

			g, err := generator.NewGenerator(ctx, doc, cfg)
			if err != nil {
				telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-create-generator"))
			}

			rules, err := g.Generate(ctx)
			if err != nil {
				telemetryx.Fatal(cmd, telemetryx.Reason(err, "failed-to-generate-rules"))
			}

			var outputBuf *bytes.Buffer
			var encodeErr error
			format := "json"
			if yamlOutput && !jsonOutput {
				format = "yaml"
				outputBuf, encodeErr = encodeYAML(rules)
			} else {
				outputBuf, encodeErr = encodeJSON(rules)
			}

			if encodeErr != nil {
				telemetryx.Fatal(cmd, telemetryx.Reason(encodeErr, "failed-to-encode-rules"))
			}
			telemetryx.Counter(ctx, "generate.output", 1, format)

			if outputpath != "" {
				if err = fsx.WriteFile(fpath.Clean(outputpath), outputBuf.Bytes(), 0600); err != nil {
					telemetryx.Fatal(cmd, err)
				}
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

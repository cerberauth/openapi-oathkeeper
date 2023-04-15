package generate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"os"

	"github.com/cerberauth/openapi-oathkeeper/generator"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var (
	filepath   string
	fileurl    string
	prefixId   string
	outputpath string

	jwksUris         map[string]string
	allowedIssuers   map[string]string
	allowedAudiences map[string]string
	serverUrls       []string

	upstreamUrl       string
	upstreamStripPath string
)

func NewGenerateCmd() (generateCmd *cobra.Command) {
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate Ory Oathkeeper rules from an OpenAPI 3 to file or Std output",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			var doc *openapi3.T
			var err error

			if fileurl != "" {
				uri, urlerr := url.Parse(fileurl)
				if urlerr != nil {
					panic(urlerr)
				}

				doc, err = openapi3.NewLoader().LoadFromURI(uri)
			}

			if filepath != "" {
				doc, err = openapi3.NewLoader().LoadFromFile(filepath)
			}

			if err != nil {
				panic(err)
			}

			g := generator.NewGenerator(prefixId, jwksUris, allowedIssuers, allowedAudiences, serverUrls, upstreamUrl, upstreamStripPath)
			if loadErr := g.LoadOpenAPI3Doc(ctx, doc); loadErr != nil {
				panic(loadErr)
			}

			rules, err := g.Generate()
			if err != nil {
				panic(err)
			}

			jsonBuf := new(bytes.Buffer)
			enc := json.NewEncoder(jsonBuf)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "    ")

			if encodeErr := enc.Encode(rules); encodeErr != nil {
				panic(encodeErr)
			}

			if outputpath != "" {
				os.WriteFile(outputpath, jsonBuf.Bytes(), 0644)
				return
			}

			os.Stdout.Write(jsonBuf.Bytes())
		},
	}

	generateCmd.PersistentFlags().StringVarP(&prefixId, "prefix", "p", "", "OpenAPI Prefix Id")
	generateCmd.PersistentFlags().StringToStringVarP(&jwksUris, "jwks-uris", "", nil, "JWKS Uris")
	generateCmd.PersistentFlags().StringToStringVarP(&allowedIssuers, "allowed-issuers", "", nil, "Allowed Issuers")
	generateCmd.PersistentFlags().StringToStringVarP(&allowedAudiences, "allowed-audiences", "", nil, "Allowed Audiences")
	generateCmd.PersistentFlags().StringArrayVarP(&serverUrls, "server-url", "", nil, "API Server Urls")
	generateCmd.PersistentFlags().StringVarP(&fileurl, "url", "u", "", "OpenAPI URL")
	generateCmd.PersistentFlags().StringVarP(&filepath, "file", "f", "", "OpenAPI File Path")
	generateCmd.PersistentFlags().StringVarP(&upstreamUrl, "upstream-url", "", "", "The Upstream URL the request will be forwarded to")
	generateCmd.PersistentFlags().StringVarP(&upstreamStripPath, "upstream-strip-path", "", "", "Replaces the provided path prefix when forwarding the requested URL to the upstream URL")
	generateCmd.PersistentFlags().StringVarP(&outputpath, "output", "o", "", "Oathkeeper Rules output path")

	return generateCmd
}

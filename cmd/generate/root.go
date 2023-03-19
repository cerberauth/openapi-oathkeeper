package generate

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/cerberauth/openapi-oathkeeper/generator"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

var (
	_, b, _, _      = runtime.Caller(0)
	basepath        = filepath.Dir(b)
	swaggerFilepath string
	swaggerUrl      string
	outputPath      string
)

func NewGenerateCmd() (generateCmd *cobra.Command) {
	generateCmd = &cobra.Command{
		Use: "generate",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			var doc *openapi3.T
			var err error

			if swaggerUrl != "" {
				uri, urlerr := url.Parse(swaggerUrl)
				if urlerr != nil {
					fmt.Print(urlerr)
					return
				}

				doc, err = openapi3.NewLoader().LoadFromURI(uri)
			}

			if swaggerFilepath != "" {
				doc, err = openapi3.NewLoader().LoadFromFile(swaggerFilepath)
			}

			if err != nil {
				fmt.Print(err)
				return
			}

			rules, err := generator.New().Document(doc).Generate(ctx)
			if err != nil {
				fmt.Print(err)
				return
			}

			fmt.Print(string(rules))
		},
	}
	generateCmd.PersistentFlags().StringVarP(&swaggerUrl, "url", "u", "", "OpenAPI URL")
	generateCmd.PersistentFlags().StringVarP(&swaggerFilepath, "file", "f", "", "OpenAPI File Path")
	generateCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", ".", "OAthKeeper Rules output path")

	return generateCmd
}

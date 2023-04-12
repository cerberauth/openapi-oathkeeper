package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/cerberauth/openapi-oathkeeper/cmd/generate"
)

func NewRootCmd() (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{
		Use:   "openapi-oathkeeper",
		Short: "Generate Ory Oathkeeper Rules from OpenAPI 3.0 files",
	}
	rootCmd.AddCommand(generate.NewGenerateCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	c := NewRootCmd()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

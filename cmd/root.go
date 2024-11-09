package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cerberauth/openapi-oathkeeper/cmd/generate"
	"github.com/cerberauth/openapi-oathkeeper/internal/analytics"
)

var sqaOptOut bool

func NewRootCmd(projectVersion string) (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{
		Use:   "openapi-oathkeeper",
		Short: "Generate Ory Oathkeeper Rules from OpenAPI 3.0 files",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !sqaOptOut {
				_, err := analytics.NewAnalytics(cmd.Context(), projectVersion)
				if err != nil {
					fmt.Println("Failed to initialize analytics:", err)
				}
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !sqaOptOut {
				analytics.Close()
			}
		},
	}
	rootCmd.AddCommand(generate.NewGenerateCmd())

	rootCmd.PersistentFlags().BoolVarP(&sqaOptOut, "sqa-opt-out", "", false, "Opt out of sending anonymous usage statistics and crash reports to help improve the tool")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(projectVersion string) {
	c := NewRootCmd(projectVersion)

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

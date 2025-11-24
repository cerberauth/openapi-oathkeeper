package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/cerberauth/openapi-oathkeeper/cmd/generate"
	"github.com/cerberauth/x/telemetryx"
)

var (
	sqaOptOut    bool
	otelShutdown func(context.Context) error
)

var name = "openapi-oathkeeper"

func NewRootCmd(projectVersion string) (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{
		Use:   name,
		Short: "Generate Ory Oathkeeper Rules from OpenAPI 3.0 files",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !sqaOptOut {
				otelShutdown, _ = telemetryx.New(cmd.Context(), name, projectVersion)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if otelShutdown != nil {
				_ = otelShutdown(cmd.Context())
				otelShutdown = nil
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
	defer func() {
		if otelShutdown != nil {
			_ = otelShutdown(context.Background())
			otelShutdown = nil
		}
	}()

	if err := c.Execute(); err != nil {
		if otelShutdown != nil {
			_ = otelShutdown(context.Background())
			otelShutdown = nil
		}

		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		// nolint: gocritic // false positive
		os.Exit(1)
	}
}

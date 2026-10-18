package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/cerberauth/openapi-oathkeeper/cmd/generate"
	"github.com/cerberauth/x/telemetryx"
)

// argvioAPIKey is openapi-oathkeeper's public_ingest-scoped Argvio API key, selecting the
// "openapi-oathkeeper" tenant its usage analytics land in. It only allows sending usage
// data, so it is safe to embed.
const argvioAPIKey = "argv_live_pub_MJEWEivQJp4PnYwNi0-V5NhE3mU6CpiO" //nolint:gosec // ingest-only key, not a secret

var (
	sqaOptOut bool
)

var name = "openapi-oathkeeper"

func NewRootCmd(projectVersion, commit, date string) (cmd *cobra.Command) {
	var rootCmd = &cobra.Command{
		Use:     name,
		Version: projectVersion + " (commit=" + commit + ", built=" + date + ")",
		Short:   "Generate Ory Oathkeeper Rules from OpenAPI 3.0 files",
	}
	rootCmd.AddCommand(generate.NewGenerateCmd())

	rootCmd.PersistentFlags().BoolVarP(&sqaOptOut, "sqa-opt-out", "", false, "Opt out of sending anonymous usage statistics and crash reports to help improve the tool")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(projectVersion, commit, date string) {
	c := NewRootCmd(projectVersion, commit, date)
	argvioShutdown := telemetryx.Instrument(c, projectVersion, argvioAPIKey)
	defer func() { _ = argvioShutdown(context.Background()) }()

	if err := c.Execute(); err != nil {
		_ = argvioShutdown(context.Background())

		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		// nolint: gocritic // false positive
		os.Exit(1)
	}
}

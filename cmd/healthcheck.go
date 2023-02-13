package cmd

import (
	"sophos_tenant_cli/sophos"

	"github.com/spf13/cobra"
)

var healthCheck = &cobra.Command{
	Use:     "healthsummary",
	Short:   "Provides a summary of the health of the tenant, only tenants with something to report are shown.",
	Aliases: []string{"summary"},
	Run: func(cmd *cobra.Command, args []string) {
		sophos.Execute()
	},
}

func init() {
	rootCmd.AddCommand(healthCheck)
}

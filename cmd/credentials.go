package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCredentals = &cobra.Command{
	Use:                   "setcredentials [client_id] [client_secret]",
	Short:                 "Sets the Partner/Enterprise tenant credentials",
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.RangeArgs(2, 2),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("client_id", args[0])
		viper.Set("client_secret", args[1])
		viper.WriteConfigAs("config.json")
	},
}

func init() {
	rootCmd.AddCommand(setCredentals)
}

package podfetcher

import (
	"log"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gregf/podfetcher/lib/commands"
)

func Execute() {
	initConfig()

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "Updates the database with the latest episodes to be fetched.",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Update()
		},
	}

	var cmdFetch = &cobra.Command{
		Use:   "fetch",
		Short: "Fetches podcast episodes that have not been downloaded.",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Fetch()
		},
	}

	var cmdCatchUp = &cobra.Command{
		Use:   "catchup",
		Short: "Marks all podcast episodes as downloaded",
		Run: func(cmd *cobra.Command, args []string) {
			commands.CatchUp()
		},
	}
	var rootCmd = &cobra.Command{Use: "podfetcher"}
	rootCmd.AddCommand(cmdUpdate, cmdFetch, cmdCatchUp)
	rootCmd.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.podfetcher")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error config file %s \n", err)
		return
	}
}

package podfetcher

import (
	"log"

	"github.com/gregf/podfetcher/lib/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	var rootCmd = &cobra.Command{Use: "podfetcher"}
	rootCmd.AddCommand(cmdUpdate)
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

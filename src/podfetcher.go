package podfetcher

import (
	"log"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gregf/podfetcher/src/commands"
)

// Execute parses command line args and fires up commands
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

	var cmdLsNew = &cobra.Command{
		Use:   "lsnew",
		Short: "Display new episodes to be downloaded.",
		Run: func(cmd *cobra.Command, args []string) {
			commands.LsNew()
		},
	}

	var cmdImport = &cobra.Command{
		Use:   "import",
		Short: "Import feeds from a opml file.",
		Run: func(cmd *cobra.Command, args []string) {
			commands.Import(args)
		},
	}

	var rootCmd = &cobra.Command{Use: "podfetcher"}
	rootCmd.AddCommand(cmdUpdate, cmdFetch, cmdCatchUp, cmdLsNew, cmdImport)
	rootCmd.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/podfetcher")
	viper.AddConfigPath("$HOME/.config/podfetcher")
	viper.AddConfigPath("$XDG_CONFIG_HOME/podfetcher")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error config file %s \n", err)
		return
	}
}
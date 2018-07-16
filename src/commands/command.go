package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gregf/podfetcher/src/database"

	"github.com/nightlyone/lockfile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Podfetcher Version number
const podFetcherVersion = "v0.5"

// Env struct
type Env struct {
	db database.Datastore
}

// Execute parses command line args and fires up commands
func Execute() {
	lf := "/tmp/podfetcher.lock"

	initConfig()
	createLock(lf)
	defer unLock(lf)
	trapInit(lf)

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db}

	var podcastID int

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Podfetcher: %s\n", podFetcherVersion)
		},
	}

	var cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "Updates the database with the latest episodes to be fetched.",
		Run:   env.Update,
	}

	var cmdFetch = &cobra.Command{
		Use:   "fetch",
		Short: "Fetches podcast episodes that have not been downloaded.",
		Run:   env.Fetch,
	}

	var cmdCatchUp = &cobra.Command{
		Use:   "catchup",
		Short: "Marks all podcast episodes as downloaded",
		Run: func(cmd *cobra.Command, args []string) {
			env.CatchUp(podcastID)
		},
	}

	var cmdLsNew = &cobra.Command{
		Use:   "lsnew",
		Short: "Display new episodes to be downloaded.",
		Run:   env.LsNew,
	}

	var cmdImport = &cobra.Command{
		Use:   "import",
		Short: "Import feeds from a opml file.",
		Run: func(cmd *cobra.Command, args []string) {
			Import(args)
		},
	}

	var cmdAdd = &cobra.Command{
		Use:   "add",
		Short: "Add a feed to your feeds file.",
		Run: func(cmd *cobra.Command, args []string) {
			Add(args)
		},
	}

	var cmdLsCasts = &cobra.Command{
		Use:   "lscasts",
		Short: "Displays a list of subscribed podcasts",
		Run:   env.LsCasts,
	}

	var cmdPause = &cobra.Command{
		Use:   "pause",
		Short: "Toggles between paused states.",
		Long:  "Toggles between paused states. Paused Podcasts are ignored.",
		PreRun: func(cmd *cobra.Command, args []string) {
			if podcastID == 0 {
				fmt.Println("It looks like you forget to set --cast ID")
				fmt.Println("You can get --cast ID from lspodcasts")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			env.Pause(podcastID)
		},
	}

	cmdPause.Flags().IntVarP(&podcastID, "cast", "c", 0, "Podcast ID")
	cmdCatchUp.Flags().IntVarP(&podcastID, "cast", "c", 0, "Podcast ID")
	var rootCmd = &cobra.Command{Use: "podfetcher"}
	rootCmd.AddCommand(
		cmdAdd,
		cmdCatchUp,
		cmdFetch,
		cmdImport,
		cmdLsCasts,
		cmdLsNew,
		cmdPause,
		cmdUpdate,
		cmdVersion)
	rootCmd.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/podfetcher")
	viper.AddConfigPath("$HOME/.config/podfetcher")
	viper.AddConfigPath("$XDG_CONFIG_HOME/podfetcher")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error config file %s\n", err)
		Setup()
		return
	}
}

func createLock(lockFile string) {
	lock, err := lockfile.New(lockFile)
	if err != nil {
		log.Fatalf("Cannot init lock. reason: %v", err)
	}

	err = lock.TryLock()
	if err != nil {
		log.Fatalf("Podfetcher: %s\n", err)
	}
}

func unLock(lockFile string) {
	lock, err := lockfile.New(lockFile)
	if err != nil {
		log.Fatalf("Podfetcher: %s\n", err)
	}

	lock.Unlock()
}

func trapInit(lockFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		unLock(lockFile)
		os.Exit(1)
	}()

}

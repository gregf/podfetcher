package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gregf/podfetcher/lib/database"
)

var slash = string(os.PathSeparator)

// Fetch loops through episodes where downloaded = false and downloads them.
func Fetch() {
	urls, err := database.FindNewEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		fmt.Printf("Fetching: %s - %s\n", database.FindPodcastTitleByURL(url),
			database.FindEpisodeTitleByURL(url))
		download(url)
	}
}

func expandPath(path string) string {
	length := len(path)

	if length == 0 {
		return path
	}

	// replace env variables
	expandedPath := os.ExpandEnv(path)

	// replace ~ with $HOME
	if (length == 1 && path[0] == '~') || (length > 1 && path[:2] == "~/") {
		usr, _ := user.Current()

		expandedPath = strings.Replace(expandedPath, "~", usr.HomeDir, 1)
	} else if path[:1] == "~" {
		// replace ~user with their $HOME

		firstSlash := strings.Index(path, slash)
		if firstSlash < 0 {
			firstSlash = len(path)
		}

		if firstSlash > 1 {
			usr, err := user.Lookup(path[1:firstSlash])
			if err == nil {
				expandedPath = usr.HomeDir + path[firstSlash:]
			}
		}
	}

	// get an absolute path, ignoring errors
	if absPath, err := filepath.Abs(expandedPath); err == nil {
		expandedPath = absPath
	}

	// cleanup
	expandedPath = filepath.Clean(expandedPath)

	return expandedPath
}

func download(url string) {
	if strings.Contains(url, "youtube.com") {
		ytdl(url)
	} else {
		wget(url)
	}
	database.SetDownloadedByURL(url)
}

func run(cmdName string, cmdArgs []string) {
	cmd := exec.Command(cmdName, cmdArgs...)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

func wget(url string) {
	title := makeTitle(database.FindPodcastTitleByURL(url))
	saveLoc := filepath.Join(expandPath(viper.GetString("download")), title, getFileName(url, false))
	err := os.MkdirAll(filepath.Join(expandPath(viper.GetString("download")), title), 0755)
	if err != nil {
		log.Fatal(err)
	}
	cmdName := "wget"
	cmdArgs := []string{"-c", url, "-O", saveLoc}
	run(cmdName, cmdArgs)
}

func ytdl(url string) {
	title := makeTitle(database.FindPodcastTitleByURL(url))
	saveLoc := filepath.Join(viper.GetString("download"), title, getFileName(url, true))
	err := os.MkdirAll(filepath.Join(viper.GetString("download"), title), 0755)
	if err != nil {
		log.Fatal(err)
	}
	cmdName := "youtube-dl"
	cmdArgs := []string{"--no-playlist", "--continue", "--no-part", "-f", viper.GetString("youtube-quality"), "-o", saveLoc, url}
	run(cmdName, cmdArgs)
}

func getFileName(enclosureURL string, youtube bool) (filename string) {
	if youtube {
		ytdlCmd := exec.Command("youtube-dl", "--get-filename", enclosureURL)
		ytdlOut, err := ytdlCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		filename = string(ytdlOut)
		return filename
	}

	url, err := url.Parse(enclosureURL)
	if err != nil {
		log.Fatal(err)
	}
	filename = filepath.Base(url.Path)

	return filename
}

func makeTitle(title string) (newtitle string) {
	newtitle = strings.Replace(title, " ", "", -1)
	newtitle = strings.Replace(newtitle, "%20", "", -1)

	return newtitle
}

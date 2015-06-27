package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/filter"
	"github.com/gregf/podfetcher/src/helpers"
)

// Fetch loops through episodes where downloaded = false and downloads them.
func Fetch() {
	urls, err := database.FindNewEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		podcastTitle := database.FindPodcastTitleByURL(url)
		episodeTitle := database.FindEpisodeTitleByURL(url)
		download(podcastTitle, episodeTitle, url)
	}
}

func download(podcastTitle string, episodeTitle string, url string) {
	if filter.Run(podcastTitle, episodeTitle) {
		fmt.Printf("Filtered: %s - %s\n", podcastTitle, episodeTitle)
		database.SetDownloadedByURL(url)
		return
	}
	fmt.Printf("Fetching: %s - %s\n", podcastTitle, episodeTitle)
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
	download := helpers.ExpandPath(viper.GetString("main.download"))
	saveLoc := filepath.Join(
		download,
		title,
		getFileName(url, false))
	err := os.MkdirAll(filepath.Join(
		download,
		title),
		0755)
	if err != nil {
		log.Fatal(err)
	}
	cmdName := "wget"
	cmdArgs := []string{"-c", url, "-O", saveLoc}
	run(cmdName, cmdArgs)
}

func ytdl(url string) {
	title := makeTitle(database.FindPodcastTitleByURL(url))
	download := helpers.ExpandPath(viper.GetString("main.download"))
	saveLoc := filepath.Join(
		download,
		title,
		getFileName(url, true))
	err := os.MkdirAll(filepath.Join(
		download,
		title), 0755)
	if err != nil {
		log.Fatal(err)
	}
	cmdName := "youtube-dl"
	cmdArgs := []string{
		"--no-playlist",
		"--continue",
		"--no-part",
		"-f",
		viper.GetString("main.youtube-quality"),
		"-o",
		saveLoc,
		url}
	run(cmdName, cmdArgs)
}

func getFileName(enclosureURL string, youtube bool) (filename string) {
	if youtube {
		ytdlCmd := exec.Command("youtube-dl", "--get-filename", enclosureURL)
		ytdlOut, err := ytdlCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		filename = strings.Split(string(ytdlOut), "\n")[0]
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

package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/cavaliercoder/grab"
	"github.com/juju/deputy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gregf/podfetcher/src/filter"
	"github.com/gregf/podfetcher/src/helpers"
)

// Params struct for downloader
type Params struct {
	url     string
	yturl   string
	youtube bool
}

// Fetch loops through episodes where downloaded = false and downloads them.
func (env *Env) Fetch(cmd *cobra.Command, args []string) {
	urls, err := env.db.FindNewEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		podcastTitle := env.db.FindPodcastTitleByURL(url)
		episodeTitle := env.db.FindEpisodeTitleByURL(url)
		env.download(podcastTitle, episodeTitle, url)
	}
}

func (env *Env) download(podcastTitle, episodeTitle, url string) {
	if filter.Run(podcastTitle, episodeTitle) {
		fmt.Printf("Filtered: %s - %s\n", podcastTitle, episodeTitle)
		env.db.SetDownloadedByURL(url)
		return
	}
	s := false
	fmt.Printf("Downloading: %s - %s\n", podcastTitle, episodeTitle)
	if strings.Contains(strings.ToLower(url), "youtube.com") {
		s = env.downloader(Params{url: getYoutubeURL(url), yturl: url, youtube: true})
	} else {
		s = env.downloader(Params{url: url, youtube: false})
	}
	if s {
		env.db.SetDownloadedByURL(url)
		notify(fmt.Sprintf("Fetched: %s - %s", podcastTitle, episodeTitle))
	} else {
		notify(fmt.Sprintf("Failed: %s - %s", podcastTitle, episodeTitle))
	}
}

func run(cmdName string, cmdArgs []string) (cmdOut string) {
	d := deputy.Deputy{
		Errors:    deputy.FromStderr,
		StdoutLog: func(b []byte) { cmdOut = string(b) },
	}
	cmd := exec.Command(cmdName, cmdArgs...)
	err := d.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}
	return cmdOut
}

func notify(msg string) {
	cmdName := viper.GetString("main.notify-program")
	cmdArgs := []string{
		"podfetcher:",
		fmt.Sprintf("%s", msg),
	}
	if len(cmdName) <= 0 {
		return
	}
	run(cmdName, cmdArgs)
}

func (env *Env) downloader(p Params) bool {
	var fileName string
	var title string
	if p.youtube {
		fileName = getFileName(p.yturl, p.youtube)
		title = makeTitle(env.db.FindPodcastTitleByURL(p.yturl))
	} else {
		fileName = getFileName(p.url, p.youtube)
		title = makeTitle(env.db.FindPodcastTitleByURL(p.url))
	}

	dlDir := helpers.ExpandPath(viper.GetString("main.download"))
	saveLoc := filepath.Join(dlDir, title, fileName)
	err := os.MkdirAll(filepath.Join(dlDir, title), 0755)
	if err != nil {
		log.Fatalf("mkdir failed %s\n", err)
	}
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(saveLoc, p.url)

	// start download
	resp := client.Do(req)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\rTransfered %v / %v (%d%%) - %.0fKBp/s  ",
				datasize.ByteSize(resp.BytesComplete()).HumanReadable(),
				datasize.ByteSize(resp.Size).HumanReadable(),
				int(100*resp.Progress()),
				resp.BytesPerSecond()/1024)
		case <-resp.Done:
			fmt.Printf("\rTransfered %v / %v (100%%) - 0KBp/s  ",
				datasize.ByteSize(resp.BytesComplete()).HumanReadable(),
				datasize.ByteSize(resp.Size).HumanReadable())
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "\nDownload failed: %v\n", err)
		return false
	}
	fmt.Println("")
	return true
}

func getYoutubeURL(url string) (yturl string) {
	cmdName := "youtube-dl"
	cmdArgs := []string{
		"--format",
		viper.GetString("main.youtube-quality"),
		"--get-url",
		url}
	cmdOut := run(cmdName, cmdArgs)
	yturl = strings.Split(string(cmdOut), "\n")[0]
	return yturl
}

func getFileName(enclosureURL string, youtube bool) (filename string) {
	if youtube {
		cmdName := "youtube-dl"
		cmdArgs := []string{
			"--get-filename",
			enclosureURL}
		cmdOut := run(cmdName, cmdArgs)
		filename = strings.Split(string(cmdOut), "\n")[0]
		return filename
	}

	url, err := url.Parse(enclosureURL)
	if err != nil {
		log.Fatal(err)
	}
	filename = filepath.Base(url.Path)

	return filename
}

func makeTitle(title string) (t string) {
	title = strings.Replace(title, " ", "", -1)
	title = strings.Replace(title, "%20", "", -1)
	title = strings.Replace(title, "/", "-", -1)
	title = strings.Replace(title, "\\", "-", -1)
	title = strings.Replace(title, "'", "", -1)
	title = strings.Replace(title, "\"", "", -1)
	title = strings.Replace(title, ",", "", -1)

	return title
}

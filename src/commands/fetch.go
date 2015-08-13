package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
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
	fmt.Printf("Downloading: %s - %s\n", podcastTitle, episodeTitle)
	if strings.Contains(strings.ToLower(url), "youtube.com") {
		env.downloader(Params{url: getYoutubeURL(url), yturl: url, youtube: true})
	} else {
		env.downloader(Params{url: url, youtube: false})
	}
	env.db.SetDownloadedByURL(url)
	notify(podcastTitle, episodeTitle)
}

func run(cmdName string, cmdArgs []string) (cmdOut string) {
	d := deputy.Deputy{
		Errors:    deputy.FromStderr,
		StdoutLog: func(b []byte) { cmdOut = string(b) },
		Timeout:   time.Second * 30,
	}
	cmd := exec.Command(cmdName, cmdArgs...)
	err := d.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}
	return cmdOut
}

func notify(podcastTitle, episodeTitle string) {
	cmdName := viper.GetString("main.notify-program")
	cmdArgs := []string{
		"podfetcher:",
		fmt.Sprintf("Fetched: %s - %s", podcastTitle, episodeTitle),
	}
	if len(cmdName) <= 0 {
		return
	}
	run(cmdName, cmdArgs)
}

func (env *Env) downloader(p Params) {
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
	/*
		Create new file.
		Filename from fileName variable
	*/
	file, err := os.Create(saveLoc)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/*
		check status and CheckRedirect
	*/
	checkStatus := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	/*
		Get Response: 200 OK?
	*/
	response, err := checkStatus.Get(p.url)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	/*
		fileSize example: 12572 bytes
	*/
	filesize := response.ContentLength
	go func() {
		n, err := io.Copy(file, response.Body)
		if n != filesize {
			fmt.Println("Truncated")
		}
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	countSize := int(filesize)
	bar := pb.StartNew(countSize)
	var fi os.FileInfo
	for fi == nil || fi.Size() < filesize {
		fi, _ = file.Stat()
		bar.Set(int(fi.Size()))
		bar.ShowBar = false
		bar.ShowSpeed = true
		bar.SetUnits(pb.U_BYTES)
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("")
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

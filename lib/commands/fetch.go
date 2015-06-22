package commands

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gregf/podfetcher/lib/database"
	"github.com/spf13/viper"
)

import ()

func Fetch() {
	urls, err := database.FindNewEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		fmt.Printf("Fetching: %s ...\n", url)
		download(url)
	}
}

func download(url string) {
	if strings.Contains(url, "youtube.com") {
		ytdl(url)
	} else {
		wget(url)
	}
	database.SetDownloadedByUrl(url)
}

func wget(url string) {
	title := makeTitle(database.FindTitleByUrl(url))
	saveLoc := filepath.Join(viper.GetString("download"), title, getFileName(url, false))
	err := os.MkdirAll(filepath.Join(viper.GetString("download"), title), 0755)
	if err != nil {
		log.Fatal(err)
	}
	wgetCmd := exec.Command("wget", "-c", url, "-O", saveLoc)
	wgetOut, err := wgetCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(wgetOut))
}

func ytdl(url string) {
	title := makeTitle(database.FindTitleByUrl(url))
	saveLoc := filepath.Join(viper.GetString("download"), title, getFileName(url, true))
	err := os.MkdirAll(filepath.Join(viper.GetString("download"), title), 0755)
	if err != nil {
		log.Fatal(err)
	}
	ytdlCmd := exec.Command("youtube-dl", "--no-playlist", "--continue", "--no-part", "-f", viper.GetString("youtube-quality"), "-o", saveLoc, url)

	ytdlOut, err := ytdlCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ytdlOut))
}

func getFileName(enclosureUrl string, youtube bool) (filename string) {
	if youtube {
		ytdlCmd := exec.Command("youtube-dl", "--get-filename", enclosureUrl)
		ytdlOut, err := ytdlCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		filename = string(ytdlOut)
		return filename
	} else {
		url, err := url.Parse(enclosureUrl)
		if err != nil {
			log.Fatal(err)
		}
		filename = filepath.Base(url.Path)

		return filename
	}
}

func makeTitle(title string) (newtitle string) {
	newtitle = strings.Replace(title, " ", "", -1)
	newtitle = strings.Replace(newtitle, "%20", "", -1)

	return newtitle
}

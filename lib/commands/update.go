package commands

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	rss "github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/jteeuwen/go-pkg-rss"
	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
	"github.com/gregf/podfetcher/lib/database"
)

var EnclosureError = "item %s has no enclosure url"

func feedsPath() (path string) {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "feeds")
}

func Update() {
	feeds, err := readFeeds(feedsPath())
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for _, feedURL := range feeds {
		feed := rss.New(0, true, chanHandler, itemHandler)
		err := feed.Fetch(feedURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", feedURL, err)
			return
		}
	}
}

func readLines(path string) (lines []string, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	lines = strings.Fields(string(content))
	return lines, err
}

func readFeeds(path string) (lines []string, err error) {
	feeds, err := readLines(path)
	if err != nil {
		log.Fatal(err)
	}

	comment, err := regexp.Compile(`\A#`)
	if err != nil {
		log.Fatalf("comment %s\n", err)
	}
	for _, line := range feeds {
		if comment.Match([]byte(line)) {
			return
		}
		if len(line) <= 0 {
			return
		}
		lines = append(lines, line)
	}
	return lines, err
}

func itemHandler(f *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	database.AddPodcast(ch.Title, f.Url)

	var maxEpisodes = viper.GetInt("episodes")
	var items []*rss.Item
	if len(newitems) < maxEpisodes {
		items = newitems[0:len(newitems)]
	} else {
		items = newitems[0:maxEpisodes]
	}
	for _, item := range items {
		var enclosureUrl string
		if strings.Contains(f.Url, "youtube.com") {
			if len(item.Links) > 0 {
				enclosureUrl = item.Links[0].Href
			} else {
				log.Printf(EnclosureError, item.Title)
				return
			}
		} else {
			if len(item.Enclosures) > 0 {
				enclosureUrl = item.Enclosures[0].Url
			} else {
				log.Printf(EnclosureError, item.Title)
				return
			}
		}
		items := make(map[string]string)
		items["title"] = item.Title
		items["rssUrl"] = f.Url
		items["enclosureUrl"] = enclosureUrl
		items["pubdate"] = item.PubDate
		if item.Guid != nil {
			items["guid"] = *item.Guid
		} else {
			items["guid"] = item.PubDate + f.Url
		}
		database.AddItem(items)
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
}

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "iso-8859-1" {
		return r, nil
	}
	return nil, errors.New("Unsupported character set encoding: " + charset)
}

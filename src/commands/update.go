package commands

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/gregf/podfetcher/src/helpers"
	rss "github.com/mattn/go-pkg-rss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var enclosureError = "%s is missing a enclosure url"

// Update loops over the feeds file and inserts podcasts + episodes into the
// database.
func (env *Env) Update(cmd *cobra.Command, args []string) {
	feeds, err := helpers.ReadFeeds(helpers.FeedsPath())
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for _, feedURL := range feeds {
		feed := rss.New(0, true, chanHandler, env.itemHandler)
		err := feed.Fetch(feedURL, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func (env *Env) itemHandler(f *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	env.db.AddPodcast(ch.Title, f.Url)

	var maxEpisodes = viper.GetInt("main.episodes")
	var items []*rss.Item
	if len(newitems) < maxEpisodes {
		items = newitems[0:]
	} else {
		items = newitems[0:maxEpisodes]
	}
	for _, item := range items {
		var enclosureURL string
		if strings.Contains(f.Url, "youtube.com") {
			if len(item.Links) > 0 {
				enclosureURL = item.Links[0].Href
			} else {
				log.Printf(enclosureError, item.Title)
				return
			}
		} else {
			if len(item.Enclosures) > 0 {
				enclosureURL = item.Enclosures[0].Url
			} else {
				log.Printf(enclosureError, item.Title)
				return
			}
		}
		items := make(map[string]string)
		items["title"] = item.Title
		items["rssURL"] = f.Url
		items["enclosureURL"] = enclosureURL
		items["pubdate"] = item.PubDate
		if item.Guid != nil {
			items["guid"] = *item.Guid
		} else {
			items["guid"] = item.PubDate + f.Url
		}
		env.db.AddItem(items)
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

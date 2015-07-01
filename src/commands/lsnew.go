package commands

import (
	"fmt"

	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/filter"
	"github.com/gregf/podfetcher/src/helpers"
)

// LsNew prints out episodes where downloaded = false
func LsNew() {
	new := database.FindEpisodesWithPodcastTitle()
	podcastCount := 0
	episodeCount := 0

	if len(new) != 0 {
		fmt.Printf("Episodes marked with [*] have been filtered\n\n")
	}
	for podcastTitle, episodeTitle := range new {
		podcastCount++
		for _, t := range episodeTitle {
			episodeCount++
			var filtered string
			if filter.Run(podcastTitle, t) {
				filtered = "[*]"
			}
			w1 := int(helpers.GetWidth() / 4)
			w2 := int(helpers.GetWidth() / 2)
			fmt.Printf("%-3s %-*.*s - %-.*s\n",
				filtered,
				w1,
				w1,
				podcastTitle,
				w2,
				t)
		}
	}
	if len(new) != 0 {
		fmt.Printf("\n%d episode(s) to consider from %d podcast(s)\n",
			episodeCount,
			podcastCount)
	}
}

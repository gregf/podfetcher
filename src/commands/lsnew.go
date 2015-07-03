package commands

import (
	"fmt"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/apcera/termtables"

	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/filter"
)

// LsNew prints out episodes where downloaded = false
func LsNew() {
	new := database.FindEpisodesWithPodcastTitle()
	podcastCount := 0
	episodeCount := 0

	table := termtables.CreateTable()
	table.AddHeaders("Filtered", "Podcast", "Episode Title")
	var ts = &termtables.TableStyle{
		SkipBorder:  true,
		PaddingLeft: 0, PaddingRight: 2,
		Width:     80,
		Alignment: termtables.AlignLeft,
	}
	table.Style = ts
	for podcastTitle, episodeTitle := range new {
		podcastCount++
		for _, t := range episodeTitle {
			episodeCount++
			var filtered string
			if filter.Run(podcastTitle, t) {
				filtered = "[*]"
			}
			table.AddRow(filtered, podcastTitle, t)
		}
	}

	if len(new) != 0 {
		fmt.Println(table.Render())
		fmt.Printf("\n%d episode(s) to consider from %d podcast(s)\n",
			episodeCount,
			podcastCount)
	}
}

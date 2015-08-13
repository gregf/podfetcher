package commands

import (
	"fmt"
	"log"

	"github.com/apcera/termtables"
	"github.com/spf13/cobra"

	"github.com/gregf/podfetcher/src/filter"
)

// LsNew prints out episodes where downloaded = false
func (env *Env) LsNew(cmd *cobra.Command, args []string) {
	eps, err := env.db.FindEpisodesWithPodcastTitle()
	if err != nil {
		log.Fatal(err)
	}
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
	for podcastTitle, episodeTitle := range eps {
		podcastCount++
		for _, t := range episodeTitle {
			episodeCount++
			var filtered string
			if filter.Run(podcastTitle, t) {
				filtered = "✓"
			}
			table.AddRow(filtered, podcastTitle, t)
		}
	}

	if len(eps) != 0 {
		fmt.Println(table.Render())
		fmt.Printf("\n%d episode(s) to consider from %d podcast(s)\n",
			episodeCount,
			podcastCount)
	}
}

package commands

import (
	"fmt"

	"github.com/apcera/termtables"

	"github.com/gregf/podfetcher/src/database"
)

// LsCasts displays a list of podcasts you are subscribed to.
func LsCasts() {
	ids, titles := database.FindAllPodcasts()
	table := termtables.CreateTable()
	var ts = &termtables.TableStyle{
		SkipBorder:  true,
		PaddingLeft: 0, PaddingRight: 2,
		Width:     80,
		Alignment: termtables.AlignLeft,
	}
	table.Style = ts
	table.AddHeaders("ID", "Podcast Title")
	for i, id := range ids {
		table.AddRow(id, titles[i])
	}
	fmt.Println(table.Render())
}

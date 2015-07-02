package commands

import (
	"fmt"

	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/helpers"
)

// LsCasts displays a list of podcasts you are subscribed to.
func LsCasts() {
	ids, titles := database.FindAllPodcasts()

	for i, id := range ids {
		w1 := int(helpers.GetWidth() - 5)
		fmt.Printf("%d - %.*s\n",
			id,
			w1,
			titles[i])
	}
}

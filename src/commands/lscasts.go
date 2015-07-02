package commands

import (
	"fmt"

	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/helpers"

	// Required for Iter function
	_ "github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/cevaris/ordered_map"
)

// LsCasts displays a list of podcasts you are subscribed to.
func LsCasts() {
	casts := database.FindAllPodcasts()

	for kv := range casts.Iter() {
		w1 := int(helpers.GetWidth() - 5)
		id := kv.Key
		title := kv.Value
		fmt.Printf("%d - %.*s\n",
			id,
			w1,
			title)
	}
}

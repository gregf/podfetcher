package commands

import (
	"fmt"

	"github.com/gregf/podfetcher/src/database"
)

// Pause pauses toggles the pause state of a podcast, which pauses new downloads.
func Pause(id int) {
	state := database.TogglePause(id)
	title := database.FindPodcastTitle(id)
	pausestate := pausestate(state)

	fmt.Printf("%s is now %s\n", title, pausestate)
}

func pausestate(state bool) (pausestate string) {
	if state {
		pausestate = "paused"
	} else {
		pausestate = "unpaused"
	}

	return pausestate
}

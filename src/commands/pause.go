package commands

import "fmt"

// Pause pauses toggles the pause state of a podcast, which pauses new downloads.
func (env *Env) Pause(id int) {
	state := env.db.TogglePause(id)
	title := env.db.FindPodcastTitle(id)
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

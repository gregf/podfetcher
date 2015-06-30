package commands

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gregf/podfetcher/src/database"
	"github.com/gregf/podfetcher/src/filter"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

const (
	_TIOCGWINSZ = 0x5413 // OSX 1074295912
)

// LsNew prints out episodes where downloaded = false
func LsNew() {
	new := database.FindEpisodesWithPodcastTitle()

	fmt.Printf("Episodes marked with [*] have been filtered\n\n")
	for podcastTitle, episodeTitle := range new {
		for _, t := range episodeTitle {
			var filtered string
			if filter.Run(podcastTitle, t) {
				filtered = "[*]"
			}
			w1 := int(getWidth() / 4)
			w2 := int(getWidth() / 2)
			fmt.Printf("%-3s %-*.*s - %-.*s\n",
				filtered,
				w1,
				w1,
				podcastTitle,
				w2,
				t)
		}
	}
}

func getWidth() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
}

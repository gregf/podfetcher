package commands

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gregf/podfetcher/src/database"
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

	for title, eptitle := range new {
		for _, t := range eptitle {
			w := int(getWidth() / 3)
			fmt.Printf("%-*.*s - %-.*s\n",
				w,
				w,
				title,
				w,
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

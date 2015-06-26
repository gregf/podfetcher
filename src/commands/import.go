package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/gilliek/go-opml/opml"
	"github.com/gregf/podfetcher/src/helpers"
)

// Import imports feeds from a opml file.
func Import(args []string) {
	if len(args) <= 0 {
		fmt.Println("podfetcher import <path>")
		return
	}

	doc, err := opml.NewOPMLFromFile(helpers.ExpandPath(args[0]))
	if err != nil {
		log.Fatal(err)
	}

	outlines := doc.Outlines()
	f, err := os.OpenFile(helpers.FeedsPath(),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, outline := range outlines {
		if len(outline.XMLURL) <= 0 {
			return
		}
		line := strings.Split(outline.XMLURL, "\n")
		url := []byte(fmt.Sprintf("%s\n", line[0]))
		f.Write(url)
	}
}

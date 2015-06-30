package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/gregf/podfetcher/src/helpers"
)

// Add adds a url to your feeds file
func Add(args []string) {
	if len(args) <= 0 {
		fmt.Println("podfetcher add <url>")
		return
	}

	file, err := os.OpenFile(helpers.FeedsPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	if _, err = file.WriteString(args[0] + "\n"); err != nil {
		log.Fatal(err)
	}
}

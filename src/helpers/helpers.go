package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	gohome "gopkg.in/caarlos0/gohome.v1"
)

var slash = string(os.PathSeparator)

const appName = "podfetcher"

// ConfigPath returns the path to the config file
func ConfigPath() (path string) {
	return filepath.Join(gohome.Config(appName), "config.yml")
}

// FeedsPath returns the path to the feeds path
func FeedsPath() (path string) {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), "feeds")
}

// ReadLines reads a file into an array of lines
func ReadLines(path string) (lines []string, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	lines = strings.Fields(string(content))
	return lines, err
}

// ReadFeeds reads a feeds file stripping out blank lines and comments
func ReadFeeds(path string) (lines []string, err error) {
	feeds, err := ReadLines(path)
	if err != nil {
		log.Fatal(err)
	}

	comment, err := regexp.Compile(`\A#`)
	if err != nil {
		log.Fatalf("comment %s\n", err)
	}
	for _, line := range feeds {
		if comment.Match([]byte(line)) {
			return
		}
		if len(line) <= 0 {
			return
		}
		lines = append(lines, line)
	}
	return lines, err
}

// ExpandPath expands a relative path to absolute path
func ExpandPath(path string) string {
	length := len(path)

	if length == 0 {
		return path
	}

	// replace env variables
	expandedPath := os.ExpandEnv(path)

	// replace ~ with $HOME
	if (length == 1 && path[0] == '~') || (length > 1 && path[:2] == "~/") {
		usr, _ := user.Current()

		expandedPath = strings.Replace(expandedPath, "~", usr.HomeDir, 1)
	} else if path[:1] == "~" {
		// replace ~user with their $HOME

		firstSlash := strings.Index(path, slash)
		if firstSlash < 0 {
			firstSlash = len(path)
		}

		if firstSlash > 1 {
			usr, err := user.Lookup(path[1:firstSlash])
			if err == nil {
				expandedPath = usr.HomeDir + path[firstSlash:]
			}
		}
	}

	// get an absolute path, ignoring errors
	if absPath, err := filepath.Abs(expandedPath); err == nil {
		expandedPath = absPath
	}

	// cleanup
	expandedPath = filepath.Clean(expandedPath)

	return expandedPath
}

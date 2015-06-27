package filter

import (
	"fmt"
	"strings"

	"github.com/gregf/podfetcher/Godeps/_workspace/src/github.com/spf13/viper"
)

// Run filters podcasts based on filters from config.yml
func Run(podcastTitle string, episodeTitle string) bool {
	filterName := fmt.Sprintf("filters.%s", podcastTitle)
	filters := viper.GetStringSlice(filterName)
	globalFilters := viper.GetStringSlice("filters.Global")

	// Global filters
	for _, filter := range globalFilters {
		if compare(episodeTitle, filter) {
			return true
		}
	}

	// Podcast specific filters
	for _, filter := range filters {
		if compare(episodeTitle, filter) {
			return true
		}
	}
	return false
}

func compare(episodeTitle, filter string) bool {
	e := strings.ToLower(episodeTitle)
	f := strings.ToLower(filter)
	if strings.Contains(e, f) {
		return true
	}
	return false
}

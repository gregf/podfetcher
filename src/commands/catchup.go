package commands

import (
	"github.com/gregf/podfetcher/src/database"
)

// CatchUp marks all episodes downloaded = true
func CatchUp() {
	database.CatchUp()
}

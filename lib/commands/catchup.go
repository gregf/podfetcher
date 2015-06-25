package commands

import (
	"github.com/gregf/podfetcher/lib/database"
)

// CatchUp marks all episodes downloaded = true
func CatchUp() {
	database.CatchUp()
}

package commands

import (
	"github.com/gregf/podfetcher/lib/database"
)

func CatchUp() {
	database.CatchUp()
}

package sqlite3

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/sqlite3/internal"
	"github.com/gofiber/storage/sqlite3/v2"
)

// SQLite3ConfigOverride is an interface for a Component must implements to adapt the [sqlite3.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type SQLite3ConfigOverride interface {
	SQLite3ConfigOverride(instance string, config *sqlite3.Config)
}

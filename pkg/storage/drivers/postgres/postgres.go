package postgres

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/postgres/internal"

	"github.com/gofiber/storage/postgres/v3"
)

// PostgresConfigOverride is an interface for a Component must implements to adapt the [postgres.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type PostgresConfigOverride interface {
	PostgresConfigOverride(instance string, config *postgres.Config)
}

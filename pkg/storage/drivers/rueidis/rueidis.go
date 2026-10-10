package rueidis

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis/internal"

	"github.com/gofiber/storage/rueidis"
)

// RueidisConfigOverride is an interface for a Component must implements to adapt the [rueidis.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type RueidisConfigOverride interface {
	RueidisConfigOverride(instance string, config *rueidis.Config)
}

package valkey

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey/internal"

	"github.com/gofiber/storage/valkey"
)

// ValkeyConfigOverride is an interface for a Component must implements to adapt the [valkey.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type ValkeyConfigOverride interface {
	ValkeyConfigOverride(instance string, config *valkey.Config)
}

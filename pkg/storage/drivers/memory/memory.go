package memory

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memory/internal" // init package

	"github.com/gofiber/storage/memory/v2"
)

// MemoryConfigOverride is an interface for a Component must implements to adapt the [memory.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type MemoryConfigOverride interface {
	MemoryConfigOverride(instance string, config *memory.Config)
}

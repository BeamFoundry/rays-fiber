package memcache

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memcache/internal"

	"github.com/gofiber/storage/memcache/v2"
)

// MemcacheConfigOverride is an interface for a Component must implements to adapt the [memcache.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type MemcacheConfigOverride interface {
	MemcacheConfigOverride(instance string, config *memcache.Config)
}

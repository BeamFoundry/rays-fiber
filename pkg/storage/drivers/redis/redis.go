package redis

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis/internal" // init package

	"github.com/gofiber/storage/redis/v3"
)

// RedisConfigOverride is an interface for a Component must implements to adapt the [redis.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type RedisConfigOverride interface {
	RedisConfigOverride(instance string, config *redis.Config)
}

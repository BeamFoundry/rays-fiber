package mongodb

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mongodb/internal"

	"github.com/gofiber/storage/mongodb/v2"
)

// MongoDBConfigOverride is an interface for a Component must implements to adapt the [mongodb.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type MongoDBConfigOverride interface {
	MongoDBConfigOverride(instance string, config *mongodb.Config)
}

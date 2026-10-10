package aerospike

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/aerospike/internal" // init package

	"github.com/gofiber/storage/aerospike"
)

// AerospikeConfigOverride is an interface for a Component must implements to adapt the [aerospike.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type AerospikeConfigOverride interface {
	AerospikeConfigOverride(instance string, config *aerospike.Config)
}

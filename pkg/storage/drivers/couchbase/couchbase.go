package couchbase

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/couchbase/internal" // init package

	"github.com/gofiber/storage/couchbase/v2"
)

// CouchbaseConfigOverride is an interface for a Component must implements to adapt the [couchbase.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type CouchbaseConfigOverride interface {
	CouchbaseConfigOverride(instance string, config *couchbase.Config)
}

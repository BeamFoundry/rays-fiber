package valkey

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey/internal" // init package

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	"github.com/gofiber/storage/valkey"
)

// ValkeyConfigOverride is an interface for a Component must implements to adapt the [valkey.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type ValkeyConfigOverride interface {
	ValkeyConfigOverride(instance string, config *valkey.Config)
}

// ValkeyClientTLSProvider is an interface a Component must implement to provide
// a [certs.ClientTLSProvider] for a Valkey storage instance.
type ValkeyClientTLSProvider interface {
	ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider
}

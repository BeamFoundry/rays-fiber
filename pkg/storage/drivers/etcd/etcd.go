package etcd

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd/internal" // init package

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	"github.com/gofiber/storage/etcd/v2"
)

// EtcdConfigOverride is an interface for a Component must implements to adapt the [etcd.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type EtcdConfigOverride interface {
	EtcdConfigOverride(instance string, config *etcd.Config)
}

// EtcdClientTLSProvider is an interface a Component must implement to provide
// a [certs.ClientTLSProvider] for an Etcd storage instance.
type EtcdClientTLSProvider interface {
	EtcdClientTLSProvider(instance string) certs.ClientTLSProvider
}

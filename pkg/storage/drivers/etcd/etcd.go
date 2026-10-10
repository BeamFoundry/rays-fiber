package etcd

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd/internal"

	"github.com/gofiber/storage/etcd/v2"
)

// EtcdConfigOverride is an interface for a Component must implements to adapt the [etcd.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type EtcdConfigOverride interface {
	EtcdConfigOverride(instance string, config *etcd.Config)
}

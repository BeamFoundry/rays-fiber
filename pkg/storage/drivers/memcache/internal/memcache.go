package internal

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/memcache/v2"

	"github.com/gofiber/fiber/v3"
)

type MemcacheConfigOverride interface {
	MemcacheConfigOverride(instance string, config *memcache.Config)
}
type memcacheConfig struct {
	servers string
	reset   bool
}

func (this memcacheConfig) Config() memcache.Config {
	config := memcache.ConfigDefault
	config.Servers = this.servers
	config.Reset = this.reset

	return config
}

// ---

type MemcacheConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	memcache map[string]memcacheConfig
}

func (this *MemcacheConfigurationProperties) Config(name string) memcache.Config {
	if c, ok := this.memcache[name]; ok {
		return c.Config()
	} else {
		return this.memcache["primary"].Config()
	}
}

var _, _ = (any)(&MemcacheConfigurationProperties{}).(Stereotype)

// ---

type memcacheStorageProvider struct {
	cfg      *MemcacheConfigurationProperties
	override Optional[MemcacheConfigOverride]
}

func (this *memcacheStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.MemcacheConfigOverride(name, &cfg)
	}
	return memcache.New(cfg)
}

type MemcacheAutoConfig struct {
	Configuration

	config   *MemcacheConfigurationProperties `@:"Inject"`
	override Optional[MemcacheConfigOverride] `@:"Inject"`
}

func (this *MemcacheAutoConfig) Memcache(_ struct {
	Qualifier `Name:"memcache"`
}) StorageEngine {
	return &memcacheStorageProvider{this.config, this.override}
}

func (this *MemcacheAutoConfig) MemcacheDefault(
	that struct {
		Qualifier             `Name:"memcache"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"memcache"`
		provider              StorageEngine `@:"Inject,Qualifier{memcache}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&MemcacheAutoConfig{}).(Stereotype)

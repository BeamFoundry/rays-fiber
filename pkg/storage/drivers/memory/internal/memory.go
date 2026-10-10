package memory

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/memory/v2"

	"github.com/gofiber/fiber/v3"
)

type MemoryConfigOverride interface {
	MemoryConfigOverride(instance string, config *memory.Config)
}

// ---

type memoryConfig struct {
	gcInternal time.Duration
}

func (this memoryConfig) Config() memory.Config {
	config := memory.ConfigDefault
	config.GCInterval = this.gcInternal
	return config
}

// ---

type MemoryConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	memory map[string]memoryConfig
}

var _, _ = any(&MemoryConfigurationProperties{}).(Stereotype)

func (this MemoryConfigurationProperties) Config(name string) memory.Config {
	c, ok := this.memory[name]
	if !ok {
		c = this.memory["primary"]
	}
	return c.Config()
}

type memoryStorageProvider struct {
	cfg      *MemoryConfigurationProperties
	override Optional[MemoryConfigOverride]
}

func (this *memoryStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.MemoryConfigOverride(name, &cfg)
	}
	return memory.New(cfg)
}

// ---

type MemoryAutoConfig struct {
	Configuration
	cfg      *MemoryConfigurationProperties `@:"Inject"`
	override Optional[MemoryConfigOverride] `@:"Inject"`
}

func (this *MemoryAutoConfig) Memory(
	_ struct {
		Qualifier `Name:"memory"`
	},
) StorageEngine {
	return &memoryStorageProvider{this.cfg, this.override}
}

func (this *MemoryAutoConfig) MemoryDefault(
	that struct {
		Qualifier             `Name:"memory"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"memory"`
		provider              StorageEngine `@:"Inject,Qualifier{memory}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&MemoryAutoConfig{}).(Stereotype)

package internal

import (
	"strings"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	. "github.com/BeamFoundry/rays/pkg/core"

	"github.com/gofiber/fiber/v3"
)

type storageProviderInstance struct {
	Provider string
	Instance string
}

type StorageProviderService struct {
	Service
	//
	primary   fiber.Storage            `@:"Inject"`
	providers map[string]StorageEngine `@:"Inject"`
	//
	cache map[storageProviderInstance]fiber.Storage
}

var _ StorageProvider = &StorageProviderService{}

var _, _ = any(&StorageProviderService{}).(Stereotype)

func (this StorageProviderService) New() *StorageProviderService {
	this.cache = map[storageProviderInstance]fiber.Storage{}
	return &this
}

func (this *StorageProviderService) Provide(name string) fiber.Storage {
	if name == "default" {
		return this.primary
	}

	c := strings.Split(name, ":")
	cfg := storageProviderInstance{}

	cfg.Provider = c[0]

	if len(c) > 1 {
		cfg.Instance = c[1]
	} else {
		cfg.Instance = "primary"
	}

	if s, ok := this.cache[cfg]; ok {
		return s
	} else if p, ok := this.providers[cfg.Provider]; ok {
		st := p.Storage(cfg.Instance)
		this.cache[cfg] = st
		return st
	}
	return this.primary
}

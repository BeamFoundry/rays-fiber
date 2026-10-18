package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// CacheNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cache], to indicates when to skip the creation of a new cache entry.
// If an entry already exists, it will be used.
//
// Please see the [cache] middleware configuration documentation.
//
// [cache]: https://docs.gofiber.io/middleware/cache#config
type CacheNext interface {
	CacheNext() func(fiber.Ctx) bool
}

// CacheInvalidator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cache], to indicates a function that is executed before checking the cache entry.
// It can be used to invalidate the existing cache manually by returning true.
//
// Please see the [cache] middleware configuration documentation.
//
// [cache]: https://docs.gofiber.io/middleware/cache#config
type CacheInvalidator interface {
	CacheInvalidator() func(fiber.Ctx) bool
}

// CacheKeyGenerator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cache], to generate custom keys. The HTTP method and a key-format version
// are partitioned internally by the middleware.
//
// Please see the [cache] middleware configuration documentation.
//
// [cache]: https://docs.gofiber.io/middleware/cache#config
type CacheKeyGenerator interface {
	CacheKeyGenerator() func(fiber.Ctx) string
}

// CacheExpirationGenerator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cache], to generate custom expiration keys based on the request.
//
// Please see the [cache] middleware configuration documentation.
//
// [cache]: https://docs.gofiber.io/middleware/cache#config
type CacheExpirationGenerator interface {
	CacheExpirationGenerator() func(fiber.Ctx, *cache.Config) time.Duration
}

// CacheConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Cache], to review the [cache.Config] that will be configured by [Cache]
//
// Please see the [cache] middleware configuration documentation.
//
// [cache]: https://docs.gofiber.io/middleware/cache/#config
type CacheConfig interface {
	CacheConfig(config *cache.Config)
}

type cacheConfigurationProperties struct {
	storage               *string
	cacheHeader           string
	methods               []string
	keyHeaders            []string
	expiration            time.Duration
	maxBytes              uint
	disableCacheControl   bool
	disableQueryKeys      bool
	disableValueRedaction bool
	disableVaryHeaders    bool
	storeResponseHeaders  bool
}

type cacheConfig struct {
	cache cacheConfigurationProperties
}

var _, _ = any(&FiberMwConfig[cacheConfig]{}).(Stereotype)

func (this *cacheConfig) Config() cache.Config {
	config := cache.ConfigDefault
	config.Expiration = this.cache.expiration
	config.CacheHeader = this.cache.cacheHeader
	config.KeyHeaders = this.cache.keyHeaders
	config.DisableCacheControl = this.cache.disableCacheControl
	config.DisableQueryKeys = this.cache.disableQueryKeys
	config.DisableValueRedaction = this.cache.disableValueRedaction
	config.DisableVaryHeaders = this.cache.disableVaryHeaders
	config.StoreResponseHeaders = this.cache.storeResponseHeaders
	config.MaxBytes = this.cache.maxBytes
	config.Methods = this.cache.methods
	return config
}

// [Cache] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Cache
//	}
//
//	func (this *ApiController) CacheConfig(config *cache.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [CacheNext], [CacheInvalidator], [CacheKeyGenerator] and [CacheExpirationGenerator].
//
// Also, the parent container may implements [CacheConfig] to receive a pointer to the [cache.Config]
// that will be used.
//
// Please see those interface definitions and the [cache] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//	        cache:
//	          expiration: 5m
//	          cacheHeader: "X-Cache"
//	          disablelCacheControl: false
//	          disableQueryKeys: false
//	          disableValueRedaction: false
//	          disableVaryHeaders: false
//	          keyHeaders: ["accept", "accept-encoding", "accept-language"]
//	          keyCookies: []
//	          storage: null
//	          storeResponseHeaders: false
//	          maxBytes: 0
//	          methods: [ "GET", "HEAD" ]
//
// [cache]: https://docs.gofiber.io/middleware/cache/#config
type Cache struct {
	config          *FiberMwConfig[cacheConfig] `@:"Inject"`
	storageProvider StorageProvider             `@:"Inject"`
	cfg             *cache.Config
}

func (this *Cache) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if cfg.cache.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.cache.storage)
	}
	if next, ok := parent.(CacheNext); ok {
		this.cfg.Next = next.CacheNext()
	}
	if invalidator, ok := parent.(CacheInvalidator); ok {
		this.cfg.CacheInvalidator = invalidator.CacheInvalidator()
	}
	if keyGenerator, ok := parent.(CacheKeyGenerator); ok {
		this.cfg.KeyGenerator = keyGenerator.CacheKeyGenerator()
	}
	if expirationGenerator, ok := parent.(CacheExpirationGenerator); ok {
		this.cfg.ExpirationGenerator = expirationGenerator.CacheExpirationGenerator()
	}
	if config, ok := parent.(CacheConfig); ok {
		config.CacheConfig(this.cfg)
	}
	return nil
}

func (this *Cache) Handler() fiber.Handler {
	return cache.New(*this.cfg)
}

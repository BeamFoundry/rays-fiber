---
title: "Cache"
weight: 2
---

# Cache Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/cache](https://docs.gofiber.io/api/middleware/cache)*

> **Storage Requirement:** This middleware uses a persistent storage backend. If `storage: null`, it defaults to the primary engine (usually `memory`). To configure a database like Redis, see [Storage Drivers](/docs/6-storage/).


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Cache` middleware will automatically bind them:

* `CacheNext() func(fiber.Ctx) bool`
* `CacheInvalidator() func(fiber.Ctx) bool`
* `CacheKeyGenerator() func(fiber.Ctx) string`
* `CacheExpirationGenerator() func(fiber.Ctx, *cache.Config) time.Duration`
* `CacheConfig(config *cache.Config)`

### Example: Configuring `Cache`

Here is an example of a controller embedding the middleware and implementing the `CacheConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/cache"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyCacheController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Cache
}

// 2. Implement the configuration interface
func (this *MyCacheController) CacheConfig(config *cache.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyCacheController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.cache`.

```yaml
rays:
  fiber:
    servers:
      primary:
        cache:
          expiration: 5m
          cacheHeader: "X-Cache"
          disablelCacheControl: false
          disableQueryKeys: false
          disableValueRedaction: false
          disableVaryHeaders: false
          keyHeaders: ["accept", "accept-encoding", "accept-language"]
          keyCookies: []
          storage: null
          storeResponseHeaders: false
          maxBytes: 0
          methods: [ "GET", "HEAD" ]
```


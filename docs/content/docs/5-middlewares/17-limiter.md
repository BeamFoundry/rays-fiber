---
title: "Limiter"
weight: 17
---

# Limiter Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/limiter](https://docs.gofiber.io/api/middleware/limiter)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Limiter` middleware will automatically bind them:

* `LimiterNext() func(fiber.Ctx) bool`
* `LimiterMaxFunc() func(c fiber.Ctx) int`
* `LimiterKeyGenerator() func(fiber.Ctx) string`
* `LimiterLimitReached() fiber.Handler`
* `LimiterConfig(config *limiter.Config)`

### Example: Configuring `Limiter`

Here is an example of a controller embedding the middleware and implementing the `LimiterConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/limiter"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyLimiterController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Limiter
}

// 2. Implement the configuration interface
func (this *MyLimiterController) LimiterConfig(config *limiter.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyLimiterController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.limiter`.

```yaml
rays:
  fiber:
    servers:
      primary:
        limiter:
          handler: fixed
          maxConn: 5 # Renamed from max
          expiration: 1m
          skipFailedRequests: false
          skipSuccessfulRequests: false
```


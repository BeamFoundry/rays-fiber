---
title: "HealthCheck"
weight: 12
---

# HealthCheck Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/healthcheck](https://docs.gofiber.io/api/middleware/healthcheck)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `HealthCheck` middleware will automatically bind them:

* `HealthCheckNext() func(fiber.Ctx) bool`
* `HealthCheckProbe() func(fiber.Ctx) bool`
* `HealthCheckConfig(config *healthcheck.Config)`

### Example: Configuring `HealthCheck`

Here is an example of a controller embedding the middleware and implementing the `HealthCheckConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/healthcheck"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyHealthCheckController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    HealthCheck
}

// 2. Implement the configuration interface
func (this *MyHealthCheckController) HealthCheckConfig(config *healthcheck.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyHealthCheckController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.healthcheck`.

```yaml
rays:
  fiber:
    servers:
      primary:
        healthcheck:
          format: "TEXT" # Possible values : TEXT, JSON, XML, MSGPACK, CBOR
```


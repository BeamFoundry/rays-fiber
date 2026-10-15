---
title: "Timeout"
weight: 28
---

# Timeout Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/timeout](https://docs.gofiber.io/api/middleware/timeout)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Timeout` middleware will automatically bind them:

* `TimeoutNext() func(fiber.Ctx) bool`
* `TimeoutOnTimeout() fiber.Handler`
* `TimeoutErrors() []error`
* `TimeoutConfig(config *timeout.Config)`

### Example: Configuring `Timeout`

Here is an example of a controller embedding the middleware and implementing the `TimeoutConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/timeout"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyTimeoutController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Timeout
}

// 2. Implement the configuration interface
func (this *MyTimeoutController) TimeoutConfig(config *timeout.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyTimeoutController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.timeout`.

```yaml
rays:
  fiber:
    servers:
      primary:
        timeout:
          timeout: 10s
```


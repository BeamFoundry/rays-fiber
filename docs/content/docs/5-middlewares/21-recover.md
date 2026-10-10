---
title: "Recover"
weight: 21
---

# Recover Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/recover](https://docs.gofiber.io/api/middleware/recover)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Recover` middleware will automatically bind them:

* `RecoverNext() func(fiber.Ctx) bool`
* `RecoverPanicHandler() func(c fiber.Ctx, e any) error`
* `RecoverStackTraceHandler() func(c fiber.Ctx, e any)`
* `RecoverConfig(config *recover.Config)`

### Example: Configuring `Recover`

Here is an example of a controller embedding the middleware and implementing the `RecoverConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/recover"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyRecoverController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Recover
}

// 2. Implement the configuration interface
func (this *MyRecoverController) RecoverConfig(config *recover.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyRecoverController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.recover`.

```yaml
rays:
  fiber:
    servers:
      primary:
        recover:
          enableStackTrace: false
```


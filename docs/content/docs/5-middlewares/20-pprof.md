---
title: "Pprof"
weight: 20
---

# Pprof Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/pprof](https://docs.gofiber.io/api/middleware/pprof)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Pprof` middleware will automatically bind them:

* `PprofNext() func(fiber.Ctx) bool`
* `PprofConfig(config *pprof.Config)`

### Example: Configuring `Pprof`

Here is an example of a controller embedding the middleware and implementing the `PprofConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/pprof"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyPprofController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Pprof
}

// 2. Implement the configuration interface
func (this *MyPprofController) PprofConfig(config *pprof.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyPprofController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.pprof`.

```yaml
rays:
  fiber:
    servers:
      primary:
        pprof:
          prefix: ""
```


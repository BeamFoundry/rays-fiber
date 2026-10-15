---
title: "Static"
weight: 26
---

# Static Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/static](https://docs.gofiber.io/api/middleware/static)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Static` middleware will automatically bind them:

* `StaticNext() func(fiber.Ctx) bool`
* `StaticFS() fs.FS`
* `StaticModifyResponse() fiber.Handler`
* `StaticNotFoundHandler() fiber.Handler`
* `StaticConfig(config *static.Config)`

### Example: Configuring `Static`

Here is an example of a controller embedding the middleware and implementing the `StaticConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/static"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyStaticController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Static
}

// 2. Implement the configuration interface
func (this *MyStaticController) StaticConfig(config *static.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyStaticController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.static`.

```yaml
rays:
  fiber:
    servers:
      primary:
        static:
          indexNames: ["index.html"]
          cacheDuration: 10s
          maxAge: 0
          compress: false
          byteRange: false
          browse: false
          download: false
```


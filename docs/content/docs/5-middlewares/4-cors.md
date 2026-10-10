---
title: "Cors"
weight: 4
---

# Cors Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/cors](https://docs.gofiber.io/api/middleware/cors)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Cors` middleware will automatically bind them:

* `CorsNext() func(fiber.Ctx) bool`
* `CorsAllowOriginsFunc() func(origin string) bool`
* `CorsConfig(config *cors.Config)`

### Example: Configuring `Cors`

Here is an example of a controller embedding the middleware and implementing the `CorsConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/cors"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyCorsController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Cors
}

// 2. Implement the configuration interface
func (this *MyCorsController) CorsConfig(config *cors.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyCorsController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.cors`.

```yaml
rays:
  fiber:
    servers:
      primary:
        cors:
          allowOrigins: ["*"]
          allowMethods: ["GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"]
          allowHeaders: []
          allowCredentials: false
          allowPrivateNetwork: false
          disableValueRedaction: false
          exposeHeaders: []
          maxAge: 0
```


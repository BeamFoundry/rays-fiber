---
title: "Etag"
weight: 9
---

# Etag Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/etag](https://docs.gofiber.io/api/middleware/etag)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Etag` middleware will automatically bind them:

* `EtagNext() func(fiber.Ctx) bool`
* `EtagConfig(config *etag.Config)`

### Example: Configuring `Etag`

Here is an example of a controller embedding the middleware and implementing the `EtagConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/etag"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyEtagController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Etag
}

// 2. Implement the configuration interface
func (this *MyEtagController) EtagConfig(config *etag.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyEtagController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.etag`.

```yaml
rays:
  fiber:
    servers:
      primary:
        etag:
          weak: false
```


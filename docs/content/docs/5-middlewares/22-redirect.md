---
title: "Redirect"
weight: 22
---

# Redirect Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/redirect](https://docs.gofiber.io/api/middleware/redirect)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Redirect` middleware will automatically bind them:

* `RedirectNext() func(fiber.Ctx) bool`
* `RedirectConfig(config *redirect.Config)`

### Example: Configuring `Redirect`

Here is an example of a controller embedding the middleware and implementing the `RedirectConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/redirect"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyRedirectController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Redirect
}

// 2. Implement the configuration interface
func (this *MyRedirectController) RedirectConfig(config *redirect.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyRedirectController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.redirect`.

```yaml
rays:
  fiber:
    servers:
      primary:
        redirect:
          rules:
            /api/*: null # null values will be ignored.
          statusCode: 302
```


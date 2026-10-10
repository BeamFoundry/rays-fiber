---
title: "Rewrite"
weight: 24
---

# Rewrite Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/rewrite](https://docs.gofiber.io/api/middleware/rewrite)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Rewrite` middleware will automatically bind them:

* `RewriteNext() func(fiber.Ctx) bool`
* `RewriteConfig(config *rewrite.Config)`

### Example: Configuring `Rewrite`

Here is an example of a controller embedding the middleware and implementing the `RewriteConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/rewrite"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyRewriteController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Rewrite
}

// 2. Implement the configuration interface
func (this *MyRewriteController) RewriteConfig(config *rewrite.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyRewriteController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.rewrite`.

```yaml
rays:
  fiber:
    servers:
      primary:
        rewrite:
          rules:
            /api/*: null # null values will be ignored.
```


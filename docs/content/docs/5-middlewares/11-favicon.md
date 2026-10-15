---
title: "Favicon"
weight: 11
---

# Favicon Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/favicon](https://docs.gofiber.io/api/middleware/favicon)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Favicon` middleware will automatically bind them:

* `FaviconNext() func(fiber.Ctx) bool`
* `FaviconFileSystem() fs.FS`
* `FaviconConfig(config *favicon.Config)`

### Example: Configuring `Favicon`

Here is an example of a controller embedding the middleware and implementing the `FaviconConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/favicon"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyFaviconController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Favicon
}

// 2. Implement the configuration interface
func (this *MyFaviconController) FaviconConfig(config *favicon.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyFaviconController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.favicon`.

```yaml
rays:
  fiber:
    servers:
      primary:
        favicon:
          file: ""
          url: "/favicon.ico"
          cacheControl: "public, max-age=31536000"
```


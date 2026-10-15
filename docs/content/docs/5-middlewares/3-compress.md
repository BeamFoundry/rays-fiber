---
title: "Compress"
weight: 3
---

# Compress Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/compress](https://docs.gofiber.io/api/middleware/compress)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Compress` middleware will automatically bind them:

* `CompressNext() func(fiber.Ctx) bool`
* `CompressConfig(config *compress.Config)`

### Example: Configuring `Compress`

Here is an example of a controller embedding the middleware and implementing the `CompressConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/compress"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyCompressController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Compress
}

// 2. Implement the configuration interface
func (this *MyCompressController) CompressConfig(config *compress.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyCompressController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.compress`.

```yaml
rays:
  fiber:
    servers:
      primary:
        compress:
          level: "default" # Case Unsensitive string. Possible values: Default, Disabled, Speed, BestSpeed, Compression, BestCompression.
```


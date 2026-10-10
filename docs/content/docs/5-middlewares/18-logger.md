---
title: "Logger"
weight: 18
---

# Logger Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/logger](https://docs.gofiber.io/api/middleware/logger)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Logger` middleware will automatically bind them:

* `LoggerNext() func(fiber.Ctx) bool`
* `LoggerSkip() func(fiber.Ctx) bool`
* `LoggerDone() func(c fiber.Ctx, logString []byte)`
* `LoggerCustomTags() map[string]logger.LogFunc`
* `LoggerBeforeHandlerFunc() func(*logger.Config)`
* `LoggerLoggerFunc() func(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error`
* `LoggerConfig(config *logger.Config)`

### Example: Configuring `Logger`

Here is an example of a controller embedding the middleware and implementing the `LoggerConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/logger"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyLoggerController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Logger
}

// 2. Implement the configuration interface
func (this *MyLoggerController) LoggerConfig(config *logger.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyLoggerController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.logger`.

```yaml
rays:
  fiber:
    servers:
      primary:
        logger:
          format: DefaultFormat
          timeFormat: 15:04:05
          timeZone: Local
          timeInterval: 500ms
          disableColors: false
          forceColors: false
```


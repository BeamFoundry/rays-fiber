---
title: "RequestID"
weight: 23
---

# RequestID Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/requestid](https://docs.gofiber.io/api/middleware/requestid)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `RequestID` middleware will automatically bind them:

* `RequestIDNext() func(fiber.Ctx) bool`
* `RequestIDGenerator() func() string`
* `RequestIDConfig(config *requestid.Config)`

### Example: Configuring `RequestID`

Here is an example of a controller embedding the middleware and implementing the `RequestIDConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/requestid"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyRequestIDController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    RequestID
}

// 2. Implement the configuration interface
func (this *MyRequestIDController) RequestIDConfig(config *requestid.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyRequestIDController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.requestid`.

```yaml
rays:
  fiber:
    servers:
      primary:
        requestid:
          header: "X-Request-ID"
```


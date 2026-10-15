---
title: "HostAuthorization"
weight: 14
---

# HostAuthorization Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/hostauthorization](https://docs.gofiber.io/api/middleware/hostauthorization)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `HostAuthorization` middleware will automatically bind them:

* `HostAuthorizationNext() func(fiber.Ctx) bool`
* `HostAuthorizationAllowedHostsFunc() func(host string) bool`
* `HostAuthorizationErrorHandler() fiber.ErrorHandler`
* `HostAuthorizationConfig(config *hostauthorization.Config)`

### Example: Configuring `HostAuthorization`

Here is an example of a controller embedding the middleware and implementing the `HostAuthorizationConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/hostauthorization"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyHostAuthorizationController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    HostAuthorization
}

// 2. Implement the configuration interface
func (this *MyHostAuthorizationController) HostAuthorizationConfig(config *hostauthorization.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyHostAuthorizationController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.hostauthorization`.

```yaml
rays:
  fiber:
    servers:
      primary:
        hostauthorization:
          allowedHosts: []
```


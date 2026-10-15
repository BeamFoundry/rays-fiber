---
title: "BasicAuth"
weight: 1
---

# BasicAuth Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/basicauth](https://docs.gofiber.io/api/middleware/basicauth)*

> **⚠️ Web Security Note:** This middleware will be **ignored** if used on a server that has an `HttpSecurity` configured. See [Web Security Integration](/docs/3-controllers/2-security-integration/#auto-disabled-middlewares) for details.


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `BasicAuth` middleware will automatically bind them:

* `BasicAuthNext() func(fiber.Ctx) bool`
* `BasicAuthAuthorizer() func(string, string, fiber.Ctx) bool`
* `BasicAuthUnauthorized() fiber.Handler`
* `BasicAuthBadRequest() fiber.Handler`
* `BasicAuthConfig(config *basicauth.Config)`

### Example: Configuring `BasicAuth`

Here is an example of a controller embedding the middleware and implementing the `BasicAuthConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/basicauth"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyBasicAuthController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    BasicAuth
}

// 2. Implement the configuration interface
func (this *MyBasicAuthController) BasicAuthConfig(config *basicauth.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}
```

## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.basicauth`.

```yaml
rays:
  fiber:
    servers:
      primary:
        basicauth:
          charset: "UTF-8"
          headerLimit: 8192
          users: # dictionary of username/password
            admin: null # Note: null passworded users will be ignored.
          realm: "Restricted"
```


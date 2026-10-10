---
title: "KeyAuth"
weight: 16
---

# KeyAuth Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/keyauth](https://docs.gofiber.io/api/middleware/keyauth)*

> **⚠️ Web Security Note:** This middleware will be **ignored** if used on a server that has an `HttpSecurity` configured. See [Web Security Integration](/docs/3-controllers/2-security-integration/#auto-disabled-middlewares) for details.


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `KeyAuth` middleware will automatically bind them:

* `KeyAuthNext() func(fiber.Ctx) bool`
* `KeyAuthSuccessHandler() fiber.Handler`
* `KeyAuthErrorHandler() fiber.ErrorHandler`
* `KeyAuthValidator() func(fiber.Ctx, string) (bool, error)`
* `KeyAuthExtractor() extractors.Extractor`
* `KeyAuthConfig(config *keyauth.Config)`

### Example: Configuring `KeyAuth`

Here is an example of a controller embedding the middleware and implementing the `KeyAuthConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/keyauth"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyKeyAuthController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    KeyAuth
}

// 2. Implement the configuration interface
func (this *MyKeyAuthController) KeyAuthConfig(config *keyauth.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyKeyAuthController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.keyauth`.

```yaml
rays:
  fiber:
    servers:
      primary:
        keyauth:
          realm: Restricted
          challenge: ApiKey realm="Restricted"
          errorCode: "" # Allowed values are invalid_request, invalid_token, or insufficient_scope.
          errorDescription: ""
          errorURI: ""
          scope: ""
```


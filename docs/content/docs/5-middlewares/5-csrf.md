---
title: "Csrf"
weight: 5
---

# Csrf Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/csrf](https://docs.gofiber.io/api/middleware/csrf)*

> **⚠️ Web Security Note:** This middleware will be **ignored** if used on a server that has an `HttpSecurity` configured. See [Web Security Integration](/docs/3-controllers/2-security-integration/#auto-disabled-middlewares) for details.

> **⚠️ Execution Order:** Because this middleware can depend on the `Session` middleware to store its tokens, you must ensure that **`Session` is included before `Csrf`** in your struct embedding order so the session is initialized first.


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Csrf` middleware will automatically bind them:

* `CsrfNext() func(fiber.Ctx) bool`
* `CsrfExtractor() extractors.Extractor`
* `CsrfKeyGenerator() func() string`
* `CsrfErrorHandler() func(fiber.Ctx, error) error`
* `CsrfConfig(config *csrf.Config)`

### Example: Configuring `Csrf`

Here is an example of a controller embedding the middleware and implementing the `CsrfConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/csrf"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyCsrfController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Csrf
}

// 2. Implement the configuration interface
func (this *MyCsrfController) CsrfConfig(config *csrf.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyCsrfController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.csrf`.

```yaml
rays:
  fiber:
    servers:
      primary:
        csrf:
          cookieName: "csrf_"
          cookieDomain: ""
          cookiePath: ""
          cookieSameSite: "Lax"
          trustedOrigins: []
          idleTimeout: 30m
          cookieSecure: false
          cookieHTTPOnly: false
          cookieSessionOnly: false
          disableValueRedaction: false
          singleUseToken: false
```


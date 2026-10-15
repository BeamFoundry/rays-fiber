---
title: "Session"
weight: 25
---

# Session Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/session](https://docs.gofiber.io/api/middleware/session)*

> **⚠️ Web Security Note:** This middleware will be **ignored** if used on a server that has an `HttpSecurity` configured. See [Web Security Integration](/docs/3-controllers/2-security-integration/#auto-disabled-middlewares) for details.

> **Storage Requirement:** This middleware uses a persistent storage backend. If `storage: null`, it defaults to the primary engine (usually `memory`). To configure a database like Redis, see [Storage Drivers](/docs/6-storage/).


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Session` middleware will automatically bind them:

* `SessionNext() func(fiber.Ctx) bool`
* `SessionErrorHandler() func(fiber.Ctx, error)`
* `SessionKeyGenerator() func() string`
* `SessionExtractor() extractors.Extractor`
* `SessionConfig(config *session.Config)`
* `SessionSecurityConfig(server string, config *session.Config)`
* `SessionStore() *session.Store`

### Example: Configuring `Session`

Here is an example of a controller embedding the middleware and implementing the `SessionConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/session"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MySessionController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Session
}

// 2. Implement the configuration interface
func (this *MySessionController) SessionConfig(config *session.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MySessionController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.session`.

```yaml
rays:
  fiber:
    servers:
      primary:
        session:
          storage: null
          cookieDomain: ""
          cookiePath: ""
          cookieSameSite: "Lax"
          idleTimeout: 30m
          absoluteTimeout: 0s
          cookieSecure: false
          cookieHTTPOnly: false
          cookieSessionOnly: false
```


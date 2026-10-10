---
title: "Idempotency"
weight: 15
---

# Idempotency Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/idempotency](https://docs.gofiber.io/api/middleware/idempotency)*

> **Storage Requirement:** This middleware uses a persistent storage backend. If `storage: null`, it defaults to the primary engine (usually `memory`). To configure a database like Redis, see [Storage Drivers](/docs/6-storage/).


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Idempotency` middleware will automatically bind them:

* `IdempotencyNext() func(fiber.Ctx) bool`
* `IdempotencyKeyHeaderValidate() func(string) error`
* `IdempotencyLocker() idempotency.Locker`
* `IdempotencyConfig(config *idempotency.Config)`

### Example: Configuring `Idempotency`

Here is an example of a controller embedding the middleware and implementing the `IdempotencyConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/idempotency"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyIdempotencyController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Idempotency
}

// 2. Implement the configuration interface
func (this *MyIdempotencyController) IdempotencyConfig(config *idempotency.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyIdempotencyController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.idempotency`.

```yaml
rays:
  fiber:
    servers:
      primary:
        idempotency:
          lifetime: 30m
          keyHeader: "X-Idempotency-Key"
          keepResponseHeaders: []
          storage: null
```


---
title: "5. Middlewares"
weight: 5
bookCollapseSection: true
---

# Declarative Middlewares

In standard Go Fiber applications, middlewares are registered imperatively using `app.Use(middleware())`. This often leads to scattered routing logic where it's difficult to see at a glance which endpoints are protected by which middlewares.

Rays Fiber solves this by adapting most of the original Fiber middlewares into standard **structs**.

> **Important:** To use these middlewares, you should dot-import the middleware package so you can use the bare struct names:
> `import . "github.com/BeamFoundry/rays-fiber/pkg/middleware"`

## Applying Middlewares

You apply middlewares declaratively by simply **embedding them** into your controllers.

### 1. Controller-Level (Group)
By embedding the middleware struct directly into your `FiberController` struct, the middleware will execute for **every endpoint** defined in that controller.

```go
import (
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware" // Dot-import required!
)

type ApiController struct {
    FiberController
    RequestMapping `Path:"/api"`

    // Embed the middleware structs directly!
    Logger
    Cors
}
```

### 2. Endpoint-Level (Route)
By embedding the middleware into the anonymous struct used for the endpoint method, the middleware will execute **only for that specific route**.

```go
func (this *ApiController) GetSensitiveData(
    _ struct {
        RequestMapping `Method:"GET" Path:"/sensitive"`

        // This middleware only runs for this specific GET request
        RateLimiter
    },
    ctx *fiber.Ctx,
) error {
    return ctx.JSON(fiber.Map{"data": "top-secret"})
}
```

> **💡 Middleware Shadowing:** If a middleware is defined at the `FiberController` level, it applies to all routes within that group. However, if you *also* embed that exact same middleware type into a specific method's anonymous struct, the method-level middleware **shadows** the controller-level one. The framework handles it as if it were only defined on the method for that specific route, allowing you to seamlessly override global middleware behavior for a single endpoint.

## Programmatic Configuration via Interfaces

While most middleware settings are controlled globally via `application.yaml`, sometimes you need to configure a middleware dynamically at the controller level (e.g., providing a custom fallback function, or altering the config struct directly).

Rays Fiber achieves this via **InnerBeamPostConstruct** introspection. If your `FiberController` (or `FiberPathFilter`) implements specific configuration interfaces associated with the embedded middleware, the framework automatically invokes them.

### Example: Customizing `EarlyData`

If you embed the `EarlyData` middleware in your controller, you can dynamically configure it by implementing `EarlyDataConfig` directly on that same controller:

```go
type FastApiController struct {
    FiberController
    RequestMapping `Path:"/fast"`

    EarlyData // 1. Embed the middleware
}

// 2. Implement the associated configuration interface
func (this *FastApiController) EarlyDataConfig(config *earlydata.Config) {
    // 3. Mutate the specific config just for this controller!
    config.IsEarlyData = func(c *fiber.Ctx) bool {
        return c.Get("Early-Data") == "1"
    }
}
```

Check the documentation page for each individual middleware to see a complete list of the programmatic interfaces they support.

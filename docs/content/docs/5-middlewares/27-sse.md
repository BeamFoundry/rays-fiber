---
title: "SSE"
weight: 27
---

# SSE Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/middleware/sse](https://docs.gofiber.io/middleware/sse)*


## Programmatic Configuration Interfaces

> **⚠️ FiberPathFilter only:** Due to ti's nature this middleware can only be used on a FiberPathFilter with a RequestMapping containing
a "Path" and "Method" struct tag.

> **⚠️ REQUIRED:** The following method should be implemented

* `SSEHandler() func(ctx fiber.Ctx, stream *sse.Stream) error`

If your Controller implements any of the following methods, the `SSE` middleware will automatically bind them:

* `SSEOnClose() func(ctx fiber.Ctx, err error)`
* `SSEConfig(config *sse.Config)`

### Example: Configuring `SSE`

Here is an example of a controller embedding the middleware and implementing the `SSEConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/sse"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MySSEEventHandler struct {
    FiberPathFilter
    // 1. Include a RequestMapping with a Method and Path struct tags.
    RequestMapping `Method:"GET" Path:"/event"`

    // Z. Embed the middleware
    SSE
}

// 3. Implement the SSEHandler interface
func (this *MySSEEventHandler) SSEHandler() func(ctx fiber.Ctx, stream *sse.Stream) error {
    // Implements the event handler
}

var _, _ = any(&MySSEController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.sse`.

```yaml
rays:
  fiber:
    servers:
      primary:
        sse:
          retry: 0
          heartbeatInterval: 15s
          disableHeartbeat: false
```

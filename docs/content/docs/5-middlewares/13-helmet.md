---
title: "Helmet"
weight: 13
---

# Helmet Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/helmet](https://docs.gofiber.io/api/middleware/helmet)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Helmet` middleware will automatically bind them:

* `HelmetNext() func(fiber.Ctx) bool`
* `HelmetConfig(config *helmet.Config)`

### Example: Configuring `Helmet`

Here is an example of a controller embedding the middleware and implementing the `HelmetConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v2/middleware/helmet"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyHelmetController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Helmet
}

// 2. Implement the configuration interface
func (this *MyHelmetController) HelmetConfig(config *helmet.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyHelmetController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.helmet`.

```yaml
rays:
  fiber:
    servers:
      primary:
        helmet:
          xSSProtection: "0"
          contentTypeNosniff: "nosniff"
          xFrameOptions: "SAMEORIGIN"
          contentSecurityPolicy: ""
          referrerPolicy: "no-referrer"
          permissionPolicy: ""
          crossOriginEmbedderPolicy: "require-corp"
          crossOriginOpenerPolicy: "same-origin"
          crossOriginResourcePolicy: "same-origin"
          originAgentCluster: "?1"
          xDNSPrefetchControl: "off"
          xDownloadOptions: "noopen"
          xPermittedCrossDomain: "none"
          hstsMaxAge: 0
          hstsExcludeSubdomains: false
          cspReportOnly: false
          hstsPreloadEnabled: false
```


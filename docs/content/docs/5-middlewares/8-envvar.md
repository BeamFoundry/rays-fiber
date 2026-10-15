---
title: "EnvVar"
weight: 8
---

# EnvVar Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/envvar](https://docs.gofiber.io/api/middleware/envvar)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `EnvVar` middleware will automatically bind them:

* `EnvVarConfig(config *envvar.Config)`

### Example: Configuring `EnvVar`

Here is an example of a controller embedding the middleware and implementing the `EnvVarConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/envvar"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyEnvVarController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    EnvVar
}

// 2. Implement the configuration interface
func (this *MyEnvVarController) EnvVarConfig(config *envvar.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyEnvVarController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.envvar`.

```yaml
rays:
  fiber:
    servers:
      primary:
        envvar:
          exportVars: {"HOME":null} # Using null as a map value will not register it. Use "" to have an empty value
```


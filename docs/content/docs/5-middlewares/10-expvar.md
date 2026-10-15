---
title: "ExpVar"
weight: 10
---

# ExpVar Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/expvar](https://docs.gofiber.io/api/middleware/expvar)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `ExpVar` middleware will automatically bind them:

* `ExpVarNext() func(fiber.Ctx) bool`
* `ExpVarConfig(config *expvar.Config)`

### Example: Configuring `ExpVar`

Here is an example of a controller embedding the middleware and implementing the `ExpVarConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/expvar"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyExpVarController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    ExpVar
}

// 2. Implement the configuration interface
func (this *MyExpVarController) ExpVarConfig(config *expvar.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyExpVarController{}).(Stereotype)
```


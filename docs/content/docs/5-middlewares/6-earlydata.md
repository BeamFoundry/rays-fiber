---
title: "EarlyData"
weight: 6
---

# EarlyData Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/earlydata](https://docs.gofiber.io/api/middleware/earlydata)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `EarlyData` middleware will automatically bind them:

* `EarlyDataNext() func(fiber.Ctx) bool`
* `EarlyDataIsEarlyData() func(c fiber.Ctx) bool`
* `EarlyDataAllowEarlyData() func(c fiber.Ctx) bool`
* `EarlyDataError() error`
* `EarlyDataConfig(config *earlydata.Config)`

### Example: Configuring `EarlyData`

Here is an example of a controller embedding the middleware and implementing the `EarlyDataConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/earlydata"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyEarlyDataController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    EarlyData
}

// 2. Implement the configuration interface
func (this *MyEarlyDataController) EarlyDataConfig(config *earlydata.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyEarlyDataController{}).(Stereotype)
```


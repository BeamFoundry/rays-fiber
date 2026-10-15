---
title: "Paginate"
weight: 19
---

# Paginate Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/paginate](https://docs.gofiber.io/api/middleware/paginate)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `Paginate` middleware will automatically bind them:

* `PaginateNext() func(fiber.Ctx) bool`
* `PaginateConfig(config *paginate.Config)`

### Example: Configuring `Paginate`

Here is an example of a controller embedding the middleware and implementing the `PaginateConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/paginate"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyPaginateController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Paginate
}

// 2. Implement the configuration interface
func (this *MyPaginateController) PaginateConfig(config *paginate.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyPaginateController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.paginate`.

```yaml
rays:
  fiber:
    servers:
      primary:
        paginate:
          pageKey: "page"
          limitKey: "limit"
          sortKey: ""
          defaultSort: "id"
          cursorKey: "cursor"
          offsetKey: "offset"
          cursorParam: ""
          allowedSorts: []
          defaultPage: 1
          defaultLimit: 10
          maxLimit: 100
```


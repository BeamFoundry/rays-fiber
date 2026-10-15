---
title: "EncryptCookie"
weight: 7
---

# EncryptCookie Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/api/middleware/encryptcookie](https://docs.gofiber.io/api/middleware/encryptcookie)*


## Programmatic Configuration Interfaces

If your Controller implements any of the following methods, the `EncryptCookie` middleware will automatically bind them:

* `EncryptCookieNext() func(fiber.Ctx) bool`
* `EncryptCookieEncryptor() func(name, decryptedString, key string) (string, error)`
* `EncryptCookieDecryptor() func(name, encryptedString, key string) (string, error)`
* `EncryptCookieConfig(config *encryptcookie.Config)`

### Example: Configuring `EncryptCookie`

Here is an example of a controller embedding the middleware and implementing the `EncryptCookieConfig` interface to mutate the underlying Fiber configuration:

```go
package controllers

import (
    "github.com/gofiber/fiber/v3/middleware/encryptcookie"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware"
)

type MyEncryptCookieController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    EncryptCookie
}

// 2. Implement the configuration interface
func (this *MyEncryptCookieController) EncryptCookieConfig(config *encryptcookie.Config) {
    // Mutate the config directly for this controller
    // config.SomeProperty = SomeValue
}

var _, _ = any(&MyEncryptCookieController{}).(Stereotype)
```
## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.encryptcookie`.

```yaml
rays:
  fiber:
    servers:
      primary:
        encryptcookie:
          key: null
          except: []
```


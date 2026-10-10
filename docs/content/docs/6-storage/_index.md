---
title: "6. Storage Drivers"
weight: 6
bookCollapseSection: true
---

# Storage Drivers

Many Fiber middlewares (such as Cache, Session, and Idempotency) rely on temporary or persistent storage engines to function. Rays Fiber provides a standardized way to configure multiple storage backends and cleanly inject them into your middlewares.

## Configuration Syntax

When configuring a middleware in YAML, you will often see a `storage` property.

* If `storage: null`, the middleware will use the server's default storage engine (usually in-memory).
* If you want to use a specific storage engine, you provide a string in the format `<driver>:<instance>` (e.g., `redis:primary` or `mongodb:cache-cluster`).

By default, the framework defines the `"memory"` driver as the default storage engine.

```yaml
rays:
  fiber:
    storage:
      primary: memory # Specify which Storage engine to default to.
```

## Importing Drivers

The default `memory` driver is included out-of-the-box in the `rays-fiber` module.

All other drivers are separated into their own packages to keep your binary size small. To use a database driver like Redis or Postgres, you must import its specific module in your `go.mod` and underscore-import it in your code to trigger its registration.

Modules are located under:
`github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/<driver>`

**Example:**
```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis"
)
```

## Multiple storage instances

Storage drivers are not coupled to web server instance in order to share them, if needed.

To define another instance you have to duplicate it's configuration and adapt it.

Use `CopyDefaultConfig` (or `CopyDefaultConfigForProfile`) in your `main.go` before calling `Initialize()`:

```go
func main() {
    ctx := rays.ApplicationContext()

    // Duplicate the "primary" redis storage configuration.
    ctx.CopyDefaultConfig("rays.fiber.storage.redis.primary", "rays.fiber.storage.redis.secondary")

    if err := ctx.Initialize(); err != nil {
        err.LogErrorWithStackTrace()
        os.Exit(1)
    }

    ctx.Run()
}
```

Now you can change the defaults in your YAML configuration.

```yaml
rays:
  fiber:
    storage:
      redis:
        secondary:
          port: 6380
```

## Available Drivers

Rays Fiber supports (or is currently implementing) a wide variety of storage drivers. Click on a driver in the menu to see its default configuration and properties.

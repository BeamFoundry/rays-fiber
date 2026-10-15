---
title: "4. Configuration"
weight: 4
---

# Application Configuration

Rays Fiber binds configuration natively to the Rays IoC ecosystem. Instead of imperatively defining server timeouts, body limits, or network preferences in code, you define them externally in your `application.yaml`.

## The Default Configuration

At module initialization (via `init()`), Rays Fiber automatically loads a comprehensive default configuration. It defines a default server named `"primary"` that listens on port `3000`.

Here is the exact internal YAML map that the framework uses by default:

```yaml
rays:
  fiber:
    servers:
      primary:
        address: ":3000"
        config:
          bodyLimit: 4194304 # 4 MB
          compressedFileSuffixes: {gzip: ".fiber.gz", br: ".fiber.br", zstd: ".fiber.zst"}
          concurrency: 262144
          readBufferSize: 4096
          requestMethods: [GET,HEAD,POST,PUT,DELETE,CONNECT,OPTIONS,TRACE,PATCH]
          writeBufferSize: 4096
          appName: ""
          disableDefaultContentType: false
          disableDefaultDate: false
          disableHeaderNormalizing: false
          disableKeepalive: false
          disablePreParseMultipartForm: false
          enableIPValidation: false
          enableSplittingOnParsers: false
          getOnly: false
          idleTimeout: 0
          immutable: false
          passLocalsToViews: false
          readTimeout: 0
          reduceMemoryUsage: false
          proxyHeader: ""
          serverHeader: ""
          streamRequestBody: false
          strictRouting: false
          trustProxy: false
          unescapePath: false
          viewsLayout: ""
          writeTimeout: 0
        listener:
          listenerNetwork: "tcp4"
          shutdownTimeout: 10s
          tlsMinVersion: 0x0303 # tls.VersionTLS12
          disableStartupMessage: false
          enablePrefork: false
          enablePrintRoutes: false
```

## Overriding Configurations

Because Rays utilizes `ConfigurationProperties`, you only need to declare the values you want to change in your application's actual `application.yaml`. Everything else will gracefully fall back to the defaults shown above.

For example, to change the default port to `8080` and enable Fiber's prefork feature for higher throughput, your `application.yaml` would look like this:

```yaml
rays:
  fiber:
    servers:
      primary:
        address: ":8080"
        listener:
          enablePrefork: true
```

---

## Multiple Servers

Rays Fiber supports defining and running multiple distinct web servers simultaneously (e.g., an internal API on port 8080 and a public API on port 443).

To do this, you first instruct the Rays `ApplicationContext` to duplicate the default fiber configuration tree under a new server name during bootstrap. Then, you can override its specific settings in your `application.yaml`.

### 1. Bootstrapping the Secondary Server
Use `CopyDefaultConfig` (or `CopyDefaultConfigForProfile`) in your `main.go` before calling `Initialize()`:

```go
func main() {
    ctx := rays.ApplicationContext()

    // Duplicate the "primary" default configuration tree into a new server named "admin-api"
    ctx.CopyDefaultConfig("rays.fiber.servers.primary", "rays.fiber.servers.admin-api")

    if err := ctx.Initialize(); err != nil {
        err.LogErrorWithStackTrace()
        os.Exit(1)
    }

    ctx.Run()
}
```

### 2. Configuring the Secondary Server
Now you can target this new server in your YAML:

```yaml
rays:
  fiber:
    servers:
      primary:
        address: ":8080" # Public API
      admin-api:
        address: ":9090" # Internal Admin API
```

*Remember: To mount a Controller to this new server, embed the `core.Server` marker with the tag `Name:"admin-api"` in your controller struct (see the Controllers documentation).*

---

## Programmatic Configuration

Sometimes YAML isn't enough. If you need to inject complex Go objects into the Fiber engine (like custom error handlers, JSON encoders, or specific network interfaces), you can programmatically hook into the Fiber configuration using the `FiberConfig` and `FiberListenConfig` interfaces.

By implementing these on a Beam, you gain direct access to the underlying `fiber.Config` and `fiber.ListenConfig` structs before the server starts.

```go
package configurations

import (
    "github.com/gofiber/fiber/v3"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
)

type MyCustomFiberConfig struct {
    Component
    // Server `Name:"admin-api"` // Optional: Use this to apply configs to a specific server
}

// Intercepts the fiber.Config before the App is instantiated
func (this *MyCustomFiberConfig) FiberConfig(cfg *fiber.Config) {
    cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "A custom global error occurred: " + err.Error(),
        })
    }
}

// Intercepts the fiber.ListenConfig before the server starts listening
func (this *MyCustomFiberConfig) FiberListenConfig(cfg *fiber.ListenConfig) {
    cfg.DisableStartupMessage = true
}

```

## Lifecycle Management (`BeamRunner`)

Internally, the `fiber.App` instance is managed as a `BeamRunner` component.

This means that you never need to manually call `app.Listen()` or handle operating system signals to gracefully shut down the web server. When you execute `ctx.Run()` in your `main.go`, the Rays IoC container automatically spins up all configured listeners, and upon receiving an interrupt signal, safely tears them down adhering to the `shutdownTimeout` configuration.

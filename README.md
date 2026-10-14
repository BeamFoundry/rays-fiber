# BeamFoundry Rays Fiber 🚀

**This is an extract of the full documentation available at [rays.bdfy.dev](https://rays.bfdy.dev/rays-fiber/).**

**Rays Fiber** seamlessly integrates the popular, Express-inspired [Go Fiber](https://gofiber.io/)
web framework with the **Rays IoC Container**.

By bringing Fiber into the IoC ecosystem, you can build blazing-fast web APIs using clean,
object-oriented controllers, effortless dependency injection, and declarative configuration.

## Introduction and Architecture

While Go Fiber is incredibly fast and easy to use, large enterprise applications often struggle with state management,
dependency injection, and scattered route definitions across dozens of files.

**Rays Fiber** solves this by treating your web layer as standard IoC components.

### Core Benefits

1. **Injected Controllers:** Define your HTTP handlers as methods on a struct, allowing you to seamlessly inject repositories,
  services, and configurations using the Rays `@:"Inject"` tag.
2. **Centralized Lifecycle:** The HTTP server lifecycle (startup and graceful shutdown) is automatically managed by the Rays
  application context.
3. **Declarative Routing:** Define endpoint HTTP methods and paths directly via struct tags, keeping your routing logic exactly
  where the handler code lives.

## Quickstart and Setup

Integrating Fiber into your Rays application takes only a few lines of code.

### 📦 Installation

```bash
go get github.com/BeamFoundry/rays-fiber
```

### The Entrypoint

To bootstrap the web server, underscore-import the `rays-fiber` package in your `main.go`. This triggers the auto-configuration that provisions the underlying Fiber App into the container.

```go
package main

import (
    "os"
    "github.com/BeamFoundry/rays"

    // Auto-configures and starts the Fiber application
    _ "github.com/BeamFoundry/rays-fiber"
)

func main() {
    ctx := rays.ApplicationContext()

    if err := ctx.Initialize(); err != nil {
        err.LogErrorWithStackTrace()
        os.Exit(1)
    }

    // The container will automatically start the Fiber server
    // and listen for incoming HTTP requests.
    ctx.Run()
}
```

## Controllers and Routing

Rays Fiber drops the traditional, scattered `app.Get("/path", handler)` routing in favor of highly cohesive, object-oriented controllers.

By utilizing marker interfaces and struct tags, you can declaratively map paths and HTTP methods directly to your Controller methods.

### The `FiberController` Stereotype

To create a controller, you must embed the `FiberController` interface (found in `rays-fiber/pkg/core`). Because `FiberController` derives from the base Rays `Controller` Stereotype, your struct instantly becomes an IoC-managed Beam. This means you can inject services and repositories into it using `@:"Inject"`.

### Declarative Routing with `RequestMapping`

Rays uses the `RequestMapping` marker interface (from `github.com/BeamFoundry/rays/pkg/web`) combined with struct tags to build the routing tree.

#### 1. Group Routing (Struct Level)
By embedding `RequestMapping` in your controller struct and applying a `Path` tag, the framework automatically mounts all the controller's endpoints under a Fiber `Group()`.

> **💡 Optional Group Paths:** The struct-level `RequestMapping` (or its `Path` tag) is entirely optional. If you omit it, the controller behaves as if it were mounted to the root path (`Path:"/"`) to ease single URL definition. In this case, every endpoint method *must* provide its own explicit `Path` tag.

#### 2. Endpoint Routing (Method Level)
To map a specific controller method to an HTTP endpoint, the method's **first argument** must be an anonymous struct containing the `RequestMapping` marker. You apply `Method` and `Path` tags to this marker.

The **second argument** must then be the standard `*fiber.Ctx` pointer.

### Complete Example

Here is a full example demonstrating dependency injection, a controller-level route group, and multiple endpoint handlers.

```go
package controllers

import (
    "github.com/gofiber/fiber/v2"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
)

type UserController struct {
    FiberController                    // 1. Mark as a Fiber Controller Beam
    RequestMapping  `Path:"/api/users"` // 2. Maps to fiber.Group("/api/users")

    // 3. Standard IoC Dependency Injection
    userService *UserService `@:"Inject"`
}

// GetUser handles GET /api/users/:id
func (this *UserController) GetUser(
    // 4. Endpoint mapping via anonymous struct argument
    _ struct {
        RequestMapping `Method:"GET" Path:"/:id"`
    },
    // 5. The Fiber context
    ctx *fiber.Ctx,
) error {

    id := ctx.Params("id")
    user := this.userService.FindByID(id)

    if user == nil {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
    }

    return ctx.JSON(user)
}

// CreateUser handles POST /api/users
func (this *UserController) CreateUser(
    _ struct {
        RequestMapping `Method:"POST" Path:"/"`
    },
    ctx *fiber.Ctx,
) error {

    // Custom business logic...
    return ctx.SendStatus(fiber.StatusCreated)
}

```

> **Behind the scenes:** When the application context initializes, Rays scans all Beams for the `FiberController` stereotype. It reads the struct-level tags to create the Fiber Group, uses reflection to parse the anonymous struct arguments on the methods, and wires everything into the underlying Fiber App natively!

### Multi-Server Controllers

Rays Fiber supports running multiple web servers simultaneously (e.g., a "primary" public API and a "secondary" internal admin API).

By default, all `FiberController` components are mounted to the `"primary"` server. If you want a controller to mount to a different server, simply embed the `Server` marker interface (from `github.com/BeamFoundry/rays/pkg/core`) into your struct and provide the target server's name via a struct tag:

```go
type AdminController struct {
    FiberController
    Server          `Name:"admin-api"`   // Binds this controller to the "admin-api" server
    RequestMapping  `Path:"/internal/admin"`
}
```

## Global and Path Filters

While embedding middlewares in `FiberController` structs is great for endpoint-specific logic, sometimes you need to apply middlewares across your entire application or to broad URL prefixes.

For this, Rays Fiber provides the `FiberPathFilter` stereotype.

### The `FiberPathFilter` Component

Unlike a Controller, a `FiberPathFilter` is based on the standard Rays `Component` stereotype. It does not define HTTP endpoints; its sole purpose is to register the middlewares embedded within it.

You can configure its scope using the `RequestMapping` struct tag:

#### 1. Global Scope (`app.Use`)
If you do not define a `RequestMapping` at all, the embedded middlewares are applied globally to every request on the server (equivalent to `app.Use(middlewares...)`).

```go
import (
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware" // Dot-import gives access to all middlewares
)

type GlobalSecurityFilter struct {
    FiberPathFilter

    // Applied globally to the "primary" server
    Cors
    Helmet
}
```

#### 2. Path Prefix Scope
If you define a `Path`, the middlewares act as a prefix filter (equivalent to `app.Use("/api", middlewares...)`).

```go
type ApiPathFilter struct {
    FiberPathFilter
    RequestMapping `Path:"/api"` // Only applies to /api/*

    Logger
}
```

#### 3. Method-Specific Pattern Scope
If you define *both* a `Method` and a `Path`, the underlying Fiber framework will bind it specifically (e.g., `app.Post("/api/*", middlewares...)`). **Note:** When using HTTP Methods on a PathFilter, you must use Fiber's wildcard syntax (`/*`) to capture sub-paths.

```go
type PostMetricsFilter struct {
    FiberPathFilter
    RequestMapping `Method:"POST" Path:"/api/*"`

    // Only runs on POST requests to /api/*
    Metrics
}
```

### Multi-Server Support

Just like Controllers, `FiberPathFilter` components apply to the `"primary"` server by default. You can easily bind them to secondary servers using the `Server` marker interface:

```go
type AdminGlobalFilter struct {
    FiberPathFilter
    Server `Name:"admin-api"` // Bind globally to the admin server

    Logger
}
```

## Web Security Integration

Rays Fiber is designed to deeply integrate with the **Rays Web Security** framework. When an `HttpSecurity` object is configured for a specific server, it assumes control over authentication and authorization.

### Binding `HttpSecurity` to Servers

By default, an `HttpSecurity` beam returned by your configuration targets the `"primary"` server.

To bind an `HttpSecurity` object to a secondary server (e.g., `"admin-api"`), you supply an anonymous struct as the *first argument* to your configuration method, applying the `Server` struct tag:

```go
func (this *SecurityConfig) AdminHttpSecurity(
    server struct {
        Server `Name:"admin-api"` // Secures the admin-api server
    },
) *HttpSecurity {
    return (&HttpSecurityBuilder{}).
        // ... rules ...
        Build()
}
```

### Auto-Disabled Middlewares

When Rays Fiber detects that an `HttpSecurity` object has been bound to a server, it intelligently **ignores** the following standard Fiber middlewares if you attempt to embed them directly inside your controllers:

* `Csrf`
* `Session`
* `BasicAuth`
* `KeyAuth`

*Why?* Because `rays-web-security` handles CSRF protection, session policies (`SessionCreationPolicy`), and authentication internally at a global filter level. Applying them at the controller level would cause conflicts.

### Customizing Security `Csrf` and `Session`

Even though `rays-web-security` orchestrates CSRF and Sessions, you may still want to tweak their underlying Fiber configurations (e.g., changing the cookie name or timeout).

To do this, create a component that implements the `CsrfSecurityConfig` or `SessionSecurityConfig` interfaces:

```go
package configurations

import (
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/fiber/v2/middleware/csrf"
    . "github.com/BeamFoundry/rays/pkg/core"
)

type MySecurityTweaks struct {
    Component
}

// Customizes the Fiber Session middleware utilized by rays-web-security
func (this *MySecurityTweaks) SessionSecurityConfig(server string, config *session.Config) {
    if server == "primary" {
        config.CookieSecure = true // Force HTTPS cookies
    }
}

// Customizes the Fiber CSRF middleware utilized by rays-web-security
func (this *MySecurityTweaks) CsrfSecurityConfig(server string, config *csrf.Config) {
    if server == "primary" {
        config.CookieName = "my_custom_csrf"
    }
}

```

## Application Configuration

Rays Fiber binds configuration natively to the Rays IoC ecosystem. Instead of imperatively defining server timeouts, body limits, or network preferences in code, you define them externally in your `application.yaml`.

### The Default Configuration

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

### Overriding Configurations

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

### Multiple Servers

Rays Fiber supports defining and running multiple distinct web servers simultaneously (e.g., an internal API on port 8080 and a public API on port 443).

To do this, you first instruct the Rays `ApplicationContext` to duplicate the default fiber configuration tree under a new server name during bootstrap. Then, you can override its specific settings in your `application.yaml`.

#### 1. Bootstrapping the Secondary Server
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

#### 2. Configuring the Secondary Server
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

### Programmatic Configuration

Sometimes YAML isn't enough. If you need to inject complex Go objects into the Fiber engine (like custom error handlers, JSON encoders, or specific network interfaces), you can programmatically hook into the Fiber configuration using the `FiberConfig` and `FiberListenConfig` interfaces.

By implementing these on a Beam, you gain direct access to the underlying `fiber.Config` and `fiber.ListenConfig` structs before the server starts.

```go
package configurations

import (
    "github.com/gofiber/fiber/v2"
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

### Lifecycle Management (`BeamRunner`)

Internally, the `fiber.App` instance is managed as a `BeamRunner` component.

This means that you never need to manually call `app.Listen()` or handle operating system signals to gracefully shut down the web server. When you execute `ctx.Run()` in your `main.go`, the Rays IoC container automatically spins up all configured listeners, and upon receiving an interrupt signal, safely tears them down adhering to the `shutdownTimeout` configuration.


## Declarative Middlewares

In standard Go Fiber applications, middlewares are registered imperatively using `app.Use(middleware())`. This often leads to scattered routing logic where it's difficult to see at a glance which endpoints are protected by which middlewares.

Rays Fiber solves this by adapting most of the original Fiber middlewares into standard **structs**.

> **Important:** To use these middlewares, you should dot-import the middleware package so you can use the bare struct names:
> `import . "github.com/BeamFoundry/rays-fiber/pkg/middleware"`

### Applying Middlewares

You apply middlewares declaratively by simply **embedding them** into your controllers.

#### 1. Controller-Level (Group)
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

#### 2. Endpoint-Level (Route)
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

### Programmatic Configuration via Interfaces

While most middleware settings are controlled globally via `application.yaml`, sometimes you need to configure a middleware dynamically at the controller level (e.g., providing a custom fallback function, or altering the config struct directly).

Rays Fiber achieves this via **InnerBeamPostConstruct** introspection. If your `FiberController` (or `FiberPathFilter`) implements specific configuration interfaces associated with the embedded middleware, the framework automatically invokes them.

#### Example: Customizing `EarlyData`

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

## Storage Drivers

Many Fiber middlewares (such as Cache, Session, and Idempotency) rely on temporary or persistent storage engines to function. Rays Fiber provides a standardized way to configure multiple storage backends and cleanly inject them into your middlewares.

### Configuration Syntax

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

### Importing Drivers

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

### Multiple storage instances

Storage drivers are not coupled to web server instance in order to share them between, if needed.

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


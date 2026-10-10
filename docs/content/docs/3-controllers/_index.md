---
title: "3. Controllers & Routing"
weight: 3
bookCollapseSection: false
---

# Controllers & Routing

Rays Fiber drops the traditional, scattered `app.Get("/path", handler)` routing in favor of highly cohesive, object-oriented controllers.

By utilizing marker interfaces and struct tags, you can declaratively map paths and HTTP methods directly to your Controller methods.

## The `FiberController` Stereotype

To create a controller, you must embed the `FiberController` interface (found in `rays-fiber/pkg/core`). Because `FiberController` derives from the base Rays `Controller` Stereotype, your struct instantly becomes an IoC-managed Beam. This means you can inject services and repositories into it using `@:"Inject"`.

## Declarative Routing with `RequestMapping`

Rays uses the `RequestMapping` marker interface (from `github.com/BeamFoundry/rays/pkg/web`) combined with struct tags to build the routing tree.

### 1. Group Routing (Struct Level)
By embedding `RequestMapping` in your controller struct and applying a `Path` tag, the framework automatically mounts all the controller's endpoints under a Fiber `Group()`.

> **💡 Optional Group Paths:** The struct-level `RequestMapping` (or its `Path` tag) is entirely optional. If you omit it, the controller behaves as if it were mounted to the root path (`Path:"/"`) to ease single URL definition. In this case, every endpoint method *must* provide its own explicit `Path` tag.

### 2. Endpoint Routing (Method Level)
To map a specific controller method to an HTTP endpoint, the method's **first argument** must be an anonymous struct containing the `RequestMapping` marker. You apply `Method` and `Path` tags to this marker.

The **second argument** must then be the standard `*fiber.Ctx` pointer.

## Complete Example

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

## Multi-Server Controllers

Rays Fiber supports running multiple web servers simultaneously (e.g., a "primary" public API and a "secondary" internal admin API).

By default, all `FiberController` components are mounted to the `"primary"` server. If you want a controller to mount to a different server, simply embed the `Server` marker interface (from `github.com/BeamFoundry/rays/pkg/core`) into your struct and provide the target server's name via a struct tag:

```go
type AdminController struct {
    FiberController
    Server          `Name:"admin-api"`   // Binds this controller to the "admin-api" server
    RequestMapping  `Path:"/internal/admin"`
}
```

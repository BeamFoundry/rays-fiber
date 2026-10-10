---
title: "Global & Path Filters"
weight: 1
---

# Global & Path Filters

While embedding middlewares in `FiberController` structs is great for endpoint-specific logic, sometimes you need to apply middlewares across your entire application or to broad URL prefixes.

For this, Rays Fiber provides the `FiberPathFilter` stereotype.

## The `FiberPathFilter` Component

Unlike a Controller, a `FiberPathFilter` is based on the standard Rays `Component` stereotype. It does not define HTTP endpoints; its sole purpose is to register the middlewares embedded within it.

You can configure its scope using the `RequestMapping` struct tag:

### 1. Global Scope (`app.Use`)
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

### 2. Path Prefix Scope
If you define a `Path`, the middlewares act as a prefix filter (equivalent to `app.Use("/api", middlewares...)`).

```go
type ApiPathFilter struct {
    FiberPathFilter
    RequestMapping `Path:"/api"` // Only applies to /api/*

    Logger
}
```

### 3. Method-Specific Pattern Scope
If you define *both* a `Method` and a `Path`, the underlying Fiber framework will bind it specifically (e.g., `app.Post("/api/*", middlewares...)`). **Note:** When using HTTP Methods on a PathFilter, you must use Fiber's wildcard syntax (`/*`) to capture sub-paths.

```go
type PostMetricsFilter struct {
    FiberPathFilter
    RequestMapping `Method:"POST" Path:"/api/*"`

    // Only runs on POST requests to /api/*
    Metrics
}
```

## Multi-Server Support

Just like Controllers, `FiberPathFilter` components apply to the `"primary"` server by default. You can easily bind them to secondary servers using the `Server` marker interface:

```go
type AdminGlobalFilter struct {
    FiberPathFilter
    Server `Name:"admin-api"` // Bind globally to the admin server

    Logger
}
```

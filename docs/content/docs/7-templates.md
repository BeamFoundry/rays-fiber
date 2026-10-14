---
title: "7. HTML Templates"
weight: 7
aliases:
  - /rays-fiber/docs/7-template-html/
  - /rays-fiber/docs/template-html/
---

# HTML Template Engine

Rays Fiber includes a revamped HTML template engine in [`pkg/template/html`](https://github.com/BeamFoundry/rays-fiber/tree/main/pkg/template/html),
based on Go's standard [`html/template`](https://pkg.go.dev/html/template) package.

While inspired by the original Fiber template module ([`github.com/gofiber/template/html`](https://github.com/gofiber/template)),
this package modernizes layout and partial handling to adhere strictly to idiomatic Go template semantics.

---

## 🌟 Why the Revamped Version?

In the original `gofiber/template/html` engine, layout inheritance relied on a proprietary, non-standard `{{embed}}` tag placed inside the layout file.
The engine intercepted `{{embed}}` and stitched template strings together behind the scenes.

While simple for trivial layouts, this design had several limitations:
- **Single placeholder limitation:** You could not easily define multiple customizable blocks (such as page title, custom `<head>` stylesheets, main content, and page-specific scripts).
- **Non-standard Go syntax:** `{{embed}}` is not part of standard Go [`html/template`](https://pkg.go.dev/html/template).
- **Brittle nesting:** Multi-level layout hierarchies were complex and prone to rendering issues.

### The Rays Fiber Approach

The revamped `pkg/template/html` completely removes the need for `{{embed}}`. Instead, it leverages Go's native template inheritance constructs:

1. **Layouts** define pluggable regions using standard Go `{{block "<name>" .}}...{{end}}` actions.
2. **Views** define specific content for those regions using standard Go `{{define "<name>"}}...{{end}}` blocks.
3. **Views** explicitly invoke the layout to render using `{{template "<layout_name>" .}}` (or the dynamic `{{include .Template .}}` helper).

This adheres 100% to idiomatic Go template principles while supporting multiple blocks, nested layout hierarchies, and thread-safe caching.

---

## 🚀 Installation & Import

Import the template package directly:

```go
import "github.com/BeamFoundry/rays-fiber/pkg/template/html"
```

---

## 📐 Template Architecture & Syntax

### 1. Directory Structure

A typical template directory in your project:

```text
views/
├── layouts/
│   ├── main.html
│   └── nested/
│       ├── base.html
│       └── admin.html
├── partials/
│   ├── header.html
│   └── footer.html
├── errors/
│   └── 404.html
└── index.html
```

> **Template Names:** Templates are identified by their relative path inside the views directory **without** the file extension. For example, `views/layouts/main.html` is referenced as `"layouts/main"`, and `views/partials/header.html` is referenced as `"partials/header"`.

---

### 2. Creating Layouts with `{{block}}`

In your layout (e.g. `views/layouts/main.html`), use `{{block "<name>" .}}...{{end}}` anywhere a view might supply custom content. You can also provide default content inside the block that views can optionally override:

```html
<!-- views/layouts/main.html -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{block "title" .}}Default Application Title{{end}}</title>

  <!-- Block for page-specific stylesheets or meta tags -->
  {{block "head" .}}{{end}}
</head>
<body>
  <!-- Standard partial inclusion -->
  {{template "partials/header" .}}

  <main class="container">
    <!-- Main content block -->
    {{block "content" .}}
      <p>Default content placeholder</p>
    {{end}}
  </main>

  <!-- Standard partial inclusion -->
  {{template "partials/footer" .}}

  <!-- Block for page-specific scripts -->
  {{block "scripts" .}}{{end}}
</body>
</html>
```

---

### 3. Creating Views with `{{define}}` and `{{template}}`

In your view template (e.g. `views/index.html`):
1. **Define** the content for each block using `{{define "<name>"}}...{{end}}`.
2. **Invoke** the target layout using `{{template "<layout_name>" .}}`.

```html
<!-- views/index.html -->

<!-- 1. Invoke the layout to render -->
{{template "layouts/main" .}}

<!-- 2. Supply page-specific title -->
{{define "title"}}{{.Title}} - My App{{end}}

<!-- 3. Supply optional head elements -->
{{define "head"}}
<link rel="stylesheet" href="/static/css/home.css">
{{end}}

<!-- 4. Supply main content -->
{{define "content"}}
  <div class="hero">
    <h1>{{.Title}}</h1>
    <p>{{.WelcomeMessage}}</p>
  </div>
{{end}}

<!-- 5. Supply page-specific scripts -->
{{define "scripts"}}
<script src="/static/js/home.js"></script>
{{end}}
```

---

### 4. Partials

Re-usable fragments (like navigation bars or footers) can be placed in separate files and included anywhere using the standard `{{template "<name>" .}}` directive:

```html
<!-- views/partials/header.html -->
<header>
  <nav>
    <a href="/">Home</a>
    <a href="/about">About</a>
  </nav>
</header>
```

```html
<!-- views/partials/footer.html -->
<footer>
  <p>&copy; 2026 My App. All rights reserved.</p>
</footer>
```

---

## 🌲 Multi-Level Nested Layouts

Hierarchical layout trees (e.g. Base HTML Shell → Admin Layout with Sidebar → Specific Dashboard Page) are natively supported:

### Base Layout (`views/layouts/nested/base.html`)

```html
<!-- views/layouts/nested/base.html -->
<!DOCTYPE html>
<html>
<head>
  <title>{{block "title" .}}Admin Console{{end}}</title>
</head>
<body>
  {{block "shell" .}}{{end}}
</body>
</html>
```

### Child Layout (`views/layouts/nested/admin.html`)

```html
<!-- views/layouts/nested/admin.html -->
{{template "layouts/nested/base" .}}

{{define "shell"}}
  <aside class="sidebar">
    <nav><!-- Admin navigation links --></nav>
  </aside>
  <main class="admin-body">
    {{block "content" .}}{{end}}
  </main>
{{end}}
```

### View (`views/admin.html`)

```html
<!-- views/admin.html -->
{{template "layouts/nested/admin" .}}

{{define "title"}}Dashboard - Admin Console{{end}}

{{define "content"}}
  <h1>Admin Dashboard</h1>
  <p>Welcome, {{.User}}!</p>
{{end}}
```

---

## 🧰 Built-in Functions & Dynamic Layouts

The `pkg/template/html` engine registers helpful built-in template functions:

### 1. `include`

The `include` function dynamically executes a named template and returns the rendered output as safe `template.HTML`:

```html
{{include "partials/header" .}}
```

This is particularly useful when the layout or partial name is passed dynamically in the view data (e.g., `.Template`):

```html
<!-- views/flexible.html -->
{{define "content"}}
  <h1>{{.Title}}</h1>
  <p>Rendered with dynamic layout selection.</p>
{{end}}

{{if .Template}}
  {{include .Template .}}
{{else}}
  {{block "content" .}}{{end}}
{{end}}
```

If `.Template` is set to `"layouts/main"`, `{{include .Template .}}` will render the layout, which in turn renders
the `"content"` block. If no layout is provided, it falls back to rendering `{{block "content" .}}` directly.

### 2. `toHTML`

Converts a raw string into unescaped `template.HTML`:

```html
{{printf "<h1>%s</h1>" .Title | toHTML}}
```

---

## ⚙️ Configuration in Rays Fiber

In Rays, applications configure Fiber components using `FiberConfig` on an IoC `Component` beam.

### Using Local Filesystem (`html.New`)

```go
package config

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"

    "github.com/BeamFoundry/rays-fiber/pkg/template/html"
    "github.com/gofiber/fiber/v3"
)

type ViewConfiguration struct {
    Component
}

var _ FiberConfig = &ViewConfiguration{}

func (v *ViewConfiguration) FiberConfig(cfg *fiber.Config) {
    // Initialize engine pointing to the views directory and template extension
    engine := html.New("./views", ".html")

    // Optional: Enable hot reloading for development
    // engine.Reload(true)

    // Optional: Register custom helper functions
    engine.AddFunc("isAdmin", func(user string) bool {
        return user == "admin"
    })

    // Assign to Fiber's configuration
    cfg.Views = engine
    cfg.PassLocalsToViews = true
}
```

### Using Embedded Filesystem (`html.NewFileSystem`)

For single-binary deployments, bundle your views with `embed.FS` and `html.NewFileSystem`:

```go
package config

import (
    "embed"
    "net/http"

    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"

    "github.com/BeamFoundry/rays-fiber/pkg/template/html"
    "github.com/gofiber/fiber/v3"
)

//go:embed views/*
var viewsFS embed.FS

type EmbeddedViewConfiguration struct {
    Component
}

var _ FiberConfig = &EmbeddedViewConfiguration{}

func (v *EmbeddedViewConfiguration) FiberConfig(cfg *fiber.Config) {
    engine := html.NewFileSystem(http.FS(viewsFS), ".html")
    cfg.Views = engine
}
```

---

## 🎮 Rendering in Controllers

Render views inside your `FiberController` methods using `ctx.Render`:

```go
package controllers

import (
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    "github.com/gofiber/fiber/v3"
)

type HomeController struct {
    FiberController
    RequestMapping `Path:"/"`
}

func (h *HomeController) Index(
    _ struct {
        RequestMapping `Method:"GET" Path:"/"`
    },
    ctx fiber.Ctx,
) error {
    // Renders "views/index.html" with "views/layouts/main.html"
    return ctx.Render("index", fiber.Map{
        "Title":          "Home Page",
        "WelcomeMessage": "Welcome to Rays Fiber with modern Go templates!",
    }, "layouts/main")
}

func (h *HomeController) Admin(
    _ struct {
        RequestMapping `Method:"GET" Path:"/admin"`
    },
    ctx fiber.Ctx,
) error {
    // Nested layout example: "views/admin.html" with nested layouts
    return ctx.Render("admin", fiber.Map{
        "User": "admin",
    }, "layouts/nested/main", "layouts/nested/base")
}
```

> **Note on Layout Arguments:** When passing layout names to `ctx.Render(view, binding, layouts...)`, the engine isolates and compiles
the template tree specifically for that layout combination.
The view template then invokes the layout via `{{template "layouts/main" .}}` (or `{{include .Template .}}`).

---

## 🛠️ Engine Options & API Reference

| Method | Description |
| :--- | :--- |
| `html.New(directory, extension string) *Engine` | Creates a new template engine scanning a filesystem path (e.g. `./views`, `".html"`). |
| `html.NewFileSystem(fs http.FileSystem, extension string) *Engine` | Creates a new template engine using an `http.FileSystem` (such as `http.FS(embedFS)`). |
| `engine.AddFunc(name string, fn any)` | Registers a custom helper function for templates. |
| `engine.AddFuncMap(m map[string]any)` | Registers multiple helper functions at once via a function map. |
| `engine.Delims(left, right string)` | Customizes template delimiters (defaults to `"{{"` and `"}}"`). |
| `engine.Reload(enabled bool)` | When set to `true`, reloads and reparses templates from the filesystem on every render (useful during development). |
| `engine.Load() error` | Parses all templates in the configured directory or filesystem. |
| `engine.Render(out io.Writer, name string, binding any, layout ...string) error` | Renders a template with binding data and optional layouts. |

---

## 📊 Comparison with Original `gofiber/template/html`

| Feature | Original `gofiber/template/html` | Rays Fiber `pkg/template/html` |
| :--- | :--- | :--- |
| **Layout Injection** | Proprietary `{{embed}}` marker | Standard Go `{{block "<name>" .}}{{end}}` |
| **Layout Invocation** | Implicit string replacement | Standard Go `{{template "layout" .}}` / `{{include ...}}` |
| **Multiple Content Blocks** | ❌ Not supported (only one `{{embed}}`) | ✅ Fully supported (`"title"`, `"head"`, `"content"`, `"scripts"`, etc.) |
| **Nested Layouts** | ⚠️ Fragile with multiple levels | ✅ Standard Go template composition |
| **Standard Go Tooling** | ⚠️ Linters/formatters may flag `{{embed}}` | ✅ 100% compliant with standard Go `html/template` |
| **Built-in Helpers** | Basic | Includes `include` and `toHTML` |
| **Thread-Safe Caching** | Basic | Cached per view and layout permutation with mutex safety |

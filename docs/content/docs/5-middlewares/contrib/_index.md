---
title: "Contrib"
weight: 29
bookCollapseSection: true
---

# Contrib Middlewares

Third-party Fiber middlewares from [gofiber/contrib](https://github.com/gofiber/contrib) are adapted into their own Go sub-modules under `github.com/BeamFoundry/rays-fiber/pkg/middleware/contrib/<name>` so that their external dependencies are only pulled in when you explicitly use them.

To use a contrib middleware, install its sub-module and dot-import it into your controller or path filter:

```go
import (
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware/contrib/otel"
)
```

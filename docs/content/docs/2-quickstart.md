---
title: "2. Quickstart & Setup"
weight: 2
---

# Quickstart & Setup

Integrating Fiber into your Rays application takes only a few lines of code.

## 📦 Installation

```bash
go get github.com/BeamFoundry/rays-fiber
```

## The Entrypoint

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

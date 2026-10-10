---
title: "SQLite3"
weight: 11
---

# SQLite3 Storage Driver

## 📦 Installation & Import

To use the SQLite3 driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/sqlite3
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/sqlite3"
)
```

## Programmatic Configuration (`SQLite3ConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `SQLite3ConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/sqlite3"
)

type SQLite3Customizer struct {
    Component
}

// SQLite3ConfigOverride customizes the configuration for any SQLite3 instance (e.g., "primary")
func (this *SQLite3Customizer) SQLite3ConfigOverride(instance string, config *sqlite3.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&SQLite3Customizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      sqlite3:
        primary:
          database: "file::memory:"
          table: fiber_storage
          reset: false
          gcInterval: 10s
          maxIdleConns: 100
          maxOpenConns: 100
          connMaxLifetime: 10s
```

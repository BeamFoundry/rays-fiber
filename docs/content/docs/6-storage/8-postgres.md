---
title: "Postgres"
weight: 8
---

# Postgres Storage Driver

## 📦 Installation & Import

To use the Postgres driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/postgres
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/postgres"
)
```

## Programmatic Configuration (`PostgresConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `PostgresConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/postgres"
)

type PostgresCustomizer struct {
    Component
}

// PostgresConfigOverride customizes the configuration for any Postgres instance (e.g., "primary")
func (this *PostgresCustomizer) PostgresConfigOverride(instance string, config *postgres.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&PostgresCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      postgres:
        primary:
          connectionURI: ""
          host: 127.0.0.1
          port: 5432
          username: ""
          password: ""
          database: fiber
          table: fiber_storage
          sslMode: disable
          reset: false
          gcInterval: 10s
```

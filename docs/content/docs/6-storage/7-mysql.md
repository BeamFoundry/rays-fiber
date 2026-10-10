---
title: "MySQL"
weight: 7
---

# MySQL Storage Driver

## 📦 Installation & Import

To use the MySQL driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mysql
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mysql"
)
```

## Programmatic Configuration (`MySQLConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `MySQLConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mysql"
)

type MySQLCustomizer struct {
    Component
}

// MySQLConfigOverride customizes the configuration for any MySQL instance (e.g., "primary")
func (this *MySQLCustomizer) MySQLConfigOverride(instance string, config *mysql.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&MySQLCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      mysql:
        primary:
          connectionURI: ""
          host: 127.0.0.1
          port: 3306
          username: ""
          password: ""
          database: fiber
          table: fiber_storage
          reset: false
          gcInterval: 10s
```

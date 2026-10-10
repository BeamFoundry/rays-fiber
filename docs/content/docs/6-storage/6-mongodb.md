---
title: "MongoDB"
weight: 6
---

# MongoDB Storage Driver

## 📦 Installation & Import

To use the MongoDB driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mongodb
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mongodb"
)
```

## Programmatic Configuration (`MongoDBConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `MongoDBConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mongodb"
)

type MongoDBCustomizer struct {
    Component
}

// MongoDBConfigOverride customizes the configuration for any MongoDB instance (e.g., "primary")
func (this *MongoDBCustomizer) MongoDBConfigOverride(instance string, config *mongodb.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&MongoDBCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      mongodb:
        primary:
          connectionURI: ""
          host: 127.0.0.1
          port: 27017
          username: ""
          password: ""
          database: "fiber"
          collection: "fiber_storage"
          reset: false
```

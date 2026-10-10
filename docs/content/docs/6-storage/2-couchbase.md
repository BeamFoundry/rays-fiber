---
title: "Couchbase"
weight: 2
---

# Couchbase Storage Driver

## 📦 Installation & Import

To use the Couchbase driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/couchbase
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/couchbase"
)
```

## Programmatic Configuration (`CouchbaseConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `CouchbaseConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/couchbase"
)

type CouchbaseCustomizer struct {
    Component
}

// CouchbaseConfigOverride customizes the configuration for any Couchbase instance (e.g., "primary")
func (this *CouchbaseCustomizer) CouchbaseConfigOverride(instance string, config *couchbase.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&CouchbaseCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      couchbase:
        primary:
          host: couchbase://127.0.0.1:8091
          username: "admin"
          password: "123456"
          bucket: fiber_storage
          connectionTimeout: 3s
          kvTimeout: 1s
```

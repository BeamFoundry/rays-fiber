---
title: "Etcd"
weight: 3
---

# Etcd Storage Driver

## 📦 Installation & Import

To use the Etcd driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd"
)
```

## Programmatic Configuration (`EtcdConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `EtcdConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd"
)

type EtcdCustomizer struct {
    Component
}

// EtcdConfigOverride customizes the configuration for any Etcd instance (e.g., "primary")
func (this *EtcdCustomizer) EtcdConfigOverride(instance string, config *etcd.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&EtcdCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      etcd:
        primary:
          endpoints: ["localhost:2379"]
          dialTimeout: 2s
          username: ""
          password: ""
```

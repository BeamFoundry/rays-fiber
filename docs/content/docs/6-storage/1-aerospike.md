---
title: "Aerospike"
weight: 1
---

# Aerospike Storage Driver

## 📦 Installation & Import

To use the Aerospike driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/aerospike
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/aerospike"
)
```

## Programmatic Configuration (`AerospikeConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `AerospikeConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/aerospike"
)

type AerospikeCustomizer struct {
    Component
}

// AerospikeConfigOverride customizes the configuration for any Aerospike instance (e.g., "primary")
func (this *AerospikeCustomizer) AerospikeConfigOverride(instance string, config *aerospike.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&AerospikeCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      aerospike:
        primary:
          hosts: ["localhost:3000"]
          namespace: "test"
          setName: "fiber"
          reset: false
          expiration: 1h
          schemaVersion: 1
          schemaDescription: "Default Fiber storage schema"
          forceSchemaUpdate: false
          initialConnectionTimeout: 10s
```

---
title: "Valkey"
weight: 12
---

# Valkey Storage Driver

## 📦 Installation & Import

To use the Valkey driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"
)
```

## Programmatic Configuration (`ValkeyConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `ValkeyConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"
)

type ValkeyCustomizer struct {
    Component
}

// ValkeyConfigOverride customizes the configuration for any Valkey instance (e.g., "primary")
func (this *ValkeyCustomizer) ValkeyConfigOverride(instance string, config *valkey.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&ValkeyCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      valkey:
        primary:
          username: ""
          password: ""
          clientName: ""
          url: ""
          selectDB: 0
          initAddress: ["127.0.0.1:6379"]
          cacheSizeEachConn: 134217728
          ringScaleEachConn: 10
          readBufferEachConn: 524288
          writeBufferEachConn: 524288
          blockingPoolSize: 1024
          pipelineMultiplex: 2
          disableRetry: false
          disableCache: false
          alwaysPipelining: true
          reset: false
          cacheTTL: 1m
```

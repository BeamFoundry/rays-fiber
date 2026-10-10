---
title: "Memcache"
weight: 4
---

# Memcache Storage Driver

## 📦 Installation & Import

To use the Memcache driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memcache
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memcache"
)
```

## Programmatic Configuration (`MemcacheConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `MemcacheConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memcache"
)

type MemcacheCustomizer struct {
    Component
}

// MemcacheConfigOverride customizes the configuration for any Memcache instance (e.g., "primary")
func (this *MemcacheCustomizer) MemcacheConfigOverride(instance string, config *memcache.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&MemcacheCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      memcache:
        primary:
          servers: localhost:11211
          reset: false
```

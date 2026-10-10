---
title: "Memory"
weight: 5
---

# Memory Storage Driver

## 📦 Installation & Import

The **Memory** storage driver is bundled natively in the `rays-fiber` module. No additional `go get` or underscore-import is needed.

## Programmatic Configuration (`MemoryConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `MemoryConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memory"
)

type MemoryCustomizer struct {
    Component
}

// MemoryConfigOverride customizes the configuration for any Memory instance (e.g., "primary")
func (this *MemoryCustomizer) MemoryConfigOverride(instance string, config *memory.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&MemoryCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      memory:
        primary:
          gcInternal: 10s
```

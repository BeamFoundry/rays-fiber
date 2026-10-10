---
title: "Redis"
weight: 9
---

# Redis Storage Driver

## 📦 Installation & Import

To use the Redis driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis"
)
```

## Programmatic Configuration (`RedisConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `RedisConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis"
)

type RedisCustomizer struct {
    Component
}

// RedisConfigOverride customizes the configuration for any Redis instance (e.g., "primary")
func (this *RedisCustomizer) RedisConfigOverride(instance string, config *redis.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&RedisCustomizer{}).(Stereotype)
```

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      redis:
        primary:
          host: 127.0.0.1
          port: 6379
          username: ""
          password: ""
          database: 0
          url: ""
          addrs: []
          masterName: ""
          clientName: ""
          sentinelUsername: ""
          sentinelPassword: ""
          reset: false
          poolSize: 10
          isClusterMode: false
```

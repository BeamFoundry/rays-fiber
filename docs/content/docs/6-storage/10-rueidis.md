---
title: "Rueidis"
weight: 10
---

# Rueidis Storage Driver

## 📦 Installation & Import

To use the Rueidis driver, first install its standalone package:

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis
```

Then underscore-import it in your application entrypoint (e.g. `main.go`):

```go
import (
    _ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"
)
```

## Programmatic Configuration (`RueidisConfigOverride`)

If you need to customize settings on the underlying Fiber storage configuration object programmatically (such as adding custom TLS contexts or driver-specific options), create a Rays `Component` that implements `RueidisConfigOverride`:

```go
package configurations

import (
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"
)

type RueidisCustomizer struct {
    Component
}

// RueidisConfigOverride customizes the configuration for any Rueidis instance (e.g., "primary")
func (this *RueidisCustomizer) RueidisConfigOverride(instance string, config *rueidis.Config) {
    if instance == "primary" {
        // Mutate configuration options prior to driver connection
    }
}

var _, _ = any(&RueidisCustomizer{}).(Stereotype)
```

## Client TLS Configuration (`RueidisClientTLSProvider`)

To configure TLS for Rueidis connections using `rays-web-security`, create a `Component` that implements `RueidisClientTLSProvider`. Its method receives the storage instance name and returns a `certs.ClientTLSProvider` (or `nil` to bypass TLS for that instance):

```go
package configurations

import (
    "crypto/tls"

    "github.com/BeamFoundry/rays-web-security/pkg/certs"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/lang"
)

type RueidisTLSProvider struct {
    Component
}

// RueidisClientTLSProvider returns a ClientTLSProvider for the specified instance
func (this *RueidisTLSProvider) RueidisClientTLSProvider(instance string) certs.ClientTLSProvider {
    if instance == "primary" {
        return this
    }
    return nil
}

// ProvideClientTLS provides the *tls.Config for the Rueidis connection
func (this *RueidisTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
    return &tls.Config{
        ServerName: "rueidis.example.com",
    }, nil
}

var _ RueidisClientTLSProvider = &RueidisTLSProvider{}
var _ certs.ClientTLSProvider = &RueidisTLSProvider{}
```

> **💡 Override Integration:** A `RueidisConfigOverride` component can also implement `RueidisClientTLSProvider` directly. In all cases, `RueidisConfigOverride` will still be invoked after TLS is set up, allowing you to inspect or further tune the final configuration.

## Default YAML Configuration

The framework provisions the `"primary"` instance of this driver with the following default settings. To customize it or configure copied instances, modify your `application.yaml`:

```yaml
rays:
  fiber:
    storage:
      rueidis:
        primary:
          username: ''
          password: ''
          clientName: ''
          url: ''
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
          disableAlwaysPipelining: false
          reset: false
          cacheTTL: 1m0s
```

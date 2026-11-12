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

## Client TLS Configuration (`ValkeyClientTLSProvider`)

To configure TLS for Valkey connections using `rays-web-security`, create a `Component` that implements `ValkeyClientTLSProvider`. Its method receives the storage instance name and returns a `certs.ClientTLSProvider` (or `nil` to bypass TLS for that instance):

```go
package configurations

import (
    "crypto/tls"

    "github.com/BeamFoundry/rays-web-security/pkg/certs"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/lang"
)

type ValkeyTLSProvider struct {
    Component
}

// ValkeyClientTLSProvider returns a ClientTLSProvider for the specified instance
func (this *ValkeyTLSProvider) ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider {
    if instance == "primary" {
        return this
    }
    return nil
}

// ProvideClientTLS provides the *tls.Config for the Valkey connection
func (this *ValkeyTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
    return &tls.Config{
        ServerName: "valkey.example.com",
    }, nil
}

var _ ValkeyClientTLSProvider = &ValkeyTLSProvider{}
var _ certs.ClientTLSProvider = &ValkeyTLSProvider{}
```

> **💡 Override Integration:** A `ValkeyConfigOverride` component can also implement `ValkeyClientTLSProvider` directly. In all cases, `ValkeyConfigOverride` will still be invoked after TLS is set up, allowing you to inspect or further tune the final configuration.

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

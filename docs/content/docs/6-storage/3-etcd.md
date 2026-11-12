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

## Client TLS Configuration (`EtcdClientTLSProvider`)

To configure TLS for Etcd connections using `rays-web-security`, create a `Component` that implements `EtcdClientTLSProvider`. Its method receives the storage instance name and returns a `certs.ClientTLSProvider` (or `nil` to bypass TLS for that instance):

```go
package configurations

import (
    "crypto/tls"

    "github.com/BeamFoundry/rays-web-security/pkg/certs"
    . "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/lang"
)

type EtcdTLSProvider struct {
    Component
}

// EtcdClientTLSProvider returns a ClientTLSProvider for the specified instance
func (this *EtcdTLSProvider) EtcdClientTLSProvider(instance string) certs.ClientTLSProvider {
    if instance == "primary" {
        return this
    }
    return nil
}

// ProvideClientTLS provides the *tls.Config for the Etcd connection
func (this *EtcdTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
    return &tls.Config{
        ServerName: "etcd.example.com",
    }, nil
}

var _ EtcdClientTLSProvider = &EtcdTLSProvider{}
var _ certs.ClientTLSProvider = &EtcdTLSProvider{}
```

> **💡 Override Integration:** An `EtcdConfigOverride` component can also implement `EtcdClientTLSProvider` directly. In all cases, `EtcdConfigOverride` will still be invoked after TLS is set up, allowing you to inspect or further tune the final configuration.

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

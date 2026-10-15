---
title: "Web Security Integration"
weight: 2
---

# Web Security Integration

Rays Fiber is designed to deeply integrate with the **Rays Web Security** framework. When an `HttpSecurity` object is configured for a specific server, it assumes control over authentication and authorization.

## Binding `HttpSecurity` to Servers

By default, an `HttpSecurity` beam returned by your configuration targets the `"primary"` server.

To bind an `HttpSecurity` object to a secondary server (e.g., `"admin-api"`), you supply an anonymous struct as the *first argument* to your configuration method, applying the `Server` struct tag:

```go
func (this *SecurityConfig) AdminHttpSecurity(
    server struct {
        Server `Name:"admin-api"` // Secures the admin-api server
    },
) *HttpSecurity {
    return (&HttpSecurityBuilder{}).
        // ... rules ...
        Build()
}
```

## Auto-Disabled Middlewares

When Rays Fiber detects that an `HttpSecurity` object has been bound to a server, it intelligently **ignores** the following standard Fiber middlewares if you attempt to embed them directly inside your controllers:

* `Csrf`
* `Session`
* `BasicAuth`
* `KeyAuth`

*Why?* Because `rays-web-security` handles CSRF protection, session policies (`SessionCreationPolicy`), and authentication internally at a global filter level. Applying them at the controller level would cause conflicts.

## Customizing Security `Csrf` and `Session`

Even though `rays-web-security` orchestrates CSRF and Sessions, you may still want to tweak their underlying Fiber configurations (e.g., changing the cookie name or timeout).

To do this, create a component that implements the `CsrfSecurityConfig` or `SessionSecurityConfig` interfaces:

```go
package configurations

import (
    "github.com/gofiber/fiber/v3/middleware/session"
    "github.com/gofiber/fiber/v3/middleware/csrf"
    . "github.com/BeamFoundry/rays/pkg/core"
)

type MySecurityTweaks struct {
    Component
}

// Customizes the Fiber Session middleware utilized by rays-web-security
func (this *MySecurityTweaks) SessionSecurityConfig(server string, config *session.Config) {
    if server == "primary" {
        config.CookieSecure = true // Force HTTPS cookies
    }
}

// Customizes the Fiber CSRF middleware utilized by rays-web-security
func (this *MySecurityTweaks) CsrfSecurityConfig(server string, config *csrf.Config) {
    if server == "primary" {
        config.CookieName = "my_custom_csrf"
    }
}

```

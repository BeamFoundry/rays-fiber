---
title: "Otel"
weight: 1
---

# Otel Middleware

*📖 **Official Go Fiber Documentation:** [https://docs.gofiber.io/contrib/otel/](https://docs.gofiber.io/contrib/otel/)*

OpenTelemetry instrumentation middleware adapted from `github.com/gofiber/contrib/v3/otel`.

## Installation

```bash
go get github.com/BeamFoundry/rays-fiber/pkg/middleware/contrib/otel
```

## Programmatic Configuration Interfaces

If your Controller or `FiberPathFilter` implements any of the following methods, the `Otel` middleware will automatically bind them:

* `OtelNext() func(fiber.Ctx) bool`
* `OtelTracerProvider() oteltrace.TracerProvider`
* `OtelMeterProvider() otelmetric.MeterProvider`
* `OtelPropagators() propagation.TextMapPropagator`
* `OtelResponsePropagators() propagation.TextMapPropagator`
* `OtelPublicEndpointFn() func(fiber.Ctx) bool`
* `OtelSpanStartOptions() []oteltrace.SpanStartOption`
* `OtelSpanNameFormatter() func(fiber.Ctx) string`
* `OtelCustomAttributes() func(fiber.Ctx) []attribute.KeyValue`
* `OtelCustomResponseAttributes() func(fiber.Ctx) []attribute.KeyValue`
* `OtelCustomMetricAttributes() func(fiber.Ctx) []attribute.KeyValue`
* `OtelCustomResponseMetricAttributes() func(fiber.Ctx) []attribute.KeyValue`
* `OtelOptions(options []otel.Option) []otel.Option`

### Example: Configuring `Otel`

Here is an example of a controller embedding the middleware and implementing the `OtelOptions` interface to customize the underlying `otel.Option` slice:

```go
package controllers

import (
    "github.com/gofiber/contrib/v3/otel"
    . "github.com/BeamFoundry/rays-fiber/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/core"
    . "github.com/BeamFoundry/rays/pkg/web"
    . "github.com/BeamFoundry/rays-fiber/pkg/middleware/contrib/otel"
)

type MyOtelController struct {
    FiberController
    RequestMapping `Path:"/example"`

    // 1. Embed the middleware
    Otel
}

// 2. Implement the configuration interface
func (this *MyOtelController) OtelOptions(options []otel.Option) []otel.Option {
    // Append or override otel.Option values for this controller
    return options
}

var _, _ = any(&MyOtelController{}).(Stereotype)
```

## Default YAML Configuration

The framework applies these default settings to the `"primary"` server. To override them, copy the properties you want to change into your `application.yaml` under `rays.fiber.servers.primary.otel`.

```yaml
rays:
  fiber:
    servers:
      primary:
        otel:
          withoutMetrics: false
          traceResponseHeader: ""
          publicEndpoint: false
          port: null
          clientIP: true
          capturedRequestHeaders: []
          capturedResponseHeaders: []
          redactedQueryParams: []
```

## Integration with `rays-otel` & `SecurityFilterChain`

- **Container Provider Injection**: When `github.com/BeamFoundry/rays-otel` (or custom `TracerProvider`, `MeterProvider`, or `TextMapPropagator` Beams) is registered in the container, `Otel` automatically injects them via `Optional[T]` before applying controller/filter interface overrides (`OtelTracerProvider`, `OtelMeterProvider`, `OtelPropagators`, `OtelOptions`).
- **`TransactionManager` Context Binding**: `rays-fiber` automatically binds `c.Context()` onto `TransactionManager.SetContext(...)` before `SecurityFilterChain` and before controller handler execution, so downstream `rays-orm`, `rays-data-mongo`, `rays-web-security-oauth2`, and `Logging` calls inherit the active HTTP span.
- **Middleware Ordering**: `SecurityFilterChain` runs after `FiberPathFilter`s and before `FiberController` middlewares. To include authentication/authorization inside the HTTP server span, embed `Otel` in a global `FiberPathFilter`:

```go
type GlobalOtelFilter struct {
    FiberPathFilter
    Otel
}

var _, _ = any(&GlobalOtelFilter{}).(Stereotype)
```


package otel

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// OtelNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to indicate when to skip this middleware.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelNext interface {
	OtelNext() func(fiber.Ctx) bool
}

// OtelTracerProvider is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to specify a tracer provider to use for creating a tracer.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelTracerProvider interface {
	OtelTracerProvider() oteltrace.TracerProvider
}

// OtelMeterProvider is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to specify a meter provider to use for reporting.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelMeterProvider interface {
	OtelMeterProvider() otelmetric.MeterProvider
}

// OtelPropagators is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to specify propagators to use for extracting information from HTTP requests.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelPropagators interface {
	OtelPropagators() propagation.TextMapPropagator
}

// OtelResponsePropagators is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to set the propagators that inject the trace context into responses.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelResponsePropagators interface {
	OtelResponsePropagators() propagation.TextMapPropagator
}

// OtelPublicEndpointFn is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to decide per request whether it comes from a public endpoint.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelPublicEndpointFn interface {
	OtelPublicEndpointFn() func(fiber.Ctx) bool
}

// OtelSpanStartOptions is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to add options to every server span when it starts.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelSpanStartOptions interface {
	OtelSpanStartOptions() []oteltrace.SpanStartOption
}

// OtelSpanNameFormatter is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to customize how recording spans are named once the route is known.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelSpanNameFormatter interface {
	OtelSpanNameFormatter() func(fiber.Ctx) string
}

// OtelCustomAttributes is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to add custom attributes to the span when it starts.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelCustomAttributes interface {
	OtelCustomAttributes() func(fiber.Ctx) []attribute.KeyValue
}

// OtelCustomResponseAttributes is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to add span attributes after the handler chain has run.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelCustomResponseAttributes interface {
	OtelCustomResponseAttributes() func(fiber.Ctx) []attribute.KeyValue
}

// OtelCustomMetricAttributes is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to add custom attributes to all metrics of a request.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelCustomMetricAttributes interface {
	OtelCustomMetricAttributes() func(fiber.Ctx) []attribute.KeyValue
}

// OtelCustomResponseMetricAttributes is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to add attributes to the duration and body size metrics after the handler chain has run.
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelCustomResponseMetricAttributes interface {
	OtelCustomResponseMetricAttributes() func(fiber.Ctx) []attribute.KeyValue
}

// OtelOptions is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement,
// if it uses [Otel], to review and customize the [otel.Option] slice that will be passed to [otel.New].
//
// Please see the [otel] middleware configuration documentation.
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type OtelOptions interface {
	OtelOptions(options []otel.Option) []otel.Option
}

type otelConfigurationProperties struct {
	port                    *int
	traceResponseHeader     string
	capturedRequestHeaders  []string
	capturedResponseHeaders []string
	redactedQueryParams     []string
	withoutMetrics          bool
	publicEndpoint          bool
	clientIP                bool
}

type otelConfig struct {
	otel otelConfigurationProperties
}

var _, _ = any(&FiberMwConfig[otelConfig]{}).(Stereotype)

func (this *otelConfig) Options() []otel.Option {
	opts := []otel.Option{
		otel.WithoutMetrics(this.otel.withoutMetrics),
		otel.WithClientIP(this.otel.clientIP),
	}
	if this.otel.traceResponseHeader != "" {
		opts = append(opts, otel.WithTraceResponseHeader(this.otel.traceResponseHeader))
	}
	if this.otel.publicEndpoint {
		opts = append(opts, otel.WithPublicEndpoint())
	}
	if this.otel.port != nil {
		opts = append(opts, otel.WithPort(*this.otel.port))
	}
	if len(this.otel.capturedRequestHeaders) > 0 {
		opts = append(opts, otel.WithCapturedRequestHeaders(this.otel.capturedRequestHeaders...))
	}
	if len(this.otel.capturedResponseHeaders) > 0 {
		opts = append(opts, otel.WithCapturedResponseHeaders(this.otel.capturedResponseHeaders...))
	}
	if len(this.otel.redactedQueryParams) > 0 {
		opts = append(opts, otel.WithRedactedQueryParams(this.otel.redactedQueryParams...))
	}
	return opts
}

// [Otel] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Otel
//	}
//
//	func (this *ApiController) OtelOptions(options []otel.Option) []otel.Option {
//		// changing configuration for this controller.
//		return options
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPathFilter] may implement.
// Those interfaces are [OtelNext], [OtelTracerProvider], [OtelMeterProvider], [OtelPropagators],
// [OtelResponsePropagators], [OtelPublicEndpointFn], [OtelSpanStartOptions], [OtelSpanNameFormatter],
// [OtelCustomAttributes], [OtelCustomResponseAttributes], [OtelCustomMetricAttributes], and
// [OtelCustomResponseMetricAttributes].
//
// Also, the parent container may implement [OtelOptions] to receive and modify the [otel.Option] slice
// that will be used.
//
// Please see those interface definitions and the [otel] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					otel:
//						withoutMetrics: false
//						traceResponseHeader: ""
//						publicEndpoint: false
//						port: null
//						clientIP: true
//						capturedRequestHeaders: []
//						capturedResponseHeaders: []
//						redactedQueryParams: []
//
// [otel]: https://docs.gofiber.io/contrib/otel/
type Otel struct {
	config         *FiberMwConfig[otelConfig]              `@:"Inject"`
	tracerProvider Optional[oteltrace.TracerProvider]      `@:"Inject"`
	meterProvider  Optional[otelmetric.MeterProvider]      `@:"Inject"`
	propagator     Optional[propagation.TextMapPropagator] `@:"Inject"`
	opts           []otel.Option
}

// PostConstruct implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct].
func (this *Otel) PostConstruct(parent any) Error {
	this.opts = this.config.GetConfig(parent).Options()

	if tp, ok := this.tracerProvider.Get(); ok {
		this.opts = append(this.opts, otel.WithTracerProvider(tp))
	}
	if mp, ok := this.meterProvider.Get(); ok {
		this.opts = append(this.opts, otel.WithMeterProvider(mp))
	}
	if prop, ok := this.propagator.Get(); ok {
		this.opts = append(this.opts, otel.WithPropagators(prop))
	}

	if next, ok := parent.(OtelNext); ok {
		if fn := next.OtelNext(); fn != nil {
			this.opts = append(this.opts, otel.WithNext(fn))
		}
	}
	if tp, ok := parent.(OtelTracerProvider); ok {
		if provider := tp.OtelTracerProvider(); provider != nil {
			this.opts = append(this.opts, otel.WithTracerProvider(provider))
		}
	}
	if mp, ok := parent.(OtelMeterProvider); ok {
		if provider := mp.OtelMeterProvider(); provider != nil {
			this.opts = append(this.opts, otel.WithMeterProvider(provider))
		}
	}
	if prop, ok := parent.(OtelPropagators); ok {
		if propagators := prop.OtelPropagators(); propagators != nil {
			this.opts = append(this.opts, otel.WithPropagators(propagators))
		}
	}
	if respProp, ok := parent.(OtelResponsePropagators); ok {
		if propagators := respProp.OtelResponsePropagators(); propagators != nil {
			this.opts = append(this.opts, otel.WithResponsePropagators(propagators))
		}
	}
	if pubFn, ok := parent.(OtelPublicEndpointFn); ok {
		if fn := pubFn.OtelPublicEndpointFn(); fn != nil {
			this.opts = append(this.opts, otel.WithPublicEndpointFn(fn))
		}
	}
	if startOpts, ok := parent.(OtelSpanStartOptions); ok {
		if opts := startOpts.OtelSpanStartOptions(); len(opts) > 0 {
			this.opts = append(this.opts, otel.WithSpanStartOptions(opts...))
		}
	}
	if formatter, ok := parent.(OtelSpanNameFormatter); ok {
		if fn := formatter.OtelSpanNameFormatter(); fn != nil {
			this.opts = append(this.opts, otel.WithSpanNameFormatter(fn))
		}
	}
	if attrs, ok := parent.(OtelCustomAttributes); ok {
		if fn := attrs.OtelCustomAttributes(); fn != nil {
			this.opts = append(this.opts, otel.WithCustomAttributes(fn))
		}
	}
	if respAttrs, ok := parent.(OtelCustomResponseAttributes); ok {
		if fn := respAttrs.OtelCustomResponseAttributes(); fn != nil {
			this.opts = append(this.opts, otel.WithCustomResponseAttributes(fn))
		}
	}
	if metricAttrs, ok := parent.(OtelCustomMetricAttributes); ok {
		if fn := metricAttrs.OtelCustomMetricAttributes(); fn != nil {
			this.opts = append(this.opts, otel.WithCustomMetricAttributes(fn))
		}
	}
	if respMetricAttrs, ok := parent.(OtelCustomResponseMetricAttributes); ok {
		if fn := respMetricAttrs.OtelCustomResponseMetricAttributes(); fn != nil {
			this.opts = append(this.opts, otel.WithCustomResponseMetricAttributes(fn))
		}
	}
	if options, ok := parent.(OtelOptions); ok {
		this.opts = options.OtelOptions(this.opts)
	}
	return nil
}

// Handler implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler].
func (this *Otel) Handler() fiber.Handler {
	return otel.New(this.opts...)
}

// PreTxBinding implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberPreTxHandler]
// so the OpenTelemetry middleware runs before transaction context binding.
func (this *Otel) PreTxBinding() {}

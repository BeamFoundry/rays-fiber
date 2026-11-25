package otel_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/contrib/otel"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

type OtelPathFilter struct {
	FiberPathFilter
	Otel
}

var (
	_ OtelNext                           = &OtelPathFilter{}
	_ OtelTracerProvider                 = &OtelPathFilter{}
	_ OtelMeterProvider                  = &OtelPathFilter{}
	_ OtelPropagators                    = &OtelPathFilter{}
	_ OtelResponsePropagators            = &OtelPathFilter{}
	_ OtelPublicEndpointFn               = &OtelPathFilter{}
	_ OtelSpanStartOptions               = &OtelPathFilter{}
	_ OtelSpanNameFormatter              = &OtelPathFilter{}
	_ OtelCustomAttributes               = &OtelPathFilter{}
	_ OtelCustomResponseAttributes       = &OtelPathFilter{}
	_ OtelCustomMetricAttributes         = &OtelPathFilter{}
	_ OtelCustomResponseMetricAttributes = &OtelPathFilter{}
	_ OtelOptions                        = &OtelPathFilter{}
)

func (o *OtelPathFilter) OtelNext() func(fiber.Ctx) bool {
	return func(fiber.Ctx) bool { return false }
}

func (o *OtelPathFilter) OtelTracerProvider() oteltrace.TracerProvider {
	return &wrappedTP{TracerProvider: nooptrace.NewTracerProvider()}
}

func (o *OtelPathFilter) OtelMeterProvider() otelmetric.MeterProvider {
	return &wrappedMP{MeterProvider: noopmetric.NewMeterProvider()}
}

func (o *OtelPathFilter) OtelPropagators() propagation.TextMapPropagator {
	return &wrappedProp{TextMapPropagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})}
}

func (o *OtelPathFilter) OtelResponsePropagators() propagation.TextMapPropagator {
	return &wrappedProp{TextMapPropagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})}
}

func (o *OtelPathFilter) OtelPublicEndpointFn() func(fiber.Ctx) bool {
	return func(fiber.Ctx) bool { return false }
}

func (o *OtelPathFilter) OtelSpanStartOptions() []oteltrace.SpanStartOption {
	return []oteltrace.SpanStartOption{oteltrace.WithSpanKind(oteltrace.SpanKindServer)}
}

func (o *OtelPathFilter) OtelSpanNameFormatter() func(fiber.Ctx) string {
	return func(c fiber.Ctx) string { return c.Path() }
}

func (o *OtelPathFilter) OtelCustomAttributes() func(fiber.Ctx) []attribute.KeyValue {
	return func(fiber.Ctx) []attribute.KeyValue { return []attribute.KeyValue{attribute.String("k", "v")} }
}

func (o *OtelPathFilter) OtelCustomResponseAttributes() func(fiber.Ctx) []attribute.KeyValue {
	return func(fiber.Ctx) []attribute.KeyValue { return []attribute.KeyValue{attribute.String("rk", "rv")} }
}

func (o *OtelPathFilter) OtelCustomMetricAttributes() func(fiber.Ctx) []attribute.KeyValue {
	return func(fiber.Ctx) []attribute.KeyValue { return []attribute.KeyValue{attribute.String("mk", "mv")} }
}

func (o *OtelPathFilter) OtelCustomResponseMetricAttributes() func(fiber.Ctx) []attribute.KeyValue {
	return func(fiber.Ctx) []attribute.KeyValue { return []attribute.KeyValue{attribute.String("rmk", "rmv")} }
}

func (o *OtelPathFilter) OtelOptions(options []otel.Option) []otel.Option {
	return options
}

type (
	wrappedTP   struct{ oteltrace.TracerProvider }
	wrappedMP   struct{ otelmetric.MeterProvider }
	wrappedProp struct{ propagation.TextMapPropagator }
)

type testProviderConfig struct {
	Configuration
	IgnoreScanning
}

func (this *testProviderConfig) TracerProvider() oteltrace.TracerProvider {
	return &wrappedTP{TracerProvider: nooptrace.NewTracerProvider()}
}

func (this *testProviderConfig) MeterProvider() otelmetric.MeterProvider {
	return &wrappedMP{MeterProvider: noopmetric.NewMeterProvider()}
}

func (this *testProviderConfig) Propagator() propagation.TextMapPropagator {
	return &wrappedProp{TextMapPropagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})}
}

var _ = Describe("Otel", func() {
	It("can use Otel middleware with Optional providers and config options", func() {
		port := 8443
		Expect(GetContainer().LoadConfigAt(map[string]any{
			"withoutMetrics":          false,
			"traceResponseHeader":     "X-Trace-Id",
			"publicEndpoint":          true,
			"port":                    port,
			"clientIP":                true,
			"capturedRequestHeaders":  []any{"X-Req"},
			"capturedResponseHeaders": []any{"X-Resp"},
			"redactedQueryParams":     []any{"token"},
		}, "rays.fiber.servers.primary.otel")).NotTo(HaveOccurred())

		AddBeamFor(&testProviderConfig{})
		AddBeamFor(&OtelPathFilter{})
		AddBeamFor(&OtelTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		appProvider := GetBeamScopedInstance[FiberAppProvider]()
		app := appProvider.Default()
		Expect(app).NotTo(BeNil())

		req := httptest.NewRequest(http.MethodGet, "/testing/hello", http.NoBody)
		resp, testErr := app.Test(req)
		Expect(testErr).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

		(&Otel{}).PreTxBinding()
	})
})

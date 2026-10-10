package controllers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v3"

	_ "github.com/BeamFoundry/rays"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/internal/controllers"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	_ "github.com/BeamFoundry/rays-web-security"
	. "github.com/BeamFoundry/rays-web-security/pkg/config"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Internal/Controllers Suite")
}

type TestCustomizer struct {
	Qualifier `Name:"primary"`
}

func (this *TestCustomizer) FiberConfig(*fiber.Config) {}

func (this *TestCustomizer) FiberListenConfig(*fiber.ListenConfig) {}

var _ FiberConfig = &TestCustomizer{}

var _ FiberListenConfig = &TestCustomizer{}

type TestAppMiddleware struct {
	beam *BeamWithValue `@:"Inject,Qualifier{app}"`
}

func (this *TestAppMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Test", this.beam.GetValue())
		return c.Next()
	}
}

type TestAppPathMiddleware struct{}

func (this *TestAppPathMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Path", c.Path())
		return c.Next()
	}
}

type TestRestrictedPathMiddleware struct{}

func (this *TestRestrictedPathMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Restricted", c.Path())
		return c.Next()
	}
}

type TestPathFilter struct {
	FiberPathFilter
	TestAppMiddleware
	// FiberAppSimpleCustomizer
}

//nolint:govet // Testing
type TestPathFilterPath struct {
	FiberPathFilter
	RequestMapping `Path:"/testing"`
	TestAppPathMiddleware
}

//nolint:govet // Testing
type TestPathFilterMethod struct {
	FiberPathFilter
	RequestMapping `Path:"/testing/*" Method:"GET"`
	TestAppPathMiddleware
}

//nolint:govet // Testing
type TestPathFilterMethodBAD struct {
	FiberPathFilter
	RequestMapping `Path:"/testing" Method:"BAD"`
	TestAppPathMiddleware
}

//nolint:govet // Testing
type TestPathFilterMethodALL struct {
	FiberPathFilter
	RequestMapping `Path:"/testing" Method:"ALL"`
	TestAppPathMiddleware
}

type TestGroupMiddleware struct {
	beam *BeamWithValue `@:"Inject,Qualifier{group}"`
}

func (this *TestGroupMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Append("Test", this.beam.GetValue())
		return c.Next()
	}
}

type TestPathMiddleware struct {
	beam *BeamWithValue `@:"Inject,Qualifier{path}"`
}

func (this *TestPathMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Append("Test", this.beam.GetValue())
		return c.Next()
	}
}

type TestFilterGroup struct {
	beam *BeamWithValue `@:"Inject,Qualifier{group}"`
}

func (this *TestFilterGroup) Filter(handler fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Filter", this.beam.GetValue())
		return handler(c)
	}
}

type TestFilterPath struct {
	beam *BeamWithValue `@:"Inject,Qualifier{path}"`
}

func (this *TestFilterPath) Filter(handler fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Filter", this.beam.GetValue())
		return handler(c)
	}
}

type TestControllerConfiguration struct {
	Configuration
}

func (this *TestControllerConfiguration) Controller(_ struct{ Primary }) *BeamWithValue {
	return (&BeamWithValue{}).WithValue("hello")
}

func (this *TestControllerConfiguration) App() *BeamWithValue {
	return (&BeamWithValue{}).WithValue("App")
}

func (this *TestControllerConfiguration) Group() *BeamWithValue {
	return (&BeamWithValue{}).WithValue("Group")
}

func (this *TestControllerConfiguration) Path() *BeamWithValue {
	return (&BeamWithValue{}).WithValue("Path")
}

type TestController struct {
	Controller
	Server         `Name:"primary"`
	RequestMapping `Path:"/testing"`
	TestGroupMiddleware
	TestFilterGroup
	TestRestrictedPathMiddleware
	beam *BeamWithValue `@:"Inject"`
}

func (this *TestController) Testing(ctx fiber.Ctx) error {
	return ctx.SendString(this.beam.GetValue())
}

//nolint:govet // Singleton
func (this *TestController) Hello1(
	_ struct {
		RequestMapping `Path:"/hello1" Method:"GET"`
		//
		TestRestrictedPathMiddleware
	}, ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) Hello2(
	_ struct {
		RequestMapping `Path:"/hello2"`
		TestPathMiddleware
		TestRestrictedPathMiddleware
		TestFilterPath
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) Hello3(
	_ *struct {
		RequestMapping `Path:"/hello3"`
		TestRestrictedPathMiddleware
		TestPathMiddleware
		TestFilterPath
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloError(
	_ struct {
		RequestMapping `Path:"/helloError" Method:"GET"`
	},
	_ fiber.Ctx,
) Error {
	return nil
}

func (this *TestController) HelloHead(
	_ struct {
		RequestMapping `Path:"/helloHEAD" Method:"HEAD"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloPost(
	_ struct {
		RequestMapping `Path:"/helloPOST" Method:"POST"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloPut(
	_ struct {
		RequestMapping `Path:"/helloPUT" Method:"PUT"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloDelete(
	_ struct {
		RequestMapping `Path:"/helloDELETE" Method:"DELETE"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloConnect(
	_ struct {
		RequestMapping `Path:"/helloCONNECT" Method:"CONNECT"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloOptions(
	_ struct {
		RequestMapping `Path:"/helloOPTIONS" Method:"OPTIONS"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloTrace(
	_ struct {
		RequestMapping `Path:"/helloTRACE" Method:"TRACE"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloPatch(
	_ struct {
		RequestMapping `Path:"/helloPATCH" Method:"PATCH"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) HelloAll(
	_ struct {
		RequestMapping `Path:"/helloALL" Method:"ALL"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) MethodCanonical(
	_ struct {
		RequestMapping `Path:"/canonical"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestController) MethodCanonicalVoid(
	_ struct {
		RequestMapping `Path:"/canonicalVoid"`
	},
	ctx fiber.Ctx,
) {
	//nolint:errcheck // Testing
	ctx.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpress(
	_ struct {
		RequestMapping `Path:"/express"`
	},
	_ fiber.Req, res fiber.Res,
) error {
	return res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressVoid(
	_ struct {
		RequestMapping `Path:"/expressVoid"`
	},
	_ fiber.Req, res fiber.Res,
) {
	//nolint:errcheck // Testing
	res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNext(
	_ struct {
		RequestMapping `Path:"/expressNext"`
	},
	_ fiber.Req, res fiber.Res, _ func() error,
) error {
	return res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextVoid(
	_ struct {
		RequestMapping `Path:"/expressNextVoid"`
	},
	_ fiber.Req, res fiber.Res, _ func() error,
) {
	//nolint:errcheck // Testing
	res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextNo(
	_ struct {
		RequestMapping `Path:"/expressNextNo"`
	},
	_ fiber.Req, res fiber.Res, _ func(),
) error {
	return res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextNoVoid(
	_ struct {
		RequestMapping `Path:"/expressNextNoVoid"`
	},
	_ fiber.Req, res fiber.Res, _ func(),
) {
	//nolint:errcheck // Testing
	res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextErr(
	_ struct {
		RequestMapping `Path:"/expressNextErr"`
	},
	_ fiber.Req, res fiber.Res, _ func(error) error,
) error {
	return res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextErrVoid(
	_ struct {
		RequestMapping `Path:"/expressNextErrVoid"`
	},
	_ fiber.Req, res fiber.Res, _ func(error) error,
) {
	//nolint:errcheck // Testing
	res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextErrNo(
	_ struct {
		RequestMapping `Path:"/expressNextErrNo"`
	},
	_ fiber.Req, res fiber.Res, _ func(error),
) error {
	return res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodExpressNextErrNoVoid(
	_ struct {
		RequestMapping `Path:"/expressNextErrNoVoid"`
	},
	_ fiber.Req, res fiber.Res, _ func(error),
) {
	//nolint:errcheck // Testing
	res.SendString(this.beam.GetValue())
}

func (this *TestController) MethodHttpHandler(
	_ struct {
		RequestMapping `Path:"/httpHandler"`
	},
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%s", this.beam.GetValue())
	})
}

func (this *TestController) MethodHttpHandlerFunc(
	_ struct {
		RequestMapping `Path:"/httpHandlerFunc"`
	},
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "%s", this.beam.GetValue())
	})
}

func (this *TestController) MethodHttpFunc(
	_ struct {
		RequestMapping `Path:"/httpFunc"`
	},
	w http.ResponseWriter, _ *http.Request,
) {
	fmt.Fprintf(w, "%s", this.beam.GetValue())
}

func (this *TestController) MethodFastHttp(
	_ struct {
		RequestMapping `Path:"/fastHttp"`
	},
	ctx *fasthttp.RequestCtx,
) {
	ctx.SetBodyString(this.beam.GetValue())
}

func (this *TestController) MethodFastHttpErr(
	_ struct {
		RequestMapping `Path:"/fastHttpErr"`
	},
	ctx *fasthttp.RequestCtx,
) error {
	ctx.SetBodyString(this.beam.GetValue())
	return nil
}

type TestNoMappingController struct {
	Controller
	TestGroupMiddleware
	TestFilterGroup
	beam *BeamWithValue `@:"Inject"`
}

func (this *TestNoMappingController) Hello1(
	_ struct {
		RequestMapping `Path:"/nomapping/hello1" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestNoMappingController) Invalid(ctx fiber.Ctx) error {
	return ctx.SendString(this.beam.GetValue())
}

type TestNoMappingPathController struct {
	Controller
	RequestMapping
	TestGroupMiddleware
	TestFilterGroup
	beam *BeamWithValue `@:"Inject"`
}

func (this *TestNoMappingPathController) Hello1(
	_ struct {
		RequestMapping `Path:"/nomappingpath/hello1" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestNoMappingPathController) Invalid(
	_ struct {
		RequestMapping
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

func (this *TestNoMappingPathController) InvalidMethod(
	_ struct {
		RequestMapping `Path:"/nomappingpath/foo" Method:"FOO"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString(this.beam.GetValue())
}

type TestNoMethodController struct {
	Controller
	RequestMapping `Path:"/nomethod"`
	TestGroupMiddleware
}

type TestInvalidController struct {
	Controller
	Server         `_:"invalid"`
	RequestMapping `Path:"/invalid"`
}

func (this *TestInvalidController) Hello(
	_ struct {
		RequestMapping `Path:"/hello" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString("hello")
}

var _ = Describe("Controller", func() {
	hello := func(path, header, filter string) {
		appProvider := GetBeamScopedInstance[FiberAppProvider]()
		app := appProvider.Default()
		Expect(app).NotTo(BeNil())

		GetLogger().Debug("Testing", "path", path, "header", header, "filter", filter)
		req := httptest.NewRequest(http.MethodGet, path, http.NoBody)
		resp, err := app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(resp.Header.Get("Test")).To(Equal(header))
		Expect(resp.Header.Get("Filter")).To(Equal(filter))
		if String(path).StartsWith("/testing") {
			Expect(resp.Header.Get("Path")).To(Equal(path))
		}
		if String(path).StartsWith("/testing/hello") {
			Expect(resp.Header.Get("Restricted")).NotTo(BeEmpty())
		} else {
			Expect(resp.Header.Get("Restricted")).To(BeEmpty())
		}
		body, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(body)).To(Equal("hello"))
	}
	It("can declare a controller", func() {
		AddBeamFor(&TestControllerConfiguration{})
		AddBeamFor(&TestNoMappingController{})
		AddBeamFor(&TestNoMappingPathController{})
		AddBeamFor(&TestNoMethodController{})
		AddBeamFor(&TestController{})
		AddBeamFor(&TestPathFilter{})
		AddBeamFor(&TestPathFilterPath{})
		AddBeamFor(&TestPathFilterMethod{})
		AddBeamFor(&TestPathFilterMethodBAD{})
		AddBeamFor(&TestPathFilterMethodALL{})
		httpSecurity := (&HttpSecurityBuilder{}).
			AnyRequest().PermitAll().Build()
		AddBeamFor(httpSecurity)
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		hello("/nomapping/hello1", "App, Group", "Group")
		hello("/nomappingpath/hello1", "App, Group", "Group")
		hello("/testing/hello1", "App, Group", "Group")
		hello("/testing/hello2", "App, Group, Path", "Path")
		hello("/testing/hello3", "App, Group, Path", "Path")
		hello("/testing/canonical", "App, Group", "Group")
		hello("/testing/canonicalVoid", "App, Group", "")
		hello("/testing/express", "App, Group", "")
		hello("/testing/expressVoid", "App, Group", "")
		hello("/testing/expressNext", "App, Group", "")
		hello("/testing/expressNextVoid", "App, Group", "")
		hello("/testing/expressNextNo", "App, Group", "")
		hello("/testing/expressNextNoVoid", "App, Group", "")
		hello("/testing/expressNextErr", "App, Group", "")
		hello("/testing/expressNextErrVoid", "App, Group", "")
		hello("/testing/expressNextErrNo", "App, Group", "")
		hello("/testing/expressNextErrNoVoid", "App, Group", "")
		hello("/testing/httpHandler", "App, Group", "")
		hello("/testing/httpFunc", "App, Group", "")
		hello("/testing/fastHttp", "App, Group", "")
		hello("/testing/fastHttpErr", "App, Group", "")
	})

	It("can use a customizer", func() {
		AddBeamFor(&TestCustomizer{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
		server := GetBeamScopedInstance[*FiberAppServerConfig]()
		Expect(server.Config("primary")).NotTo(BeNil())
		Expect(server.Listener("primary")).NotTo(BeNil())

		provider := GetBeamScopedInstance[FiberAppProvider]()
		Expect(provider.Provide("primary")).NotTo(BeNil())
		Expect(provider.Provide("default")).To(BeNil())
	})

	It("can and refuse to register a controller on an invalid server", func() {
		AddBeamFor(&TestInvalidController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).To(HaveOccurred())
	})

	It("can be started and halted", func() {
		cfg := map[string]any{
			"address": "/var/tmp/fiber.sock",
			"listener": map[string]any{
				"disableStartupMessage": true,
				"listenerNetwork":       "unix",
			},
		}
		//nolint:errcheck //Not tested
		GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary")

		AddBeamFor(&TestControllerConfiguration{})
		AddBeamFor(&TestNoMethodController{})
		AddBeamFor(&TestController{})
		AddBeamFor(&TestPathFilter{})
		AddBeamFor(&TestPathFilterPath{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[FiberAppProvider]()
		go func() {
			GetContainer().Run()
		}()
		time.Sleep(time.Second)
		provider.Shutdown()
	})
})

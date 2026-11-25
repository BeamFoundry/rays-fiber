package controllers_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	_ "github.com/BeamFoundry/rays-web-security"
	. "github.com/BeamFoundry/rays-web-security/pkg/config"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/testing"
)

type preTxCtxKey struct{}

type preTxOrderRecorder struct {
	txMgr  TransactionManager
	events []string
}

type TestPreTxMiddleware struct {
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestPreTxMiddleware) PreTxBinding() {}

func (this *TestPreTxMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		hasPreTx := this.rec.txMgr.GetContext().Value(preTxCtxKey{}) != nil
		Expect(hasPreTx).To(BeFalse())
		this.rec.events = append(this.rec.events, "preTx")
		c.SetContext(context.WithValue(c.Context(), preTxCtxKey{}, "enriched"))
		return c.Next()
	}
}

type TestMethodPreTxMiddleware struct {
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestMethodPreTxMiddleware) PreTxBinding() {}

func (this *TestMethodPreTxMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		this.rec.events = append(this.rec.events, "methodPreTx")
		return c.Next()
	}
}

type TestServiceCallingMiddleware struct {
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestServiceCallingMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
		Expect(val).To(Equal("enriched"))
		this.rec.events = append(this.rec.events, "ctrlMiddleware")
		return c.Next()
	}
}

type TestMethodServiceMiddleware struct {
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestMethodServiceMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
		Expect(val).To(Equal("enriched"))
		this.rec.events = append(this.rec.events, "methodMiddleware")
		return c.Next()
	}
}

type TestTxOrderController struct {
	Controller
	RequestMapping `Path:"/tx-order"`
	TestPreTxMiddleware
	TestServiceCallingMiddleware
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestTxOrderController) WithMethodMiddlewares(
	_ struct {
		RequestMapping `Path:"/method" Method:"GET"`
		TestMethodPreTxMiddleware
		TestMethodServiceMiddleware
	},
	c fiber.Ctx,
) error {
	val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
	Expect(val).To(Equal("enriched"))
	this.rec.events = append(this.rec.events, "handler")
	return c.SendString("ok")
}

func (this *TestTxOrderController) Plain(
	c fiber.Ctx,
) error {
	val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
	Expect(val).To(Equal("enriched"))
	this.rec.events = append(this.rec.events, "plain")
	return c.SendString("plain")
}

type TestPreTxPathFilter struct {
	FiberPathFilter
	TestPreTxMiddleware
}

type TestPlainController struct {
	Controller
	TestServiceCallingMiddleware
	rec *preTxOrderRecorder `@:"Inject"`
}

func (this *TestPlainController) Ping(
	_ struct {
		RequestMapping `Path:"/ping-filter" Method:"GET"`
	},
	c fiber.Ctx,
) error {
	val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
	Expect(val).To(Equal("enriched"))
	this.rec.events = append(this.rec.events, "handler")
	return c.SendString("pong")
}

type TestPreTxUserDetailsService struct {
	Service
	encoder PasswordEncoder     `@:"Inject"`
	rec     *preTxOrderRecorder `@:"Inject"`
}

var _ UserDetailsService = &TestPreTxUserDetailsService{}

func (this *TestPreTxUserDetailsService) LoadUserByUsername(username string) UserDetails {
	val, _ := this.rec.txMgr.GetContext().Value(preTxCtxKey{}).(string)
	Expect(val).To(Equal("enriched"))
	this.rec.events = append(this.rec.events, "securityFilter")
	mgr := InMemoryUserDetailsManager{}.New(this.encoder)
	mgr.NewUser().WithUsername(username).WithPassword("pass").Build()
	return mgr.LoadUserByUsername(username)
}

var _ = Describe("Transaction context middleware ordering", func() {
	It("runs FiberPreTxHandler before bindTxContextHandler and other controller/method middlewares after without httpSecurity", func() {
		rec := &preTxOrderRecorder{}
		AddBeamFor(rec)
		AddBeamFor(&TestTxOrderController{})
		GetContainer().EnableScanning()
		Expect(GetContainer().Initialize()).To(Succeed())

		rec.txMgr = GetBeamScopedInstance[TransactionManager]()
		app := GetBeamScopedInstance[FiberAppProvider]().Default()

		req := httptest.NewRequest(http.MethodGet, "/tx-order/method", http.NoBody)
		resp, err := app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(rec.events).To(Equal([]string{"preTx", "methodPreTx", "ctrlMiddleware", "methodMiddleware", "handler"}))
		Expect(rec.txMgr.GetContext().Value(preTxCtxKey{})).To(BeNil())

		rec.events = nil
		req = httptest.NewRequest(http.MethodGet, "/tx-order", http.NoBody)
		resp, err = app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(rec.events).To(Equal([]string{"preTx", "ctrlMiddleware", "plain"}))
		Expect(rec.txMgr.GetContext().Value(preTxCtxKey{})).To(BeNil())
	})

	It("runs FiberPathFilter before bindTxContextHandler and HttpSecurity filter chain", func() {
		rec := &preTxOrderRecorder{}
		AddBeamFor(rec)
		AddBeamFor(&TestPreTxPathFilter{})
		AddBeamFor(&TestPreTxUserDetailsService{})
		AddBeamFor(&TestPlainController{})
		httpSecurity := (&HttpSecurityBuilder{}).
			AnyRequest().Authenticated().
			Basic().
			Build()
		AddBeamFor(httpSecurity)
		GetContainer().EnableScanning()
		Expect(GetContainer().Initialize()).To(Succeed())

		rec.txMgr = GetBeamScopedInstance[TransactionManager]()
		app := GetBeamScopedInstance[FiberAppProvider]().Default()

		req := httptest.NewRequest(http.MethodGet, "/ping-filter", http.NoBody)
		req.SetBasicAuth("user", "pass")
		resp, err := app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(rec.events).To(Equal([]string{"preTx", "securityFilter", "ctrlMiddleware", "handler"}))
		Expect(rec.txMgr.GetContext().Value(preTxCtxKey{})).To(BeNil())
	})
})

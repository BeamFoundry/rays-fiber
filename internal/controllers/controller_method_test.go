package controllers_test

import (
	"fmt"
	"net/http"
	"reflect"

	f "github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/util"

	. "github.com/BeamFoundry/rays-fiber/internal/controllers"
)

var httpHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello World!")
})

type FunErr = func(error)

type FunErrErr = func(error) error

type TestMethod struct {
	_ string
}

func (this *TestMethod) Canonical(_ struct{}, _ f.Ctx) error                         { return nil }
func (this *TestMethod) CanonicalStrict(_ f.Ctx) error                               { return nil }
func (this *TestMethod) CanonicalVoid(_ struct{}, _ f.Ctx)                           {}
func (this *TestMethod) CanonicalVoidStrict(_ f.Ctx)                                 {}
func (this *TestMethod) Expr(_ struct{}, _ f.Req, _ f.Res) error                     { return nil }
func (this *TestMethod) ExprStrict(_ f.Req, _ f.Res) error                           { return nil }
func (this *TestMethod) ExprVoid(_ struct{}, _ f.Req, _ f.Res)                       {}
func (this *TestMethod) ExprVoidStrict(_ f.Req, _ f.Res)                             {}
func (this *TestMethod) ExprNt(_ struct{}, _ f.Req, _ f.Res, _ func() error) error   { return nil }
func (this *TestMethod) ExprNtStrict(_ f.Req, _ f.Res, _ func() error) error         { return nil }
func (this *TestMethod) ExprNtVoid(_ struct{}, _ f.Req, _ f.Res, _ func() error)     {}
func (this *TestMethod) ExprNtVoidStrict(_ f.Req, _ f.Res, _ func() error)           {}
func (this *TestMethod) ExprNtNo(_ struct{}, _ f.Req, _ f.Res, _ func()) error       { return nil }
func (this *TestMethod) ExprNtNoStrict(_ f.Req, _ f.Res, _ func()) error             { return nil }
func (this *TestMethod) ExprNtNoVoid(_ struct{}, _ f.Req, _ f.Res, _ func())         {}
func (this *TestMethod) ExprNtNoVoidStrict(_ f.Req, _ f.Res, _ func())               {}
func (this *TestMethod) ExprNtErr(_ struct{}, _ f.Req, _ f.Res, _ FunErrErr) error   { return nil }
func (this *TestMethod) ExprNtErrStrict(_ f.Req, _ f.Res, _ FunErrErr) error         { return nil }
func (this *TestMethod) ExprNtErrVoid(_ struct{}, _ f.Req, _ f.Res, _ FunErrErr)     {}
func (this *TestMethod) ExprNtErrVoidStrict(_ f.Req, _ f.Res, _ FunErrErr)           {}
func (this *TestMethod) ExprNtErrNo(_ struct{}, _ f.Req, _ f.Res, _ FunErr) error    { return nil }
func (this *TestMethod) ExprNtErrNoStrict(_ f.Req, _ f.Res, _ FunErr) error          { return nil }
func (this *TestMethod) ExprNtErrNoVoid(_ struct{}, _ f.Req, _ f.Res, _ FunErr)      {}
func (this *TestMethod) ExprNtErrNoVoidStrict(_ f.Req, _ f.Res, _ FunErr)            {}
func (this *TestMethod) HttpHandlerFunc(_ struct{}) http.HandlerFunc                 { return httpHandler }
func (this *TestMethod) HttpHandlerFuncStrict() http.HandlerFunc                     { return httpHandler }
func (this *TestMethod) HttpHandler(_ struct{}) http.Handler                         { return httpHandler }
func (this *TestMethod) HttpHandlerStrict() http.Handler                             { return httpHandler }
func (this *TestMethod) HttpFunc(_ struct{}, _ http.ResponseWriter, _ *http.Request) {}
func (this *TestMethod) HttpFuncStrict(_ http.ResponseWriter, _ *http.Request)       {}
func (this *TestMethod) Fasthttp(_ struct{}, _ *fasthttp.RequestCtx)                 {}
func (this *TestMethod) FasthttpStrict(_ *fasthttp.RequestCtx)                       {}
func (this *TestMethod) FasthttpErr(_ struct{}, _ *fasthttp.RequestCtx) error        { return nil }
func (this *TestMethod) FasthttpErrStrict(_ *fasthttp.RequestCtx) error              { return nil }
func (this *TestMethod) SkipErr(_ error)                                             {}

var _ = Describe("ControllerMethod", func() {
	GetLogger().SetLogLevel(LevelDebug)
	t := reflect.TypeFor[*TestMethod]()
	ctrl := &TestMethod{}
	It("can check a method", func() {
		for method := range t.Methods() {
			if Strings(method.Name).HasPrefix("Skip") {
				continue
			}
			m, found := ControllerMethod{}.FromMethod(&method)
			Expect(found).To(BeTrue(), "Method not found: "+method.Name)
			Expect(m).NotTo(BeNil())
			name := method.Name
			defaults := []reflect.Value{}
			if p, ok := Strings(name).CutSuffix("Strict"); ok {
				name = p
			} else if Strings(name).HasPrefix("HttpHandler") {
				defaults = append(defaults, reflect.ValueOf(struct{}{}))
			}
			name = CaseUtils(name).SnakeCase().String()
			Expect(Strings(name).ToUpper()).To(Equal(m.Name()))
			Expect(m.Handler(defaults, ctrl)).NotTo(BeNil())
		}
	})
	It("can report error on bad methods", func() {
		s, ok := t.MethodByName("SkipErr")
		Expect(ok).To(BeTrue())
		m, found := ControllerMethod{}.FromMethod(&s)
		Expect(found).To(BeFalse())
		Expect(m).To(BeNil())
	})
})

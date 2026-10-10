package testing

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	"github.com/BeamFoundry/rays-fiber/pkg/template/html"
	"github.com/go-playground/validator/v10"

	// . "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v3"
)

var CsrfCookieName = "csrf_"

type CookieJar struct {
	jar map[string]*http.Cookie
}

func (this CookieJar) New() *CookieJar {
	this.jar = map[string]*http.Cookie{}
	return &this
}

func (this *CookieJar) ReadCookies(resp *http.Response) {
	for _, c := range resp.Cookies() {
		this.jar[c.Name] = c
	}
}

func (this *CookieJar) WriteCookies(req *http.Request) {
	for _, c := range this.jar {
		req.AddCookie(c)
	}
}

func (this *CookieJar) GetValue(name string) string {
	return this.jar[name].Value
}

func (this *CookieJar) RemoveCookie(name string) {
	delete(this.jar, name)
}

type TestClient struct {
	cookieJar *CookieJar
	app       *fiber.App
}

func (this TestClient) New(app *fiber.App) *TestClient {
	this.cookieJar = CookieJar{}.New()
	Expect(app).NotTo(BeNil())
	this.app = app
	return &this
}

func (this *TestClient) GET(uri string, reqfilters ...func(*http.Request)) *http.Response {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost"+uri, http.NoBody)
	this.cookieJar.WriteCookies(req)
	for _, f := range reqfilters {
		f(req)
	}
	resp, err := this.app.Test(req)
	Expect(err).NotTo(HaveOccurred())
	this.cookieJar.ReadCookies(resp)
	GetLogger().Debug("GET "+uri, "Req", req, "Resp", resp)
	return resp
}

func (this *TestClient) POST(uri string, form url.Values, reqfilters ...func(*http.Request)) *http.Response {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost"+uri, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Form = form
	req.PostForm = form
	this.cookieJar.WriteCookies(req)
	for _, f := range reqfilters {
		f(req)
	}
	resp, err := this.app.Test(req)
	Expect(err).NotTo(HaveOccurred())
	this.cookieJar.ReadCookies(resp)
	GetLogger().Debug("POST "+uri, "Req", req, "Resp", resp)
	return resp
}

func (this *TestClient) CookieJar() *CookieJar {
	return this.cookieJar
}

func (this *TestClient) GetCookieValue(name string) string {
	return this.cookieJar.GetValue(name)
}

//go:embed testdata
var testdata embed.FS

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
//
//nolint:wrapcheck // Need unwrap error
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type TestSecurityAppConfigCustomizer struct {
	Component
	structValidator Optional[fiber.StructValidator] `@:"Inject"`
}

var _ FiberConfig = &TestSecurityAppConfigCustomizer{}

func (this *TestSecurityAppConfigCustomizer) FiberConfig(config *fiber.Config) {
	GetLogger().Debug("TestSecurityAppConfigCustomizer::CustomizeConfig called !!!")
	config.Views = html.NewFileSystem(http.FS(testdata), ".html")
	config.PassLocalsToViews = true
	config.ErrorHandler = func(c fiber.Ctx, err error) error {
		GetLogger().Error("ErrorHandler", "error", err.Error())

		code := fiber.StatusInternalServerError

		// Retrieve the custom status code if it's a *fiber.Error
		if e, ok := errors.AsType[*fiber.Error](err); ok {
			code = e.Code
		}
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(code).SendString(err.Error())
	}
	if validation, ok := this.structValidator.Get(); ok {
		config.StructValidator = validation
	} else {
		config.StructValidator = &structValidator{validate: validator.New()}
	}
}

//---

type TestSecurityAppHandlers struct {
	FiberPathFilter
	Logger
}

func (this *TestSecurityAppHandlers) LoggerNext() func(fiber.Ctx) bool {
	return func(_ fiber.Ctx) bool {
		return true
	}
}

type TestSecurityPathAppHandlers struct {
	FiberPathFilter
	RequestMapping `Path:"/envvars"`
	EnvVar
}

package security

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	. "github.com/BeamFoundry/rays/pkg/lang"
	"github.com/valyala/fasthttp"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// FiberSameSite offers a http.SameSite (int) to Fiber's Cookie SameSite (string) conversion
type FiberSameSite struct {
	TitleEnumStringer[FiberSameSite]
	DISABLED, // 0 --> disabled
	DEFAULT, // 1 (SameSiteDefaultMode) --> disabled
	LAX, // 2 (SameSiteLaxMode)
	STRICT, // 3 (SameSiteStrictMode)
	NONE Enum[FiberSameSite] `Enum:"Ordinal"` // 4 (SameSiteNoneMode)
}

// Convert translate to string a http.SameSite value to Fiber's Cookie string value.
//
// Originally FiberHTTPHandlerAdapter.SetCookie() had to do this conversion using
// a switch case for each value.
//
// This "externalized" this conversion using the same code for each, easing the test for that.
//
//nolint:gocritic // constructor
func (this FiberSameSite) Convert(samesite http.SameSite) string {
	value := "disabled"
	switch samesite {
	case http.SameSiteLaxMode, http.SameSiteStrictMode, http.SameSiteNoneMode:
		cookie := Enum[FiberSameSite]{}
		//nolint:errcheck // not needed
		_ = cookie.UnmarshalText(fmt.Appendf(nil, "%d", samesite))
		value = cookie.String()
	default:
	}
	return value
}

type FiberSecurityRouter struct {
	router fiber.Router
}

func (this FiberSecurityRouter) New(app *fiber.App) *FiberSecurityRouter {
	this.router = app
	return &this
}

// Group implements [core.SecurityRouter].
func (this *FiberSecurityRouter) Group(path string) SecurityRouter {
	return &FiberSecurityRouter{router: this.router.Group(path)}
}

// Use implements [core.SecurityRouter].
func (this *FiberSecurityRouter) Use(handler func(SecurityHandler) error) {
	this.router.Use(func(ctx fiber.Ctx) error {
		return handler(FiberSecurityHandler{}.New(ctx))
	})
}

// Get implements [core.SecurityRouter].
func (this *FiberSecurityRouter) Get(path string, handler func(SecurityHandler) error) {
	this.router.Get(path, func(ctx fiber.Ctx) error {
		return handler(FiberSecurityHandler{}.New(ctx))
	})
}

// Post implements [core.SecurityRouter].
func (this *FiberSecurityRouter) Post(path string, handler func(SecurityHandler) error) {
	this.router.Post(path, func(ctx fiber.Ctx) error {
		return handler(FiberSecurityHandler{}.New(ctx))
	})
}

var _ SecurityRouter = &FiberSecurityRouter{}

// Fiber implementation of rays-web-security/core.HTTPHandler.
type FiberSecurityHandler struct {
	ctx fiber.Ctx
	req *http.Request
}

var _ SecurityHandler = &FiberSecurityHandler{}

var _ TLSConnectionHandler = &FiberSecurityHandler{}

// FiberSecurityHandler constructor. It needs a fiber.Ctx.
func (this FiberSecurityHandler) New(ctx fiber.Ctx) *FiberSecurityHandler {
	this.ctx = ctx
	//nolint:errcheck // Later
	// TODO : Need to log/handle this one.
	req, _ := adaptor.ConvertRequest(ctx, true)
	this.req = req
	return &this
}

// TLSConnectionState implements [core.TLSConnectionHandler].
func (this *FiberSecurityHandler) TLSConnectionState() *tls.ConnectionState {
	if this.req != nil && this.req.TLS != nil {
		return this.req.TLS
	}
	if this.ctx.RequestCtx() != nil {
		return this.ctx.RequestCtx().TLSConnectionState()
	}
	return nil
}

// Cookie implements rays-web-security/core.HTTPHandler.Cookie(string)
func (this *FiberSecurityHandler) Cookie(name string) (*http.Cookie, error) {
	c, err := this.req.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return c, nil
}

// FormValue implements rays-web-security/core.HTTPHandler.FormValue(string)
func (this *FiberSecurityHandler) FormValue(name string) string {
	return this.ctx.FormValue(name)
}

// SetCookie implements rays-web-security/core.HTTPHandler.SetCookie(*http.Cookie)
func (this *FiberSecurityHandler) SetCookie(cookie *http.Cookie) {
	c := &fiber.Cookie{
		Name:        cookie.Name,
		Value:       cookie.Value,
		Path:        cookie.Path,
		Domain:      cookie.Domain,
		MaxAge:      cookie.MaxAge,
		Expires:     cookie.Expires,
		Secure:      cookie.Secure,
		HTTPOnly:    cookie.HttpOnly,
		SameSite:    FiberSameSite{}.Convert(cookie.SameSite),
		Partitioned: cookie.Partitioned,
	}
	if c.Value == "" && (c.MaxAge < 0 || c.Expires.Equal(time.Unix(0, 0))) {
		c.Expires = fasthttp.CookieExpireDelete
	}
	this.ctx.Cookie(c)
}

// GetHeader implements rays-web-security/core.HTTPHandler.Header(string)
func (this *FiberSecurityHandler) Header(name string) string {
	return this.ctx.Get(name)
}

// SetHeader implements rays-web-security/core.HTTPHandler.SetHeader(string)
func (this *FiberSecurityHandler) SetHeader(name, value string) {
	this.ctx.Set(name, value)
}

// SessionValue implements rays-web-security/core.HTTPHandler.SessionValue(string)
func (this *FiberSecurityHandler) SessionValue(key string) any {
	sess := session.FromContext(this.ctx)
	return sess.Get(key)
}

func (this *FiberSecurityHandler) SetSessionValue(key string, obj any) {
	sess := session.FromContext(this.ctx)
	sess.Set(key, obj)
}

func (this *FiberSecurityHandler) SessionRefresh() {
	sess := session.FromContext(this.ctx)
	//nolint:errcheck // Later
	// TODO : Need to log/handle this one.
	sess.Regenerate()
}

func (this *FiberSecurityHandler) SessionReset() {
	sess := session.FromContext(this.ctx)
	//nolint:errcheck // Later
	// TODO : Need to log/handle this one.
	sess.Reset()
}

func (this *FiberSecurityHandler) URL() *url.URL {
	c := this.ctx.BaseURL() + this.ctx.OriginalURL()
	//nolint:errcheck // Later
	// TODO : Need to log/handle this one.
	u, _ := url.Parse(c)
	return u
}

func (this *FiberSecurityHandler) Method() string {
	return this.ctx.Method()
}

func (this *FiberSecurityHandler) Path() string {
	return this.ctx.Path()
}

// BindBody implements [rays-web-security/core.HTTPHandler].
func (this *FiberSecurityHandler) BindBody(obj any) error {
	return this.ctx.Bind().Body(obj)
}

// Error implements [core.SecurityHandler].
func (this *FiberSecurityHandler) Error(code int, message string) error {
	return fiber.NewError(code, message)
}

// Next implements [core.SecurityHandler].
func (this *FiberSecurityHandler) Next() error {
	return this.ctx.Next()
}

// Redirect implements [core.SecurityHandler].
func (this *FiberSecurityHandler) Redirect(to string, code int) error {
	return this.ctx.Redirect().Status(code).To(to)
}

// RenderView implements [core.SecurityHandler].
func (this *FiberSecurityHandler) RenderView(view string, data map[string]any) error {
	// FIXME: This is required, but should be removed once fixed.
	// For unknown reason, "Locals" are now passed down in ctx.Render unless
	// this is done here.
	if this.ctx.App().Config().PassLocalsToViews {
		this.ctx.Request().VisitUserValues(func(key []byte, value any) {
			if _, ok := data[string(key)]; !ok {
				data[string(key)] = value
			}
		})
	}
	return this.ctx.Render(view, data)
}

// Write implements [core.SecurityHandler].
func (this *FiberSecurityHandler) Write(p []byte) (n int, err error) {
	return this.ctx.Write(p)
}

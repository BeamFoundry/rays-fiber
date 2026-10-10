package security

import (
	"context"
	"errors"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

//nolint:govet // Singleton
type FiberSecurityAppHandler struct {
	Component
	Prototype

	securityContextHolder Optional[SecurityContextHolder] `@:"Inject"`
	csrfSecurityConfig    Optional[CsrfSecurityConfig]    `@:"Inject"`
	sessionSecurityConfig Optional[SessionSecurityConfig] `@:"Inject"`
	logger                Logging                         `@:"Inject"`

	csrfSkipper      func(SecurityHandler) bool
	sessionSkipper   func(SecurityHandler) bool
	sessionExtractor *extractors.Extractor

	instance string
	app      *fiber.App
	*FiberSecurityRouter

	Session
	Csrf
}

var _ SecurityAppHandler = &FiberSecurityAppHandler{}

var _ SessionConfig = &FiberSecurityAppHandler{}

var _ CsrfConfig = &FiberSecurityAppHandler{}

var _, _ = any(&FiberSecurityAppHandler{}).(Stereotype)

func (this *FiberSecurityAppHandler) WithContext(ctx context.Context) Error {
	if proxy := ctx.Value(PrototypeContextProxyType); proxy != nil {
		app, ok := proxy.(*fiber.App)
		if !ok || this.securityContextHolder.IsEmpty() {
			return NewError("SecurityStateError", "Unable to use securityContextHolder, HttpSecurity will not be enabled")
		}
		this.app = app
		this.FiberSecurityRouter = FiberSecurityRouter{}.New(this.app)
		if instance, ok := ctx.Value(PrototypeContextProxyInstance).(string); ok {
			this.logger.Debug("FiberSecurityAppHandler", "instance", instance)
			this.instance = instance
		}
	}

	return nil
}

// EnableCsrf implements [core.SecurityAppHandler].
func (this *FiberSecurityAppHandler) EnableCsrf(skipper func(SecurityHandler) bool) {
	this.csrfSkipper = skipper
	this.app.Use(this.Csrf.Handler())
	this.app.Use(func(c fiber.Ctx) error {
		if !this.CsrfNext()(c) {
			token := csrf.TokenFromContext(c)
			fiber.Locals(c, "Csrf", token)
		}
		return c.Next()
	})
}

// EnableSession implements [core.SecurityAppHandler].
func (this *FiberSecurityAppHandler) EnableSession(skipper func(SecurityHandler) bool) {
	this.sessionSkipper = skipper
	this.app.Use(this.Session.Handler())
}

// CsrfConfig implements [middleware.CsrfConfig].
func (this *FiberSecurityAppHandler) CsrfConfig(config *csrf.Config) {
	if cfg, ok := this.csrfSecurityConfig.Get(); ok {
		cfg.CsrfSecurityConfig(this.instance, config)
	}
}

// SessionConfig implements [middleware.SessionConfig].
func (this *FiberSecurityAppHandler) SessionConfig(config *session.Config) {
	if cfg, ok := this.sessionSecurityConfig.Get(); ok {
		cfg.SessionSecurityConfig(this.instance, config)
	}
	this.sessionExtractor = &config.Extractor
}

// HasSession implements [core.SecurityAppHandler].
func (this *FiberSecurityAppHandler) HasSession(handler SecurityHandler) bool {
	h := handler.(*FiberSecurityHandler)
	if _, err := this.sessionExtractor.Extract(h.ctx); err == nil {
		return true
	}
	return false
}

// PostAuthorizationFilterHook implements [core.SecurityAppHandler].
func (this *FiberSecurityAppHandler) PostAuthorizationFilter(handler SecurityHandler) {
	if h, ok := handler.(*FiberSecurityHandler); ok {
		if c, ok := this.securityContextHolder.Get(); ok {
			fiber.Locals(h.ctx, "httpPolicy", c.GetContext().GetHTTPPolicyDebugValue())
		}
	}
}

// TODO: Move that logic into SecurityFilterChain
func (this *FiberSecurityAppHandler) SessionNext() func(fiber.Ctx) bool {
	return func(ctx fiber.Ctx) bool {
		return this.sessionSkipper(FiberSecurityHandler{}.New(ctx))
	}
}

func (this *FiberSecurityAppHandler) CsrfExtractor() extractors.Extractor {
	formKeyLookup := "_csrf"
	headerKeyLookup := "X-Csrf-Token"

	return extractors.FromCustom(
		"FiberSecurityMiddlewares", func(c fiber.Ctx) (string, error) {
			if token := c.FormValue(formKeyLookup); token != "" {
				this.logger.Trace("CsrfExtractor() Form", "token", token, "cookie", c.Cookies("csrf_"))
				return token, nil
			} else if token := c.Get(headerKeyLookup); token != "" {
				this.logger.Trace("CsrfExtractor() Header", "token", token, "cookie", c.Cookies("csrf_"))
				return token, nil
			}
			return "", errors.New("CsrfExtractor() missing csrf token from form and header")
		},
	)
}

func (this *FiberSecurityAppHandler) CsrfErrorHandler() func(fiber.Ctx, error) error {
	return func(_ fiber.Ctx, err error) error {
		this.logger.Error("CsrfErrorHandler()", "error", err)
		return fiber.ErrForbidden
	}
}

// TODO: Move that logic into SecurityFilterChain
func (this *FiberSecurityAppHandler) CsrfNext() func(fiber.Ctx) bool {
	return func(ctx fiber.Ctx) bool {
		return this.csrfSkipper(FiberSecurityHandler{}.New(ctx))
	}
}

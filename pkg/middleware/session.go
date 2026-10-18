package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// SessionNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Session], to indicates when to skip it this middleware.
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionNext interface {
	SessionNext() func(fiber.Ctx) bool
}

// SessionErrorHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Session], to indicates a different error callback.
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionErrorHandler interface {
	SessionErrorHandler() func(fiber.Ctx, error)
}

// SessionKeyGenerator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Session], to indicates a different session ID generator.
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionKeyGenerator interface {
	SessionKeyGenerator() func() string
}

// SessionExtractor is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Session], to indicates a different session ID extrator.
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionExtractor interface {
	SessionExtractor() extractors.Extractor
}

// SessionConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Session], to review the [session.Config] that will be configured by [Session]
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionConfig interface {
	SessionConfig(config *session.Config)
}

// SessionSecurityConfig is an interface to be defined by a component that will customize
// the [Session] middleware when an HttpSecurity is defined for that server instance.
//
// Please see the [session] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
type SessionSecurityConfig interface {
	SessionSecurityConfig(server string, config *session.Config)
}

// SessionStore is an interface implemented by [Session] itself to give access to its Session store.
//
// It is referenced in the [Csrf] middleware, which may depend on it.
// Please see the [session] and [csrf] middleware configuration documentation.
//
// [session]: https://docs.gofiber.io/middleware/session/#config-properties
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
type SessionStore interface {
	SessionStore() *session.Store
}

type sessionConfigurationProperties struct {
	storage           *string
	cookieDomain      string
	cookiePath        string
	cookieSameSite    string
	idleTimeout       time.Duration
	absoluteTimeout   time.Duration
	cookieSecure      bool
	cookieHTTPOnly    bool
	cookieSessionOnly bool
}

type sessionConfig struct {
	session sessionConfigurationProperties
}

var _, _ = any(&FiberMwConfig[sessionConfig]{}).(Stereotype)

func (this *sessionConfig) Config() session.Config {
	config := session.ConfigDefault

	config.CookieDomain = this.session.cookieDomain
	config.CookiePath = this.session.cookiePath
	config.CookieSameSite = this.session.cookieSameSite
	config.IdleTimeout = this.session.idleTimeout
	config.AbsoluteTimeout = this.session.absoluteTimeout
	config.CookieSecure = this.session.cookieSecure
	config.CookieHTTPOnly = this.session.cookieHTTPOnly
	config.CookieSessionOnly = this.session.cookieSessionOnly

	config.Storage = nil
	return config
}

// [Session] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// NOTE: This middleware is marked with the [github.com/BeamFoundry/rays/pkg/core.PostAuthorisationFilter] interface.
// As such it will be skipped when used with [github.com/BeamFoundry/rays-web-security] but it will be configured and
// used automatically, if needed.
//
// WARNING: The [Csrf] middleware may on thiss middleware. As such the [Session] middleware MUST
// always be configured before [Csrf]. Please see the example below, the [session] and [csrf] Fiber documentation.
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Session
//	  Csrf
//	}
//
//	func (this *ApiController) SessionConfig(config *session.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [SessionNext], [SessionErrorHandler], [SessionKeyGenerator] and [SessionExtractor].
//
// Also, the parent container may implements [SessionConfig] to receive a pointer to the [session.Config]
// that will be used.
//
// Please see those interface definitions and the [session] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					session:
//	          storage: null
//	          cookieDomain: ""
//	          cookiePath: ""
//	          cookieSameSite: "Lax"
//	          idleTimeout: 30m
//	          absoluteTimeout: 0s
//	          cookieSecure: false
//	          cookieHTTPOnly: false
//	          cookieSessionOnly: false
//
// [csrf]: https://docs.gofiber.io/middleware/csrf/#config-properties
// [session]: https://docs.gofiber.io/middleware/session/#config
type Session struct {
	PostAuthorisationFilter
	config          *FiberMwConfig[sessionConfig] `@:"Inject"`
	storageProvider StorageProvider               `@:"Inject"`
	cfg             *session.Config
	handler         fiber.Handler
	store           *session.Store

	// initOnce sync.Once
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Session) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if cfg.session.storage != nil {
		this.cfg.Storage = this.storageProvider.Provide(*cfg.session.storage)
	}
	if next, ok := parent.(SessionNext); ok {
		this.cfg.Next = next.SessionNext()
	}
	if errorHandler, ok := parent.(SessionErrorHandler); ok {
		this.cfg.ErrorHandler = errorHandler.SessionErrorHandler()
	}
	if keyGenerator, ok := parent.(SessionKeyGenerator); ok {
		this.cfg.KeyGenerator = keyGenerator.SessionKeyGenerator()
	}
	if extractor, ok := parent.(SessionExtractor); ok {
		this.cfg.Extractor = extractor.SessionExtractor()
	}
	if config, ok := parent.(SessionConfig); ok {
		config.SessionConfig(this.cfg)
	}
	this.handler, this.store = session.NewWithStore(*this.cfg)
	return nil
}

// SessionStore() implements [SessionStore].
func (this *Session) SessionStore() *session.Store {
	return this.store
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Session) Handler() fiber.Handler {
	return this.handler
}

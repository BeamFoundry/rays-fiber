package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

// KeyAuthNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to indicates when to skip it this middleware.
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthNext interface {
	KeyAuthNext() func(fiber.Ctx) bool
}

// KeyAuthSuccessHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to indicates a function which is executed for a valid key.
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthSuccessHandler interface {
	KeyAuthSuccessHandler() fiber.Handler
}

// KeyAuthErrorHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to indicates a function which is executed for an invalid key.
// By default a 401 response with a WWW-Authenticate challenge is sent.
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthErrorHandler interface {
	KeyAuthErrorHandler() fiber.ErrorHandler
}

// KeyAuthValidator is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to indicates a function to validate the key.
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthValidator interface {
	KeyAuthValidator() func(fiber.Ctx, string) (bool, error)
}

// KeyAuthExtractor is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to indicates how to retrieve the key from the request.
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthExtractor interface {
	KeyAuthExtractor() extractors.Extractor
}

// KeyAuthConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [KeyAuth], to review the [keyauth.Config] that will be configured by [KeyAuth]
//
// Please see the [keyauth] middleware configuration documentation.
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuthConfig interface {
	KeyAuthConfig(config *keyauth.Config)
}

type keyAuthConfigurationProperties struct {
	realm            string
	challenge        string
	errorCode        string
	errorDescription string
	errorURI         string
	scope            string
}

type keyAuthConfig struct {
	keyauth keyAuthConfigurationProperties
}

var _, _ = any(&FiberMwConfig[keyAuthConfig]{}).(Stereotype)

func (this *keyAuthConfig) Config() keyauth.Config {
	config := keyauth.ConfigDefault

	config.Realm = this.keyauth.realm
	config.Challenge = this.keyauth.challenge
	config.Error = this.keyauth.errorCode
	config.ErrorDescription = this.keyauth.errorDescription
	config.ErrorURI = this.keyauth.errorURI
	config.Scope = this.keyauth.scope
	return config
}

// [KeyAuth] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// NOTE: This middleware is marked with the [github.com/BeamFoundry/rays/pkg/core.PreAuthenticationFilter] interface.
// As such it will be skipped when used with [github.com/BeamFoundry/rays-web-security].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  KeyAuth
//	}
//
//	func (this *ApiController) KeyAuthConfig(config *keyauth.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [KeyAuthNext], [KeyAuthSuccessHandler], [KeyAuthErrorHandler] and [KeyAuthValidator].
//
// Also, the parent container may implements [KeyAuthConfig] to receive a pointer to the [keyauth.Config]
// that will be used.
//
// Please see those interface definitions and the [keyauth] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					keyauth:
//						realm: Restricted
//						challenge: ApiKey realm="Restricted"
//						errorCode: "" # Allowed values are `invalid_request`, `invalid_token`, or `insufficient_scope`.
//						errorDescription: ""
//						errorURI: ""
//						scope: ""
//
// [keyauth]: https://docs.gofiber.io/middleware/keyauth/#config
type KeyAuth struct {
	PreAuthenticationFilter
	config *FiberMwConfig[keyAuthConfig] `@:"Inject"`
	cfg    *keyauth.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *KeyAuth) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(KeyAuthNext); ok {
		this.cfg.Next = next.KeyAuthNext()
	}
	if successHandler, ok := parent.(KeyAuthSuccessHandler); ok {
		this.cfg.SuccessHandler = successHandler.KeyAuthSuccessHandler()
	}
	if errorHandler, ok := parent.(KeyAuthErrorHandler); ok {
		this.cfg.ErrorHandler = errorHandler.KeyAuthErrorHandler()
	}
	if validator, ok := parent.(KeyAuthValidator); ok {
		this.cfg.Validator = validator.KeyAuthValidator()
	}
	if extractor, ok := parent.(KeyAuthExtractor); ok {
		this.cfg.Extractor = extractor.KeyAuthExtractor()
	}
	if config, ok := parent.(KeyAuthConfig); ok {
		config.KeyAuthConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *KeyAuth) Handler() fiber.Handler {
	return keyauth.New(*this.cfg)
}

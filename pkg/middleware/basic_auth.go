package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

// BasicAuthNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [BasicAuth], to indicates when to skip it this middleware.
//
// Please see the [basicauth] middleware configuration documentation.
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuthNext interface {
	BasicAuthNext() func(fiber.Ctx) bool
}

// BasicAuthAuthorizer is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [BasicAuth], to indicates a function to check the credentials.
// It will be called with a username, password, and the current context
// and is expected to return true or false to indicate approval.
//
// Please see the [basicauth] middleware configuration documentation.
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuthAuthorizer interface {
	BasicAuthAuthorizer() func(string, string, fiber.Ctx) bool
}

// BasicAuthUnauthorized is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [BasicAuth], to indicates a [fiber.Handler] to generate a response body for unauthorized
// requests.
//
// Please see the [basicauth] middleware configuration documentation.
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuthUnauthorized interface {
	BasicAuthUnauthorized() fiber.Handler
}

// BasicAuthBadRequest is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [BasicAuth], to indicates a [fiber.Handler] to generate a response for malformed Authorization headers.
//
// Please see the [basicauth] middleware configuration documentation.
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuthBadRequest interface {
	BasicAuthBadRequest() fiber.Handler
}

// BasicAuthConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [BasicAuth], to review the [basicauth.Config] that will be configured by [BasicAuth]
//
// Please see the [basicauth] middleware configuration documentation.
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuthConfig interface {
	BasicAuthConfig(config *basicauth.Config)
}

type basicAuthConfigurationProperties struct {
	users       map[string]*string
	charset     string
	realm       string
	headerLimit int
}

type basicAuthConfig struct {
	basicauth basicAuthConfigurationProperties
}

var _, _ = any(&FiberMwConfig[basicAuthConfig]{}).(Stereotype)

func (this *basicAuthConfig) Config() basicauth.Config {
	config := basicauth.ConfigDefault

	users := map[string]string{}
	for u, p := range this.basicauth.users {
		if p != nil {
			users[u] = *p
		}
	}
	config.Users = users
	config.Realm = this.basicauth.realm
	config.Charset = this.basicauth.charset
	config.HeaderLimit = this.basicauth.headerLimit

	return config
}

// [BasicAuth] is a middleware to be defined directly as an anonymous member inside a
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
//	  BasicAuth
//	}
//
//	func (this *ApiController) BasicAuthConfig(config *basicauth.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [BasicAuthNext], [BasicAuthAuthorizer], [BasicAuthUnauthorized] and [BasicAuthBadRequest].
//
// Also, the parent container may implements [BasicAuthConfig] to receive a pointer to the [basicauth.Config]
// that will be used.
//
// Please see those interface definitions and the [basicauth] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					basicauth:
//						users: # dictionary of username/password
//							admin: null # Note: null passworded users will be ignored.
//						charset: "UTF-8"
//						realm: "Restricted"
//						headerLimit: 8192
//
// [basicauth]: https://docs.gofiber.io/middleware/basicauth/#config
type BasicAuth struct {
	PreAuthenticationFilter
	config *FiberMwConfig[basicAuthConfig] `@:"Inject"`
	cfg    *basicauth.Config
}

var _ FiberHandler = &BasicAuth{}

var _ InnerBeamPostConstruct = &BasicAuth{}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *BasicAuth) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := parent.(BasicAuthNext); ok {
		this.cfg.Next = next.BasicAuthNext()
	}
	if authorizer, ok := parent.(BasicAuthAuthorizer); ok {
		this.cfg.Authorizer = authorizer.BasicAuthAuthorizer()
	}
	if unauthorized, ok := parent.(BasicAuthUnauthorized); ok {
		this.cfg.Unauthorized = unauthorized.BasicAuthUnauthorized()
	}
	if badrequest, ok := parent.(BasicAuthBadRequest); ok {
		this.cfg.BadRequest = badrequest.BasicAuthBadRequest()
	}
	if config, ok := parent.(BasicAuthConfig); ok {
		config.BasicAuthConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *BasicAuth) Handler() fiber.Handler {
	return basicauth.New(*this.cfg)
}

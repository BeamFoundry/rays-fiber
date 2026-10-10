package core

import (
	"github.com/gofiber/fiber/v3"

	. "github.com/BeamFoundry/rays/pkg/core"
)

// FiberController is a Stereotype marker interface representing a Controller type.
//
// It's an equivalent to a fiber Group as it reference a URL via a RequestMapping marker interface,
// it hosts a list of middleware embedded directly on its type and define handlers with its methods.
//
// As with other Stereotype, it supports dependency injection and post-construct initialization.
// Also it must be registered with a type assertion to be discovered.
//
// As with its [FiberPathFilter] counterpart it can hosts middlewares, [FiberHandler] embedded,
// to be configured using methods implementing interface on it.
//
// On methods defining endpoint the RequestMapping must be indicated in an anonymous struct, as its
// first argument. That struct can contain other middleware [FiberHandler] and injection.
//
// All [Fiber handlers] type are supported but the function arguments must be on the [FiberController]
// method argument, not on a returned function. See example below.
//
// NOTE: As most component linked to [fiber], it supports the [core.Server] annotation to register
// [FiberController] on a different server than the one defined in "primary" configuration.
//
// Example:
//
//	type UserController struct {
//	  FiberPathFilter
//	  RequestMapping `Path:"/user"`
//	  Logger // using the Logger middleware
//	  userRepository UserRepository `@:"Inject"`
//	}
//
//	// LoggerConfig implements [middleware.LoggerConfig]
//	func (this *UserController) LoggerConfig(config *logger.Config) {
//	  // changing the default logging configuration
//	}
//
//	// UserDetails endpoint let change the user informations
//	//
//	// Endpoint will be on "/user/details" path and the method is "POST"
//	func (this *UserController) UserDetails(
//	  _ struct{
//	    RequestMapping `Path:"/details" Method:"POST"`
//	  },
//	  ctx fiber.Ctx,
//	) error {
//	  // Save user details
//	}
//
// [Fiber handlers]: https://docs.gofiber.io/guide/routing#handler-types
type FiberController interface {
	Controller
	_FiberController()
}

// FiberPathFilter is a Stereotype marker interface representing a specialized Component type.
//
// Its role is to define a list of middleware to apply to a path prefix.
//
// As with its [FiberController] counterpart it can hosts middlewares, [FiberHandler] embedded,
// to be configured using methods implementing interface on it.
//
// NOTE: As most component linked to [fiber], it supports the [core.Server] annotation to register
// [FiberController] on a different server than the one defined in "primary" configuration.
//
// NOTE: If no RequestMapping is used the middlewares will be registered for all [FiberController], like done by [fiber.App.Use],
// with no path as first argument.
//
// NOTE: If a RequestMapping is used on it with only a "Path" struct tag, it will act as with [fiber.App.Use] with a path prefix.
//
// WARNING: If a RequestMapping is used with both "Path" and "Method" struct tags, the path argument will NOT be interpreted as
// a path prefix like with [fiber.App.Use] but as a [Fiber endpoint path] and should contain a wildcard if not used for a
// specific endpoint.
//
// Example:
//
//	type AppPathFilter struct {
//	  FiberPathFilter
//	  RequestMapping `Path:"/api"`
//	  Logger
//	}
//
//	//  LoggerConfig implements [middleware.LoggerConfig]
//	func (this *AppPathFilter) LoggerConfig(config *logger.Config) {
//	  // changing the default logging configuration
//	}
//
// [Fiber endpoint path]: https://docs.gofiber.io/guide/routing#paths
type FiberPathFilter interface {
	Component
	_FiberPathFilter()
}

// FiberAppProvider is an interface to an internal object, used to define, customize, start and shutdown [fiber.App] objects.
// It is mostly made accessible to access the [fiber.App] in testing.
//
// NOTE: As with all objects it has to be injected, or a reference be obtained in testing with GetScopedBeamInstance.
//
// Example:
//
//	appProvider := GetBeamScopedInstance[FiberAppProvider]()
type FiberAppProvider interface {
	// Provide() will return the [fiber.App] instance.
	Provide(name string) *fiber.App

	// Default() will return the default [fiber.App], like Provide("primary") would.
	Default() *fiber.App

	// Shutdown() will end all instances.
	Shutdown()
}

// FiberConfig is an interface a [Component] must implements to have an object that can change the default [fiber.Config].
//
// NOTE: As most component linked to [fiber], it supports the [core.Server] annotation to register
// [FiberConfig] on a different server than the one defined in "primary" configuration.
type FiberConfig interface {
	FiberConfig(*fiber.Config)
}

// FiberListenConfig is an interface a [Component] must implements to have an object that can change the default [fiber.ListenConfig].
//
// NOTE: As most component linked to [fiber], it supports the [core.Server] annotation to register
// [FiberListenConfig] on a different server than the one defined in "primary" configuration.
type FiberListenConfig interface {
	FiberListenConfig(*fiber.ListenConfig)
}

// FiberHandler is the interface a middleware must implements to be embedded in a [FiberController] or [FiberPathFilter].
//
// NOTE: Fiber middlewares have been adapted in custom structs already. Please see the middleware package.
type FiberHandler interface {
	Handler() fiber.Handler
}

// FiberFilterHandler is the interface describing a middleware that register around a fiber.Handler.
//
// NOTE: The Timeout middleware is only [FiberFilterHandler] provided.
//
// WARNING: FiberFilterHandler can only be used is [FiberController].
type FiberFilterHandler interface {
	Filter(fiber.Handler) fiber.Handler
}

// StorageEngine is the interface a [fiber.Storage] adapter should implement to be registered inside a [StorageProvider].
//
// NOTE: The name of the [StorageEngine] driver is taken from its Qualifier.
type StorageEngine interface {
	Storage(instance string) fiber.Storage
}

// StorageProvider is the interface of an internal [Service] that detects [StorageEngine] and can provide them to middlewares.
//
// Each storage drivers register a default instance (named "primary") associated to its name. Middleware configuration
// should include a string field indicating the storage driver name and optionally followed by a colon ":" and an
// instance name. For example : "memory:primary".
//
// Using this syntax and an injected [StorageProvider] should be enough for a middleware requiring a fiber.Storage to
// be configured.
type StorageProvider interface {
	Provide(name string) fiber.Storage
}

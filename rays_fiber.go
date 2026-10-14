// # Initialization
//
// This package should be imported from a main file, alongside rays :
//
//	package main
//
//	import (
//	  "github.com/BeamFoundry/rays"
//	  _ "github.com/BeamFoundry/rays-fiber"
//	)
//
// # Packages
//
// This project divides its code between the following packages :
//
//   - [github.com/BeamFoundry/rays-fiber/pkg/core]: Contains the interfaces definintions
//   - [github.com/BeamFoundry/rays-fiber/pkg/middleware]: Contains all the Fiber middlewares, directly defined in one package.
//   - [github.com/BeamFoundry/rays-fiber/pkg/storage/drivers]: Contains the root package for all the storage drivers.
//   - [github.com/BeamFoundry/rays-fiber/pkg/template/html]: Contains a revamped html templating module for Fiber.
//   - [github.com/BeamFoundry/rays-fiber/pkg/testing]: Contains testing objects used through the packages.
package rays_fiber //nolint:revive // Ignoring package naming hint

import (
	"embed"

	"github.com/BeamFoundry/rays"

	// Needed to initialize rays-fiber
	_ "github.com/BeamFoundry/rays-fiber/internal/controllers"
	_ "github.com/BeamFoundry/rays-fiber/pkg/core"
	_ "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage"
)

//go:embed config/*.yaml
var config embed.FS

func init() {
	//nolint:errcheck // No need to test
	rays.ApplicationContext().AddDefaultConfigFS(config)
}

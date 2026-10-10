package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/envvar"
)

// EnvVarConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EnvVar], to review the [envvar.Config] that will be configured by [EnvVar]
//
// Please see the [envvar] middleware configuration documentation.
//
// [envvar]: https://docs.gofiber.io/middleware/envvar/#config
type EnvVarConfig interface {
	EnvVarConfig(config *envvar.Config)
}

type envvarConfigurationProperties struct {
	exportVars map[string]*string
}

type envvarConfig struct {
	envvar envvarConfigurationProperties
}

var _, _ = any(&FiberMwConfig[envvarConfig]{}).(Stereotype)

func (this *envvarConfig) Config() envvar.Config {
	exported := map[string]string{}
	for k, v := range this.envvar.exportVars {
		if v != nil {
			exported[k] = *v
		}
	}

	config := envvar.Config{ExportVars: exported}
	return config
}

// [EnvVar] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  EnvVar
//	}
//
//	func (this *ApiController) EnvVarConfig(config *envvar.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties.
//
// Also, the parent container may implements [EnvVarConfig] to receive a pointer to the [envvar.Config]
// that will be used.
//
// Please see those interface definitions and the [envvar] middleware documentation.
//
// # Properties mapping :
//
//	 rays:
//	 	fiber:
//	 		servers:
//	 			primary:
//	 				envvar:
//						exportVars: {"HOME":null} # Using null as a map value will not register it. Use "" to have an empty value
//
// [envvar]: https://docs.gofiber.io/middleware/envvar/#config
type EnvVar struct {
	config *FiberMwConfig[envvarConfig] `@:"Inject"`
	cfg    *envvar.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *EnvVar) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if config, ok := (parent).(EnvVarConfig); ok {
		config.EnvVarConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *EnvVar) Handler() fiber.Handler {
	return envvar.New(*this.cfg)
}

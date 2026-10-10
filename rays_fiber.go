package rays_fiber //nolint:revive // TODO: See if correctable.

import (
	"embed"

	"github.com/BeamFoundry/rays"

	// Needed to initialize rays-fiber
	_ "github.com/BeamFoundry/rays-fiber/internal/controllers"
	_ "github.com/BeamFoundry/rays-fiber/pkg/core"
	_ "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage"
	//	_ "github.com/BeamFoundry/rays-fiber/pkg/security"
)

//go:embed config/*.yaml
var config embed.FS

func init() {
	//nolint:errcheck // No need to test
	rays.ApplicationContext().AddDefaultConfigFS(config)
}

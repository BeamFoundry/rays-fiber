package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	"github.com/gofiber/fiber/v3/middleware/envvar"
)

type EnvvarTestController struct {
	Controller
	RequestMapping `Path:"/testing/envvar"`
	EnvVar
}

var _ EnvVarConfig = &EnvvarTestController{}

// EnvVarConfig implements [middleware.EnvVarConfig].
func (e *EnvvarTestController) EnvVarConfig(_ *envvar.Config) {}

var _ = Describe("Envvar", func() {
	exported := map[string]any{
		"exportVars": map[string]any{
			"HOME": "/",
		},
	}
	It("can use Envvar", func() {
		Expect(GetContainer().LoadConfigAt(exported, "rays.fiber.servers.primary.envvar")).NotTo(HaveOccurred())
		AddBeamFor(&EnvvarTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

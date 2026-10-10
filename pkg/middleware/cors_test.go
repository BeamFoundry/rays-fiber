package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type CorsPathFilter struct {
	FiberPathFilter
	Cors
}

var _ CorsNext = &CorsPathFilter{}

var _ CorsAllowOriginsFunc = &CorsPathFilter{}

var _ CorsConfig = &CorsPathFilter{}

// CorsNext() implements [middleware.CorsNext]
func (this *CorsPathFilter) CorsNext() func(fiber.Ctx) bool {
	return nil
}

// CorsAllowOriginsFunc() implements [middleware.CorsAllowOriginsFunc()]
func (this *CorsPathFilter) CorsAllowOriginsFunc() func(origin string) bool {
	return nil
}

// CorsConfig() implements [middleware.CorsConfig].
func (this *CorsPathFilter) CorsConfig(_ *cors.Config) {
}

var _ = Describe("Cors", func() {
	It("can use Cors", func() {
		AddBeamFor(&CorsPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

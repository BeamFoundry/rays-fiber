package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	. "github.com/gofiber/fiber/v3/middleware/helmet"
)

type HelmetPathFilter struct {
	FiberPathFilter
	Helmet
}

var _ HelmetNext = &HelmetPathFilter{}

var _ HelmetConfig = &HelmetPathFilter{}

// HelmetNext implements [middleware.HelmetNext].
func (h *HelmetPathFilter) HelmetNext() func(fiber.Ctx) bool {
	return nil
}

// HelmetConfig implements [middleware.HelmetConfig].
func (h *HelmetPathFilter) HelmetConfig(_ *Config) {}

var _ = Describe("Helmet", func() {
	It("can use Helmet", func() {
		AddBeamFor(&HelmetPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

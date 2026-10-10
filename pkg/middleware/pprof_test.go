package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
)

type PprofPathFilter struct {
	FiberPathFilter
	Pprof
}

var _ PprofNext = &PprofPathFilter{}

var _ PprofConfig = &PprofPathFilter{}

// PprofNext implements [middleware.PprofNext].
func (p *PprofPathFilter) PprofNext() func(fiber.Ctx) bool {
	return nil
}

// PprofConfig implements [middleware.PprofConfig].
func (p *PprofPathFilter) PprofConfig(_ *pprof.Config) {}

var _ = Describe("Pprof", func() {
	It("can use Pprof", func() {
		AddBeamFor(&PprofPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

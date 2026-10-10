package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/etag"
)

type EtagPathFilter struct {
	FiberPathFilter
	Etag
}

// EtagNext implements [middleware.EtagNext].
func (this *EtagPathFilter) EtagNext() func(fiber.Ctx) bool {
	return nil
}

// EtagConfig implements [middleware.EtagConfig].
func (this *EtagPathFilter) EtagConfig(_ *etag.Config) {}

var _ EtagNext = &EtagPathFilter{}

var _ EtagConfig = &EtagPathFilter{}

var _ = Describe("Etag", func() {
	It("can use Etag", func() {
		AddBeamFor(&EtagPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

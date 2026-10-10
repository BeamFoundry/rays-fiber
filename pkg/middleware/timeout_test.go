package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"
)

type TimeoutTestController struct {
	Controller
	RequestMapping `Path:"/testing"`
	Timeout        `Duration:"10s"`
}

var _ TimeoutNext = &TimeoutTestController{}

var _ TimeoutErrors = &TimeoutTestController{}

var _ TimeoutConfig = &TimeoutTestController{}

// TimeoutNext implements [middleware.TimeoutNext].
func (this *TimeoutTestController) TimeoutNext() func(fiber.Ctx) bool {
	return nil
}

// TimeoutErrors implements [middleware.TimeoutErrors].
func (this *TimeoutTestController) TimeoutErrors() []error {
	return nil
}

// TimeoutConfig implements [middleware.TimeoutConfig].
func (this *TimeoutTestController) TimeoutConfig(_ *timeout.Config) {}

func (this *TimeoutTestController) TimeoutTest(
	_ struct {
		RequestMapping `Path:"/timeout" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString("hello")
}

var _ = Describe("Timeout", func() {
	It("can use Timeout", func() {
		AddBeamFor(&TimeoutTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

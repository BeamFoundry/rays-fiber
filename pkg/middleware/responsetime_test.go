package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/responsetime"
)

type ResponseTimePathFilter struct {
	FiberPathFilter
	ResponseTime
}

var _ ResponseTimeNext = &ResponseTimePathFilter{}

var _ ResponseTimeConfig = &ResponseTimePathFilter{}

// ResponseTimeNext() implements [middleware.ResponseTimeNext]
func (this *ResponseTimePathFilter) ResponseTimeNext() func(fiber.Ctx) bool {
	return nil
}

// ResponseTimeConfig() implements [middleware.ResponseTimeConfig].
func (this *ResponseTimePathFilter) ResponseTimeConfig(_ *responsetime.Config) {}

var _ = Describe("ResponseTime", func() {
	It("can use ResponseTime", func() {
		AddBeamFor(&ResponseTimePathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

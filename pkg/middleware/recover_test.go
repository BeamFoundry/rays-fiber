package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type RecoverPathFilter struct {
	FiberPathFilter
	Recover
}

// RecoverNext implements [middleware.RecoverNext].
func (r *RecoverPathFilter) RecoverNext() func(fiber.Ctx) bool {
	return nil
}

// RecoverPanicHandler implements [middleware.RecoverPanicHandler].
func (r *RecoverPathFilter) RecoverPanicHandler() func(c fiber.Ctx, e any) error {
	return nil
}

// RecoverStackTraceHandler implements [middleware.RecoverStackTraceHandler].
func (r *RecoverPathFilter) RecoverStackTraceHandler() func(c fiber.Ctx, e any) {
	return nil
}

// RecoverConfig implements [middleware.RecoverConfig].
func (r *RecoverPathFilter) RecoverConfig(_ *recover.Config) {}

var _ RecoverNext = &RecoverPathFilter{}

var _ RecoverPanicHandler = &RecoverPathFilter{}

var _ RecoverStackTraceHandler = &RecoverPathFilter{}

var _ RecoverConfig = &RecoverPathFilter{}

var _ = Describe("Recover", func() {
	It("can use Recover", func() {
		AddBeamFor(&RecoverPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

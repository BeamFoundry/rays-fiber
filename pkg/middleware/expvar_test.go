package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/expvar"
)

type ExpVarPathFilter struct {
	FiberPathFilter
	ExpVar
}

var _ ExpVarNext = &ExpVarPathFilter{}

var _ ExpVarConfig = &ExpVarPathFilter{}

// ExpVarNext implements [middleware.ExpVarNext].
func (e *ExpVarPathFilter) ExpVarNext() func(fiber.Ctx) bool {
	return nil
}

// ExpVarConfig implements [middleware.ExpVarConfig].
func (e *ExpVarPathFilter) ExpVarConfig(_ *expvar.Config) {}

var _ = Describe("ExpVar", func() {
	It("can use ExpVar", func() {
		AddBeamFor(&ExpVarPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

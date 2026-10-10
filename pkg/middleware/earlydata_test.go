package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v3"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3/middleware/earlydata"
)

type EarlyDataPathFilter struct {
	FiberPathFilter
	EarlyData
}

var _ EarlyDataNext = &EarlyDataPathFilter{}

var _ EarlyDataIsEarlyData = &EarlyDataPathFilter{}

var _ EarlyDataAllowEarlyData = &EarlyDataPathFilter{}

var _ EarlyDataError = &EarlyDataPathFilter{}

var _ EarlyDataConfig = &EarlyDataPathFilter{}

// EarlyDataNext implements [middleware.EarlyDataNext].
func (e *EarlyDataPathFilter) EarlyDataNext() func(fiber.Ctx) bool {
	return nil
}

// EarlyDataIsEarlyData implements [middleware.EarlyDataIsEarlyData].
func (e *EarlyDataPathFilter) EarlyDataIsEarlyData() func(c fiber.Ctx) bool {
	return nil
}

// EarlyDataAllowEarlyData implements [middleware.EarlyDataAllowEarlyData].
func (e *EarlyDataPathFilter) EarlyDataAllowEarlyData() func(c fiber.Ctx) bool {
	return nil
}

// EarlyDataError implements [middleware.EarlyDataError].
func (e *EarlyDataPathFilter) EarlyDataError() error {
	return nil
}

// EarlyDataConfig implements [middleware.EarlyDataConfig].
func (e *EarlyDataPathFilter) EarlyDataConfig(_ *earlydata.Config) {}

var _ = Describe("EarlyData", func() {
	It("can use EarlyData", func() {
		AddBeamFor(&EarlyDataPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

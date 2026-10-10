package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

type HealthCheckPathFilter struct {
	FiberPathFilter
	HealthCheck
}

var _ HealthCheckNext = &HealthCheckPathFilter{}

var _ HealthCheckProbe = &HealthCheckPathFilter{}

var _ HealthCheckConfig = &HealthCheckPathFilter{}

// HealthCheckNext implements [middleware.HealthCheckNext].
func (h *HealthCheckPathFilter) HealthCheckNext() func(fiber.Ctx) bool {
	return nil
}

// HealthCheckProbe implements [middleware.HealthCheckProbe].
func (h *HealthCheckPathFilter) HealthCheckProbe() func(fiber.Ctx) bool {
	return nil
}

// HealthCheckConfig implements [middleware.HealthCheckConfig].
func (h *HealthCheckPathFilter) HealthCheckConfig(_ *healthcheck.Config) {}

var _ = Describe("HealthCheck", func() {
	It("can use HealthCheck", func() {
		AddBeamFor(&HealthCheckPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

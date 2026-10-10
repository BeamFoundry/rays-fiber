package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/hostauthorization"
)

type HostAuthorizationPathFilter struct {
	FiberPathFilter
	HostAuthorization
}

var _ HostAuthorizationNext = &HostAuthorizationPathFilter{}

var _ HostAuthorizationAllowedHostsFunc = &HostAuthorizationPathFilter{}

var _ HostAuthorizationErrorHandler = &HostAuthorizationPathFilter{}

var _ HostAuthorizationConfig = &HostAuthorizationPathFilter{}

// HostAuthorizationNext() implements [middleware.HostAuthorizationNext]
func (this *HostAuthorizationPathFilter) HostAuthorizationNext() func(fiber.Ctx) bool {
	return nil
}

// HostAuthorizationAllowedHostsFunc implements [middleware.HostAuthorizationAllowedHostsFunc].
func (this *HostAuthorizationPathFilter) HostAuthorizationAllowedHostsFunc() func(host string) bool {
	return nil
}

// HostAuthorizationErrorHandler implements [middleware.HostAuthorizationErrorHandler].
func (this *HostAuthorizationPathFilter) HostAuthorizationErrorHandler() fiber.ErrorHandler {
	return nil
}

// HostAuthorizationConfig() implements [middleware.HostAuthorizationConfig].
func (this *HostAuthorizationPathFilter) HostAuthorizationConfig(_ *hostauthorization.Config) {}

var _ = Describe("HostAuthorization", func() {
	config := map[string]any{
		"allowedHosts": []any{"example.com", "*.example.com"},
	}
	It("can use HostAuthorization", func() {
		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.hostauthorization")).NotTo(HaveOccurred())
		AddBeamFor(&HostAuthorizationPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

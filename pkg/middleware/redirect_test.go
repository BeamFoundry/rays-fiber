package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/redirect"
)

type RedirectPathFilter struct {
	FiberPathFilter
	Redirect
}

var _ RedirectNext = &RedirectPathFilter{}

var _ RedirectConfig = &RedirectPathFilter{}

// RedirectNext implements [middleware.RedirectNext].
func (r *RedirectPathFilter) RedirectNext() func(fiber.Ctx) bool {
	return nil
}

// RedirectConfig implements [middleware.RedirectConfig].
func (r *RedirectPathFilter) RedirectConfig(_ *redirect.Config) {}

var _ = Describe("Redirect", func() {
	It("can use Redirect", func() {
		cfg := map[string]any{
			"rules": map[string]any{
				"/old/*": "/new/$1",
			},
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.redirect")).NotTo(HaveOccurred())
		AddBeamFor(&RedirectPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

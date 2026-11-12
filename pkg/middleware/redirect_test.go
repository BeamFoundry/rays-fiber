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
	capturedConfig *redirect.Config
}

var _ RedirectNext = &RedirectPathFilter{}

var _ RedirectConfig = &RedirectPathFilter{}

// RedirectNext implements [middleware.RedirectNext].
func (r *RedirectPathFilter) RedirectNext() func(fiber.Ctx) bool {
	return nil
}

// RedirectConfig implements [middleware.RedirectConfig].
func (r *RedirectPathFilter) RedirectConfig(cfg *redirect.Config) {
	r.capturedConfig = cfg
}

var _ = Describe("Redirect", func() {
	It("can use Redirect", func() {
		cfg := map[string]any{
			"rules": []any{
				map[string]any{
					"/old/*": "/new/$1",
				},
				map[string]any{
					"/ignored/*": nil,
				},
			},
			"statusCode": 301,
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.redirect")).NotTo(HaveOccurred())
		filter := &RedirectPathFilter{}
		AddBeamFor(filter)
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
		Expect(filter.capturedConfig).NotTo(BeNil())
		Expect(filter.capturedConfig.StatusCode).To(Equal(301))
		Expect(filter.capturedConfig.RuleList).To(HaveLen(1))
		Expect(filter.capturedConfig.RuleList[0].From).To(Equal("/old/*"))
		Expect(filter.capturedConfig.RuleList[0].To).To(Equal("/new/$1"))
	})
})

package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/rewrite"
)

type RewritePathFilter struct {
	FiberPathFilter
	Rewrite
	capturedConfig *rewrite.Config
}

// RewriteNext implements [middleware.RewriteNext].
func (r *RewritePathFilter) RewriteNext() func(fiber.Ctx) bool {
	return nil
}

// RewriteConfig implements [middleware.RewriteConfig].
func (r *RewritePathFilter) RewriteConfig(cfg *rewrite.Config) {
	r.capturedConfig = cfg
}

var _ RewriteNext = &RewritePathFilter{}

var _ RewriteConfig = &RewritePathFilter{}

var _ = Describe("Rewrite", func() {
	It("can use Rewrite", func() {
		cfg := map[string]any{
			"rules": []any{
				map[string]any{
					"/old": "/new",
				},
				map[string]any{
					"/ignored": nil,
				},
			},
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.rewrite")).NotTo(HaveOccurred())
		filter := &RewritePathFilter{}
		AddBeamFor(filter)
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
		Expect(filter.capturedConfig).NotTo(BeNil())
		Expect(filter.capturedConfig.RuleList).To(HaveLen(1))
		Expect(filter.capturedConfig.RuleList[0].From).To(Equal("/old"))
		Expect(filter.capturedConfig.RuleList[0].To).To(Equal("/new"))
	})
})

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
}

// RewriteNext implements [middleware.RewriteNext].
func (r *RewritePathFilter) RewriteNext() func(fiber.Ctx) bool {
	return nil
}

// RewriteConfig implements [middleware.RewriteConfig].
func (r *RewritePathFilter) RewriteConfig(_ *rewrite.Config) {}

var _ RewriteNext = &RewritePathFilter{}

var _ RewriteConfig = &RewritePathFilter{}

var _ = Describe("Rewrite", func() {
	It("can use Rewrite", func() {
		cfg := map[string]any{
			"rules": map[string]any{
				"/old": "/new",
			},
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.rewrite")).NotTo(HaveOccurred())
		AddBeamFor(&RewritePathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

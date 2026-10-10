package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type LimiterPathFilter struct {
	FiberPathFilter
	Limiter
}

// LimiterNext implements [middleware.LimiterNext].
func (l *LimiterPathFilter) LimiterNext() func(fiber.Ctx) bool {
	return nil
}

// LimiterMaxFunc implements [middleware.LimiterMaxFunc].
func (l *LimiterPathFilter) LimiterMaxFunc() func(c fiber.Ctx) int {
	return nil
}

// LimiterKeyGenerator implements [middleware.LimiterKeyGenerator].
func (l *LimiterPathFilter) LimiterKeyGenerator() func(fiber.Ctx) string {
	return nil
}

// LimiterLimitReached implements [middleware.LimiterLimitReached].
func (l *LimiterPathFilter) LimiterLimitReached() fiber.Handler {
	return nil
}

// LimiterConfig implements [middleware.LimiterConfig].
func (l *LimiterPathFilter) LimiterConfig(_ *limiter.Config) {}

var _ LimiterNext = &LimiterPathFilter{}

var _ LimiterMaxFunc = &LimiterPathFilter{}

var _ LimiterKeyGenerator = &LimiterPathFilter{}

var _ LimiterLimitReached = &LimiterPathFilter{}

var _ LimiterConfig = &LimiterPathFilter{}

var _ = Describe("Limiter", func() {
	storage := map[string]any{
		"storage": "memory",
	}
	testLimiter := func(h string) {
		handler := map[string]any{
			"handler": h,
		}
		Expect(GetContainer().LoadConfigAt(storage, "rays.fiber.servers.primary.limiter")).NotTo(HaveOccurred())
		Expect(GetContainer().LoadConfigAt(handler, "rays.fiber.servers.primary.limiter")).NotTo(HaveOccurred())
		AddBeamFor(&LimiterPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	}
	Context("Fixed Window", func() {
		It("can use Limiter", func() {
			testLimiter("fixed")
		})
	})
	Context("Sliding Window", func() {
		It("can use Limiter", func() {
			testLimiter("sliding")
		})
	})
	Context("Default Window", func() {
		It("can use Limiter", func() {
			testLimiter("")
		})
	})
})

package middleware_test

import (
	"time"

	"github.com/gofiber/fiber/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3/middleware/cache"
)

type CachePathFilter struct {
	FiberPathFilter
	Cache
}

var _ CacheNext = &CachePathFilter{}

var _ CacheInvalidator = &CachePathFilter{}

var _ CacheKeyGenerator = &CachePathFilter{}

var _ CacheExpirationGenerator = &CachePathFilter{}

var _ CacheConfig = &CachePathFilter{}

// CacheConfig implements [middleware.CacheNext].
func (this *CachePathFilter) CacheNext() func(fiber.Ctx) bool {
	return nil
}

// CacheConfig implements [middleware.CacheInvalidator].
func (this *CachePathFilter) CacheInvalidator() func(fiber.Ctx) bool {
	return nil
}

// CacheConfig implements [middleware.CacheKeyGenerator].
func (this *CachePathFilter) CacheKeyGenerator() func(fiber.Ctx) string {
	return cache.ConfigDefault.KeyGenerator
}

// CacheConfig implements [middleware.CacheExpirationGenerator].
func (this *CachePathFilter) CacheExpirationGenerator() func(fiber.Ctx, *cache.Config) time.Duration {
	return nil
}

// CacheConfig implements [middleware.CacheConfig].
func (this *CachePathFilter) CacheConfig(_ *cache.Config) {}

var _ = Describe("Cache", func() {
	config := map[string]any{
		"storage": "memory",
	}
	It("can use BasicAuth", func() {
		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.cache")).NotTo(HaveOccurred())
		AddBeamFor(&CachePathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

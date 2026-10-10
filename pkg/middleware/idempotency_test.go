package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
)

type IdempotencyPathFilter struct {
	FiberPathFilter
	Idempotency
}

var _ IdempotencyNext = &IdempotencyPathFilter{}

var _ IdempotencyKeyHeaderValidate = &IdempotencyPathFilter{}

var _ IdempotencyLocker = &IdempotencyPathFilter{}

var _ IdempotencyConfig = &IdempotencyPathFilter{}

// IdempotencyNext implements [middleware.IdempotencyNext].
func (i *IdempotencyPathFilter) IdempotencyNext() func(fiber.Ctx) bool {
	return nil
}

// IdempotencyKeyHeaderValidate implements [middleware.IdempotencyKeyHeaderValidate].
func (i *IdempotencyPathFilter) IdempotencyKeyHeaderValidate() func(string) error {
	return idempotency.ConfigDefault.KeyHeaderValidate
}

// IdempotencyLocker implements [middleware.IdempotencyLocker].
func (i *IdempotencyPathFilter) IdempotencyLocker() idempotency.Locker {
	return nil
}

// IdempotencyConfig implements [middleware.IdempotencyConfig].
func (i *IdempotencyPathFilter) IdempotencyConfig(_ *idempotency.Config) {}

var _ = Describe("Idempotency", func() {
	It("can use Idempotency", func() {
		cfg := map[string]any{
			"storage": "memory",
			"keepResponseHeaders": []any{
				"X-Follow",
			},
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.idempotency")).NotTo(HaveOccurred())
		AddBeamFor(&IdempotencyPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

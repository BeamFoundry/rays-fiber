package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/sse"
)

type SSEPathFilter struct {
	FiberPathFilter
	RequestMapping `Method:"GET" Path:"/events"`
	SSE
}

// SSEHandler implements [middleware.SSEHandler].
func (this *SSEPathFilter) SSEHandler() func(ctx fiber.Ctx, stream *sse.Stream) error {
	return func(_ fiber.Ctx, _ *sse.Stream) error {
		return nil
	}
}

// SSEOnClose implements [middleware.SSEOnClose].
func (this *SSEPathFilter) SSEOnClose() func(ctx fiber.Ctx, err error) {
	return nil
}

var _ SSEConfig = &SSEPathFilter{}

var _ SSEHandler = &SSEPathFilter{}

var _ SSEOnClose = &SSEPathFilter{}

// SSEConfig() implements [middleware.SSEConfig].
func (this *SSEPathFilter) SSEConfig(_ *sse.Config) {}

var _ = Describe("SSE", func() {
	It("can use SSE", func() {
		AddBeamFor(&SSEPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

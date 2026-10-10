package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/testing"
)

type CompressPathFilter struct {
	FiberPathFilter
	Compress
}

var _ CompressNext = &CompressPathFilter{}

var _ CompressConfig = &CompressPathFilter{}

// CompressNext() implements [middleware.CompressNext].
func (this *CompressPathFilter) CompressNext() func(fiber.Ctx) bool {
	return nil
}

// CompressConfig() implements [middleware.CompressConfig].
func (this *CompressPathFilter) CompressConfig(_ *compress.Config) {}

var _ = Describe("Compress", func() {
	testLevel := func(level string) {
		config := map[string]any{
			"level": level,
		}
		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.compress")).NotTo(HaveOccurred())
		AddBeamFor(&CompressPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		Expect(GetContainer().Initialize()).NotTo(HaveOccurred())
		appProvider := GetBeamDefaultInstance[FiberAppProvider]()
		app := appProvider.Default()

		req := httptest.NewRequest(http.MethodGet, "/testing/hello", http.NoBody)
		req.SetBasicAuth("admin", "nimda")
		resp, err := app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
	}

	Context("disabled level", func() {
		It("can use Compress", func() {
			testLevel("disabled")
		})
	})

	Context("default level", func() {
		It("can use Compress", func() {
			testLevel("default")
		})
	})

	Context("speed level", func() {
		It("can use Compress", func() {
			testLevel("speed")
		})
	})

	Context("default level", func() {
		It("can use Compress", func() {
			testLevel("compression")
		})
	})
})

package middleware_test

import (
	////"net/http/httptest"

	"io/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type StaticPathFilter struct {
	FiberPathFilter
	Static
}

var _ StaticFS = &StaticPathFilter{}

var _ StaticNext = &StaticPathFilter{}

var _ StaticModifyResponse = &StaticPathFilter{}

var _ StaticNotFoundHandler = &StaticPathFilter{}

var _ StaticConfig = &StaticPathFilter{}

// StaticFS implements [middleware.StaticFS].
func (s *StaticPathFilter) StaticFS() fs.FS {
	return nil
}

// StaticNext implements [middleware.StaticNext].
func (s *StaticPathFilter) StaticNext() func(fiber.Ctx) bool {
	return nil
}

// StaticModifyResponse implements [middleware.StaticModifyResponse].
func (s *StaticPathFilter) StaticModifyResponse() fiber.Handler {
	return nil
}

// StaticNotFoundHandler implements [middleware.StaticNotFoundHandler].
func (s *StaticPathFilter) StaticNotFoundHandler() fiber.Handler {
	return nil
}

// StaticConfig implements [middleware.StaticConfig].
func (s *StaticPathFilter) StaticConfig(_ *static.Config) {}

var _ = Describe("Static", func() {
	Context("No tags", func() {
		It("can use Static", func() {
			AddBeamFor(&StaticPathFilter{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("With tags", func() {
		type StaticPathFilterWithTags struct {
			FiberPathFilter
			Static `Dir:"./testdata"`
		}
		It("can use Static", func() {
			AddBeamFor(&StaticPathFilterWithTags{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

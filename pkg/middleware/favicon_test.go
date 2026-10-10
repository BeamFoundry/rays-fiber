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
	"github.com/gofiber/fiber/v3/middleware/favicon"
)

type FaviconPathFilter struct {
	FiberPathFilter
	Favicon
}

var _ FaviconNext = &FaviconPathFilter{}

var _ FaviconFileSystem = &FaviconPathFilter{}

var _ FaviconConfig = &FaviconPathFilter{}

// FaviconNext implements [middleware.FaviconNext].
func (f *FaviconPathFilter) FaviconNext() func(fiber.Ctx) bool {
	return nil
}

// FaviconFileSystem implements [middleware.FaviconFileSystem].
func (f *FaviconPathFilter) FaviconFileSystem() fs.FS {
	return nil
}

// FaviconConfig implements [middleware.FaviconConfig].
func (f *FaviconPathFilter) FaviconConfig(_ *favicon.Config) {}

var _ = Describe("Favicon", func() {
	It("can use Favicon", func() {
		AddBeamFor(&FaviconPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

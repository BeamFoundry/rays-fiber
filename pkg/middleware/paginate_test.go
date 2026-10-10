package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

type PaginatePathFilter struct {
	FiberPathFilter
	Paginate
}

var _ PaginateNext = &PaginatePathFilter{}

var _ PaginateConfig = &PaginatePathFilter{}

// PaginateNext() implements [middleware.PaginateNext]
func (this *PaginatePathFilter) PaginateNext() func(fiber.Ctx) bool {
	return nil
}

// PaginateConfig() implements [middleware.PaginateConfig].
func (this *PaginatePathFilter) PaginateConfig(_ *paginate.Config) {}

var _ = Describe("Paginate", func() {
	It("can use Paginate", func() {
		AddBeamFor(&PaginatePathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

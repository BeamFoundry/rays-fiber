package middleware_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/template"
)

type TemplatePathFilter struct {
	FiberPathFilter
	Template
}

var _ TemplateNext = &TemplatePathFilter{}

var _ TemplateConfig = &TemplatePathFilter{}

// TemplateNext() implements [middleware.TemplateNext]
func (this *TemplatePathFilter) TemplateNext() func(fiber.Ctx) bool {
	return nil
}

// TemplateConfig() implements [middleware.TemplateConfig].
func (this *TemplatePathFilter) TemplateConfig(_ *template.Config) {}

var _ = Describe("Template", func() {
	config := map[string]any{
		"storage": "memory",
	}
	It("can use Template", func() {
		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.template")).NotTo(HaveOccurred())
		AddBeamFor(&TemplatePathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

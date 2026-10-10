package middleware_test

import (
	"testing"

	"github.com/gofiber/fiber/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/web"
)

func TestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Middleware Suite")
}

type MiddlewareTestController struct {
	Controller
	RequestMapping `Path:"/testing"`
}

func (this *MiddlewareTestController) Hello(
	_ struct {
		RequestMapping `Path:"/hello" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString("hello")
}

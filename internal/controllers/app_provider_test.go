package controllers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	_ "github.com/BeamFoundry/rays-web-security"
	. "github.com/BeamFoundry/rays-web-security/pkg/config"
)

type TestApiSecurityConfiguration struct {
	Configuration
}

func (this *TestApiSecurityConfiguration) ApiHttpSecurity(
	_ struct {
		Server `Name:"api"`
	},
) *HttpSecurity {
	return (&HttpSecurityBuilder{}).
		AnyRequest().DenyAll().
		Basic().
		Build()
}

type TestApiController struct {
	Controller
	Server `Name:"api"`
}

func (this *TestApiController) Ping(
	_ struct {
		RequestMapping `Path:"/ping" Method:"GET"`
	},
	ctx fiber.Ctx,
) error {
	return ctx.SendString("pong")
}

var _ = Describe("DefaultAppProvider", func() {
	It("initializes SecurityFilterChain for a server with HttpSecurity and no FiberPathFilter", func() {
		//nolint:errcheck // Not tested
		GetContainer().LoadConfigAt(map[string]any{}, "rays.fiber.servers.api")

		AddBeamFor(&TestApiSecurityConfiguration{})
		AddBeamFor(&TestApiController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		appProvider := GetBeamScopedInstance[FiberAppProvider]()
		apiApp := appProvider.Provide("api")
		Expect(apiApp).NotTo(BeNil())

		req := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
		resp, testErr := apiApp.Test(req)
		Expect(testErr).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
	})
})

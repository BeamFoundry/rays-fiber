package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
)

type BasicAuthPathFilter struct {
	FiberPathFilter
	Server `_:"primary"`
	BasicAuth
}

var _ BasicAuthNext = &BasicAuthPathFilter{}

var _ BasicAuthAuthorizer = &BasicAuthPathFilter{}

var _ BasicAuthUnauthorized = &BasicAuthPathFilter{}

var _ BasicAuthBadRequest = &BasicAuthPathFilter{}

var _ BasicAuthConfig = &BasicAuthPathFilter{}

// BasicAuthNext() implements [middleware.BasicAuthNext]
func (this *BasicAuthPathFilter) BasicAuthNext() func(fiber.Ctx) bool {
	return nil
}

// BasicAuthAuthorizer() implements [middleware.BasicAuthAuthorizer]
func (this *BasicAuthPathFilter) BasicAuthAuthorizer() func(string, string, fiber.Ctx) bool {
	return nil
}

// BasicAuthUnauthorized() implements [middlewarE.BasicAuthUnauthorized]
func (this *BasicAuthPathFilter) BasicAuthUnauthorized() fiber.Handler {
	return nil
}

// BasicAuthBadRequest() implements [middleware.BasicAuthBadRequest].
func (this *BasicAuthPathFilter) BasicAuthBadRequest() fiber.Handler {
	return nil
}

// BasicAuthConfig() implements [middleware.BasicAuthConfig].
func (this *BasicAuthPathFilter) BasicAuthConfig(_ *basicauth.Config) {}

var _ = Describe("BasicAuth", func() {
	config := map[string]any{
		"users": map[string]any{
			"admin": "$2a$10$iNyufxdpk0vZzxXB/2fumudxnXbV/BO/lZ9sQaOosQe3UKrqPFdnS",
		},
	}
	It("can use BasicAuth", func() {
		// Test for the above issue
		Expect(reflect.TypeOf(config["users"])).To(Equal(reflect.TypeFor[map[string]any]()))

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.basicauth")).NotTo(HaveOccurred())
		AddBeamFor(&BasicAuthPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		Expect(GetContainer().Initialize()).NotTo(HaveOccurred())

		pass, ok := GetContainer().ConfigValue("rays.fiber.servers.primary.basicauth.users.admin")
		Expect(ok).To(BeTrue())
		Expect(pass).NotTo(BeEmpty())

		appProvider := GetBeamScopedInstance[FiberAppProvider]()
		app := appProvider.Default()

		req := httptest.NewRequest(http.MethodGet, "/testing/hello", http.NoBody)
		req.SetBasicAuth("admin", "nimda")
		resp, err := app.Test(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
	})
})

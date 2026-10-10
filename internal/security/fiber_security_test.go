package security_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	. "github.com/BeamFoundry/rays-fiber/internal/security"
)

var _ = Describe("FiberSecurity", func() {
	Describe("CookieSameSite", func() {
		It("can be converted to string", func() {
			Expect(FiberSameSite{}.Convert(http.SameSiteDefaultMode)).To(Equal("disabled"))
			Expect(FiberSameSite{}.Convert(http.SameSiteLaxMode)).To(Equal("Lax"))
			Expect(FiberSameSite{}.Convert(http.SameSiteStrictMode)).To(Equal("Strict"))
			Expect(FiberSameSite{}.Convert(http.SameSiteNoneMode)).To(Equal("None"))
		})
	})
	Describe("FiberSecurityRouter", func() {
		It("can Group", func() {
			app := fiber.New()
			router := FiberSecurityRouter{}.New(app)
			group := router.Group("/test")
			Expect(group).NotTo(BeNil())
		})
	})
	Describe("FiberSecurityHandler", func() {
		It("can HandleSecurityFilterChainException", func() {
			app := fiber.New()
			router := FiberSecurityRouter{}.New(app)
			router.Get("/test", func(handler SecurityHandler) error {
				err := DefaultSecurityFilterChainException{}.ReportError(NewError("IllegalArgumentError", "Invalid parameter"), 400)
				return err.Handle(handler)
			})
			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})
		It("can write", func() {
			app := fiber.New()
			router := FiberSecurityRouter{}.New(app)
			router.Get("/test", func(handler SecurityHandler) error {
				//nolint:errcheck // Testing
				handler.Write([]byte("hello"))
				return nil
			})
			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("hello"))
		})
	})
})

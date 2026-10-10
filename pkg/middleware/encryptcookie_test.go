package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
)

type EncryptCookiePathFilter struct {
	FiberPathFilter
	EncryptCookie
}

var _ EncryptCookieNext = &EncryptCookiePathFilter{}

var _ EncryptCookieEncryptor = &EncryptCookiePathFilter{}

var _ EncryptCookieDecryptor = &EncryptCookiePathFilter{}

var _ EncryptCookieConfig = &EncryptCookiePathFilter{}

// EncryptCookieNext implements [middleware.EncryptCookieNext].
func (e *EncryptCookiePathFilter) EncryptCookieNext() func(fiber.Ctx) bool {
	return nil
}

// EncryptCookieEncryptor implements [middleware.EncryptCookieEncryptor].
func (e *EncryptCookiePathFilter) EncryptCookieEncryptor() func(_, _, _ string) (string, error) {
	return nil
}

// EncryptCookieDecryptor implements [middleware.EncryptCookieDecryptor].
func (e *EncryptCookiePathFilter) EncryptCookieDecryptor() func(_, _, _ string) (string, error) {
	return nil
}

// EncryptCookieConfig implements [middleware.EncryptCookieConfig].
func (e *EncryptCookiePathFilter) EncryptCookieConfig(_ *encryptcookie.Config) {}

var _ = Describe("EncryptCookie", func() {
	Context("key is empty", func() {
		It("can use EncryptCookie", func() {
			AddBeamFor(&EncryptCookiePathFilter{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("key is not empty", func() {
		config := map[string]any{
			"key": encryptcookie.GenerateKey(32),
		}
		It("can use EncryptCookie", func() {
			Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.encryptcookie")).NotTo(HaveOccurred())
			AddBeamFor(&EncryptCookiePathFilter{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

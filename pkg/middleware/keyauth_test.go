package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

type KeyAuthPathFilter struct {
	FiberPathFilter
	KeyAuth
}

var _ KeyAuthNext = &KeyAuthPathFilter{}

var _ KeyAuthSuccessHandler = &KeyAuthPathFilter{}

var _ KeyAuthErrorHandler = &KeyAuthPathFilter{}

var _ KeyAuthValidator = &KeyAuthPathFilter{}

var _ KeyAuthExtractor = &KeyAuthPathFilter{}

var _ KeyAuthConfig = &KeyAuthPathFilter{}

// KeyAuthNext implements [middleware.KeyAuthNext].
func (k *KeyAuthPathFilter) KeyAuthNext() func(fiber.Ctx) bool {
	return nil
}

// KeyAuthSuccessHandler implements [middleware.KeyAuthSuccessHandler].
func (k *KeyAuthPathFilter) KeyAuthSuccessHandler() fiber.Handler {
	return nil
}

// KeyAuthErrorHandler implements [middleware.KeyAuthErrorHandler].
func (k *KeyAuthPathFilter) KeyAuthErrorHandler() fiber.ErrorHandler {
	return nil
}

// KeyAuthValidator implements [middleware.KeyAuthValidator].
func (k *KeyAuthPathFilter) KeyAuthValidator() func(fiber.Ctx, string) (bool, error) {
	return func(fiber.Ctx, string) (bool, error) { return false, nil }
}

// KeyAuthExtractor implements [middleware.KeyAuthExtractor].
func (k *KeyAuthPathFilter) KeyAuthExtractor() extractors.Extractor {
	return extractors.Extractor{}
}

// KeyAuthConfig implements [middleware.KeyAuthConfig].
func (k *KeyAuthPathFilter) KeyAuthConfig(_ *keyauth.Config) {}

var _ = Describe("KeyAuth", func() {
	It("can use KeyAuth", func() {
		AddBeamFor(&KeyAuthPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

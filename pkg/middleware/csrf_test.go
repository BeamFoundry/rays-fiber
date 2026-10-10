package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v3"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type CsrfPathFilter struct {
	FiberPathFilter
	Session
	Csrf
	session bool
}

var _ CsrfNext = &CsrfPathFilter{}

var _ CsrfExtractor = &CsrfPathFilter{}

var _ CsrfKeyGenerator = &CsrfPathFilter{}

var _ CsrfErrorHandler = &CsrfPathFilter{}

var _ CsrfConfig = &CsrfPathFilter{}

var _ SessionStore = &CsrfPathFilter{}

// CsrfNext implements [middleware.CsrfNext].
func (this *CsrfPathFilter) CsrfNext() func(fiber.Ctx) bool {
	return nil
}

// CsrfExtractor implements [middleware.CsrfExtractor].
func (this *CsrfPathFilter) CsrfExtractor() extractors.Extractor {
	return extractors.Extractor{}
}

// CsrfKeyGenerator implements [middleware.CsrfKeyGenerator].
func (this *CsrfPathFilter) CsrfKeyGenerator() func() string {
	return nil
}

// CsrfErrorHandler implements [middleware.CsrfErrorHandler].
func (this *CsrfPathFilter) CsrfErrorHandler() func(fiber.Ctx, error) error {
	return nil
}

// CsrfConfig implements [middleware.CsrfConfig].
func (this *CsrfPathFilter) CsrfConfig(_ *csrf.Config) {}

// SessionStore implements [middleware.SessionStore].
func (this *CsrfPathFilter) SessionStore() *session.Store {
	if this.session {
		return this.Session.SessionStore()
	}
	return nil
}

var _ = Describe("Csrf", func() {
	config := map[string]any{
		"storage": "memory",
	}
	Context("With a SessionStore", func() {
		It("can use Csrf", func() {
			Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.csrf")).NotTo(HaveOccurred())
			AddBeamFor(&CsrfPathFilter{session: true})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("With nil SessionStore", func() {
		It("can not work", func() {
			Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.csrf")).NotTo(HaveOccurred())
			AddBeamFor(&CsrfPathFilter{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).To(HaveOccurred())
		})
	})
})

package middleware_test

import (
	////"net/http/httptest"

	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type SessionPathFilter struct {
	FiberPathFilter
	Session
}

var _ SessionNext = &SessionPathFilter{}

var _ SessionErrorHandler = &SessionPathFilter{}

var _ SessionKeyGenerator = &SessionPathFilter{}

var _ SessionExtractor = &SessionPathFilter{}

var _ SessionConfig = &SessionPathFilter{}

// SessionNext implements [middleware.SessionNext].
func (this *SessionPathFilter) SessionNext() func(fiber.Ctx) bool {
	return nil
}

// SessionErrorHandler implements [middleware.SessionErrorHandler].
func (this *SessionPathFilter) SessionErrorHandler() func(fiber.Ctx, error) {
	return nil
}

// SessionKeyGenerator implements [middleware.SessionKeyGenerator].
func (this *SessionPathFilter) SessionKeyGenerator() func() string {
	return nil
}

// SessionExtractor implements [middleware.SessionExtractor].
func (this *SessionPathFilter) SessionExtractor() extractors.Extractor {
	return extractors.Extractor{}
}

// SessionConfig implements [middleware.SessionConfig].
func (this *SessionPathFilter) SessionConfig(_ *session.Config) {}

func (this *SessionPathFilter) PostConstruct() Error {
	Expect(this.SessionStore()).NotTo(BeNil())
	return nil
}

var _ = Describe("Session", func() {
	It("can use Session", func() {
		cfg := map[string]any{
			"storage": "memory",
		}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.session")).NotTo(HaveOccurred())
		AddBeamFor(&SessionPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

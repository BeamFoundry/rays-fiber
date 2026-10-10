package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type RequestIDPathFilter struct {
	FiberPathFilter
	RequestID
}

var _ RequestIDNext = &RequestIDPathFilter{}

var _ RequestIDGenerator = &RequestIDPathFilter{}

var _ RequestIDConfig = &RequestIDPathFilter{}

// RequestIDNext implements [middleware.RequestIDNext].
func (r *RequestIDPathFilter) RequestIDNext() func(fiber.Ctx) bool {
	return nil
}

// RequestIDGenerator implements [middleware.RequestIDGenerator].
func (r *RequestIDPathFilter) RequestIDGenerator() func() string {
	return nil
}

// RequestIDConfig implements [middleware.RequestIDConfig].
func (r *RequestIDPathFilter) RequestIDConfig(_ *requestid.Config) {}

var _ = Describe("RequestID", func() {
	It("can use RequestID", func() {
		AddBeamFor(&RequestIDPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
})

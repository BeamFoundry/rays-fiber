package security

import (
	"crypto/tls"
	"net/http"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("FiberSecurityHandler TLS", func() {
	It("returns TLSConnectionState correctly", func() {
		// 1. When req has TLS
		tlsState := &tls.ConnectionState{ServerName: "example.com"}
		handler := &FiberSecurityHandler{
			req: &http.Request{TLS: tlsState},
		}
		Expect(handler.TLSConnectionState()).To(Equal(tlsState))

		// 2. When req has no TLS but ctx.RequestCtx() exists
		app := fiber.New()
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		handler2 := &FiberSecurityHandler{
			ctx: ctx,
		}
		Expect(handler2.TLSConnectionState()).To(BeNil())

		// 3. When req is nil and ctx is zero-value
		handler3 := &FiberSecurityHandler{}
		Expect(handler3.TLSConnectionState()).To(BeNil())
	})
})

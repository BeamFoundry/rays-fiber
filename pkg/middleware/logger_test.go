package middleware_test

import (
	////"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type LoggerPathFilter struct {
	FiberPathFilter
	Logger
}

var _ LoggerNext = &LoggerPathFilter{}

var _ LoggerSkip = &LoggerPathFilter{}

var _ LoggerDone = &LoggerPathFilter{}

var _ LoggerCustomTags = &LoggerPathFilter{}

var _ LoggerBeforeHandlerFunc = &LoggerPathFilter{}

var _ LoggerLoggerFunc = &LoggerPathFilter{}

var _ LoggerConfig = &LoggerPathFilter{}

// LoggerNext implements [middleware.LoggerNext].
func (f *LoggerPathFilter) LoggerNext() func(fiber.Ctx) bool {
	return nil
}

// LoggerSkip implements [middleware.LoggerSkip].
func (f *LoggerPathFilter) LoggerSkip() func(fiber.Ctx) bool {
	return nil
}

// LoggerDone implements [middleware.LoggerDone].
func (f *LoggerPathFilter) LoggerDone() func(c fiber.Ctx, logString []byte) {
	return nil
}

// LoggerCustomTags implements [middleware.LoggerCustomTags].
func (f *LoggerPathFilter) LoggerCustomTags() map[string]logger.LogFunc {
	return nil
}

// LoggerBeforeHandlerFunc implements [middleware.LoggerBeforeHandlerFunc].
func (f *LoggerPathFilter) LoggerBeforeHandlerFunc() func(*logger.Config) {
	return nil
}

// LoggerLoggerFunc implements [middleware.LoggerLoggerFunc].
func (f *LoggerPathFilter) LoggerLoggerFunc() func(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error {
	return nil
}

// LoggerConfig implements [middleware.LoggerConfig].
func (f *LoggerPathFilter) LoggerConfig(_ *logger.Config) {}

var _ = Describe("FiberLogger", func() {
	It("can use FiberLogger", func() {
		AddBeamFor(&LoggerPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
	It("can use a custom format", func() {
		cfg := map[string]any{"format": "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n"}
		Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.logger")).NotTo(HaveOccurred())
		AddBeamFor(&LoggerPathFilter{})
		AddBeamFor(&MiddlewareTestController{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())
	})
	It("can use Predefined format", func() {
		formats := []string{"DefaultFormat", "CommonFormat", "CombinedFormat", "JSONFormat", "ECSFormat"}

		for _, format := range formats {
			cfg := map[string]any{"format": format}
			Expect(GetContainer().LoadConfigAt(cfg, "rays.fiber.servers.primary.logger")).NotTo(HaveOccurred())
			AddBeamFor(&LoggerPathFilter{})
			AddBeamFor(&MiddlewareTestController{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
			GetContainer().Reset()
		}
	})
})

package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/sse"
)

// SSEHandler is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter]
// must implements, if it uses [SSE], to indicates the handler that will writes events to the stream.
//
// Please see the [sse] middleware configuration documentation.
//
// [sse]: https://docs.gofiber.io/middleware/sse/#config
type SSEHandler interface {
	SSEHandler() func(ctx fiber.Ctx, stream *sse.Stream) error
}

// SSEOnClose is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter]
// must implements, if it uses [SSE], to indicates the action that will be called when the stream ends,
// with nil when the handler returned successfully and no stream write failed.
//
// Please see the [sse] middleware configuration documentation.
//
// [sse]: https://docs.gofiber.io/middleware/sse/#config
type SSEOnClose interface {
	SSEOnClose() func(ctx fiber.Ctx, err error)
}

// SSEConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [SSE], to review the [sse.Config] that will be configured by [sse]
//
// Please see the [sse] middleware configuration documentation.
//
// [sse]: https://docs.gofiber.io/middleware/sse/#config
type SSEConfig interface {
	SSEConfig(config *sse.Config)
}

type sseConfigurationProperties struct {
	retry             time.Duration
	heartbeatInterval time.Duration
	disableHeartbeat  bool
}

type sseConfig struct {
	sse sseConfigurationProperties
}

var _, _ = any(&FiberMwConfig[sseConfig]{}).(Stereotype)

func (this *sseConfig) Config() sse.Config {
	config := sse.ConfigDefault
	config.Retry = this.sse.retry
	config.HeartbeatInterval = this.sse.heartbeatInterval
	config.DisableHeartbeat = this.sse.disableHeartbeat
	return config
}

// [SSE] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] with a
// [github.com/BeamFOundry/rays/pkg/web.RequestMapping] with a "Path" and "Method"
// struct tag.
//
// # Example
//
//	type SSEEventHandler struct {
//		FiberPathFilter
//		RequestMapping `Path:"/api"`
//		SSE
//	}
//
//	func (this *SSEEventHandler) SSEHandler() func(c fiber.Ctx, stream *sse.Stream) error {
//		// handling stream
//	}
//
//	func (this *SSEEventHandler) SSEConfig(config *sse.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [SSENext], [], [] and [].
//
// Also, the parent container may implements [SSEConfig] to receive a pointer to the [sse.Config]
// that will be used.
//
// Please see those interface definitions and the [sse] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//	              sse:
//	                 retry: 0
//	                 heartbeatInterval: 15s
//	                 disableHeartbeat: false
//
// [sse]: https://docs.gofiber.io/middleware/sse/#config
type SSE struct {
	config *FiberMwConfig[sseConfig] `@:"Inject"`
	cfg    *sse.Config
}

func (this *SSE) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if handler, ok := (parent).(SSEHandler); ok {
		this.cfg.Handler = handler.SSEHandler()
	}
	if onClose, ok := (parent).(SSEOnClose); ok {
		this.cfg.OnClose = onClose.SSEOnClose()
	}
	if config, ok := (parent).(SSEConfig); ok {
		config.SSEConfig(this.cfg)
	}
	return nil
}

func (this *SSE) Handler() fiber.Handler {
	return sse.New(*this.cfg)
}

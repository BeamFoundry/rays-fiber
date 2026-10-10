package middleware

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// LoggerNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates when to skip it this middleware.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerNext interface {
	LoggerNext() func(fiber.Ctx) bool
}

// LoggerSkip is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates a function to determine if logging is skipped or written to Stream.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerSkip interface {
	LoggerSkip() func(fiber.Ctx) bool
}

// LoggerDone is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates a function that is called after the log string for a request is written
// to Stream, and pass the log string as parameter.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerDone interface {
	LoggerDone() func(c fiber.Ctx, logString []byte)
}

// LoggerCustomTags is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates a custom tag actions for this logger instance. These tags override
// built-in and globally registered tags with the same name.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerCustomTags interface {
	LoggerCustomTags() map[string]logger.LogFunc
}

// LoggerBeforeHandlerFunc is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates a BeforeHandlerFunc runs once before the handler is built,
// letting you customize colors, template, etc.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerBeforeHandlerFunc interface {
	LoggerBeforeHandlerFunc() func(*logger.Config)
}

// LoggerLoggerFunc is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to indicates a custom logger function for integration with logging
// libraries (Zerolog, Zap, Logrus, etc). Defaults to Fiber's default logger if not defined.
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerLoggerFunc interface {
	LoggerLoggerFunc() func(c fiber.Ctx, data *logger.Data, cfg *logger.Config) error
}

// LoggerConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Logger], to review the [logger.Config] that will be configured by [Logger]
//
// Please see the [logger] middleware configuration documentation.
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type LoggerConfig interface {
	LoggerConfig(config *logger.Config)
}

// TODO: Add Steam (Stdout/Stderr/ (file:// ?) ) and TimeDone
type loggerConfigurationProperties struct {
	format        string
	timeFormat    string
	timeZone      string
	timeInterval  time.Duration
	disableColors bool
	forceColors   bool
}

type loggerConfig struct {
	logger loggerConfigurationProperties
}

var _, _ = any(&FiberMwConfig[loggerConfig]{}).(Stereotype)

func (this *loggerConfig) Config() logger.Config {
	config := logger.ConfigDefault

	switch this.logger.format {
	case "DefaultFormat":
		config.Format = "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n"
	case "CommonFormat":
		config.Format = "${ip} - - [${time}] \"${method} ${url} ${protocol}\" ${status} ${bytesSent}\n"
	case "CombinedFormat":
		config.Format = "${ip} - - [${time}] \"${method} ${url} ${protocol}\" ${status} ${bytesSent} \"${referer}\" \"${ua}\"\n"
	case "JSONFormat":
		config.Format = "{time: ${time}, ip: ${ip}, method: ${method}, url: ${url}, status: ${status}, bytesSent: ${bytesSent}}\n"
	case "ECSFormat":
		config.Format = "{\"@timestamp\":\"${time}\",\"ecs\":{\"version\":\"1.6.0\"},\"client\":{\"ip\":\"${ip}\"}," +
			"\"http\":{\"request\":{\"method\":\"${method}\",\"url\":\"${url}\",\"protocol\":\"${protocol}\"},\"response\":{\"status_code\":${status}," +
			"\"body\":{\"bytes\":${bytesSent}}}},\"log\":{\"level\":\"INFO\",\"logger\":\"fiber\"},\"message\":\"${method} ${url} responded with ${status}\"}\n"
	default:
		config.Format = this.logger.format
	}
	config.TimeFormat = this.logger.timeFormat
	config.TimeZone = this.logger.timeZone
	config.TimeInterval = this.logger.timeInterval
	config.DisableColors = this.logger.disableColors
	config.ForceColors = this.logger.forceColors

	return config
}

// [Logger] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  Logger
//	}
//
//	func (this *ApiController) FiberLoggerConfig(config *logger.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [LoggerNext], [LoggerSkip], [LoggerDone], [LoggerCustomTags],
// [LoggerBeforeHandlerFunc] and [LoggerLoggerFunc].
//
// Also, the parent container may implements [LoggerConfig] to receive a pointer to the [logger.Config]
// that will be used.
//
// Please see those interface definitions and the [logger] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					logger:
//						format: DefaultFormat
//						timeFormat: 15:04:05
//						timeZone: Local
//						timeInterval: 500ms
//						disableColors: false
//						forceColors: false
//
// [logger]: https://docs.gofiber.io/middleware/logger/#config
type Logger struct {
	config *FiberMwConfig[loggerConfig] `@:"Inject"`
	cfg    *logger.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *Logger) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(LoggerNext); ok {
		this.cfg.Next = next.LoggerNext()
	}
	if skip, ok := (parent).(LoggerSkip); ok {
		this.cfg.Skip = skip.LoggerSkip()
	}
	if done, ok := (parent).(LoggerDone); ok {
		this.cfg.Done = done.LoggerDone()
	}
	if customTags, ok := (parent).(LoggerCustomTags); ok {
		this.cfg.CustomTags = customTags.LoggerCustomTags()
	}
	if beforeHandlerFunc, ok := (parent).(LoggerBeforeHandlerFunc); ok {
		this.cfg.BeforeHandlerFunc = beforeHandlerFunc.LoggerBeforeHandlerFunc()
	}
	if loggerFunc, ok := (parent).(LoggerLoggerFunc); ok {
		this.cfg.LoggerFunc = loggerFunc.LoggerLoggerFunc()
	}
	if config, ok := (parent).(LoggerConfig); ok {
		config.LoggerConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *Logger) Handler() fiber.Handler {
	return logger.New(*this.cfg)
}

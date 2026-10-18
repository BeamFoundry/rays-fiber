package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

// PaginateNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Paginate], to indicates when to skip it this middleware.
//
// Please see the [paginate] middleware configuration documentation.
//
// [paginate]: https://docs.gofiber.io/middleware/paginate/#config
type PaginateNext interface {
	PaginateNext() func(fiber.Ctx) bool
}

// PaginateConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [Paginate], to review the [paginate.Config] that will be configured by [paginate]
//
// Please see the [paginate] middleware configuration documentation.
//
// [paginate]: https://docs.gofiber.io/middleware/paginate/#config
type PaginateConfig interface {
	PaginateConfig(config *paginate.Config)
}

type paginateConfigurationProperties struct {
	pageKey      string
	limitKey     string
	sortKey      string
	defaultSort  string
	cursorKey    string
	offsetKey    string
	cursorParam  string
	allowedSorts []string
	defaultPage  int
	defaultLimit int
	maxLimit     int
}

type paginateConfig struct {
	paginate paginateConfigurationProperties
}

var _, _ = any(&FiberMwConfig[paginateConfig]{}).(Stereotype)

func (this *paginateConfig) Config() paginate.Config {
	config := paginate.ConfigDefault

	config.PageKey = this.paginate.pageKey
	config.LimitKey = this.paginate.limitKey
	config.SortKey = this.paginate.sortKey
	config.DefaultSort = this.paginate.defaultSort
	config.CursorKey = this.paginate.cursorKey
	config.OffsetKey = this.paginate.offsetKey
	config.CursorParam = this.paginate.cursorParam
	config.AllowedSorts = this.paginate.allowedSorts
	config.DefaultPage = this.paginate.defaultPage
	config.DefaultLimit = this.paginate.defaultLimit
	config.MaxLimit = this.paginate.maxLimit

	return config
}

// [Paginate] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//		FiberController
//		RequestMapping `Path:"/api"`
//		Paginate
//	}
//
//	func (this *ApiController) PaginateConfig(config *paginate.Config) {
//		// changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [PaginateNext], [], [] and [].
//
// Also, the parent container may implements [PaginateConfig] to receive a pointer to the [paginate.Config]
// that will be used.
//
// Please see those interface definitions and the [paginate] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//	        paginate:
//	          expiration: 5m
//
// [paginate]: https://docs.gofiber.io/middleware/paginate/#config
type Paginate struct {
	config *FiberMwConfig[paginateConfig] `@:"Inject"`
	cfg    *paginate.Config
}

func (this *Paginate) PostConstruct(parent any) Error {
	cfg := this.config.GetConfig(parent)
	this.cfg = new(cfg.Config())
	if next, ok := parent.(PaginateNext); ok {
		this.cfg.Next = next.PaginateNext()
	}
	if config, ok := parent.(PaginateConfig); ok {
		config.PaginateConfig(this.cfg)
	}
	return nil
}

func (this *Paginate) Handler() fiber.Handler {
	return paginate.New(*this.cfg)
}

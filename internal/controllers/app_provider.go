package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"sync"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/web"

	. "github.com/BeamFoundry/rays-fiber/internal/security"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	. "github.com/BeamFoundry/rays-web-security/pkg/config"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	"github.com/gofiber/fiber/v3"
)

type preTxRoute struct {
	method   string
	path     string
	handlers []fiber.Handler
}

type DefaultAppProvider struct {
	Service
	appsConfig          *FiberAppServerConfig        `@:"Inject"`
	filters             map[string][]FiberPathFilter `@:"Inject{GroupBy=Server,Default=primary}"`
	logger              Logging                      `@:"Inject"`
	httpSecurity        map[string][]*HttpSecurity   `@:"Inject{GroupBy=Server,Default=primary}"`
	apps                map[string]*fiber.App
	preTxRoutes         map[string][]preTxRoute
	securityMiddlewares Foundry[*FiberSecurityAppHandler] `@:"Inject"`
	securityFilterChain Optional[SecurityFilterChain]     `@:"Inject"`
	txManager           Optional[TransactionManager]      `@:"Inject"`
	lock                sync.Mutex
}

var _, _ = any(&DefaultAppProvider{}).(Stereotype)

var _ FiberAppProvider = &DefaultAppProvider{}

func (this *DefaultAppProvider) Default() *fiber.App {
	return this.Provide("primary")
}

func (this *DefaultAppProvider) Provide(name string) *fiber.App {
	if app, ok := this.apps[name]; ok {
		return app
	} else if this.appsConfig.Exists(name) {
		this.lock.Lock()
		defer this.lock.Unlock()
		config := this.appsConfig.Config(name)
		app := fiber.New(config)
		this.customize(name, app)
		this.apps[name] = app
		this.logger.Debug("AppProvider adding a new app", "name", name)
		return app
	}
	return nil
}

func (this *DefaultAppProvider) customize(name string, app *fiber.App) {
	if pathFilters, ok := this.filters[name]; ok {
		for _, pathFilter := range pathFilters {
			var path, method string
			mappingMethod := RequestMappingMethod{}.New().GET
			reflectType := reflect.TypeOf(pathFilter)
			if field := FindFieldWithAnnotation[RequestMapping](reflectType.Elem()); field != nil {
				if p, ok := field.Tag.Lookup("Path"); ok && p != "" {
					path = p
					if m, ok := field.Tag.Lookup("Method"); ok && m != "" {
						if err := mappingMethod.From(m); err != nil {
							this.logger.Error("IllegalArgumentException", fmt.Sprintf("Method %s is invalid", m))
							continue
						}
						method = mappingMethod.Name()
					}
				}
			}
			handlers := this.getInstanceHandlers(reflect.ValueOf(pathFilter))
			if _, ok := this.httpSecurity[name]; ok {
				handlers = slices.DeleteFunc(handlers, func(handler FiberHandler) bool {
					var filter *reflect.StructField
					for field := range reflectType.Elem().Fields() {
						if field.Type == reflect.TypeOf(handler).Elem() {
							filter = FindFieldForInterface[SecurityEventFilter](field.Type)
						}
					}
					return filter != nil
				})
			}
			finalHandlers := []fiber.Handler{}
			for _, handler := range handlers {
				finalHandlers = append(finalHandlers, handler.Handler())
			}
			args := []any{}
			if method != "" {
				args = append(args, toAny(finalHandlers)...)
				switch mappingMethod {
				case RequestMappingMethod{}.New().ALL:
					app.All(path, args[0], args[1:]...)
				default:
					app.Add([]string{mappingMethod.String()}, path, args[0], args[1:]...)
				}
			} else {
				if path != "" {
					args = append(args, path)
				}
				args = append(args, toAny(finalHandlers)...)
				app.Use(args...)
			}
		}
	}

	if err := Try(func() {
		if _, ok := this.httpSecurity[name]; ok {
			app.Use(this.preTxDispatchHandler(name))
			if txMgr, ok := this.txManager.Get(); ok {
				app.Use(bindTxContextHandler(txMgr))
			}
			if chain, ok := this.securityFilterChain.Get(); ok {
				ctx := context.Background()
				ctx = context.WithValue(ctx, PrototypeContextProxyType, app)
				ctx = context.WithValue(ctx, PrototypeContextProxyInstance, name)
				mdw, err := this.securityMiddlewares.Forge(ctx)
				this.logger.Debug("Initializing SecurityFilterChain")
				chain.Initialize(name, mdw)
				this.logger.Debug("SecurityFilterChain initialized")
				RaiseNilErr(err)
			}
		}
	}); err != nil {
		err.Log(slog.LevelWarn)
	}
}

func (this *DefaultAppProvider) hasHttpSecurity(name string) bool {
	_, ok := this.httpSecurity[name]
	return ok
}

func (this *DefaultAppProvider) registerPreTxHandlers(name, method, path string, handlers ...fiber.Handler) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.preTxRoutes == nil {
		this.preTxRoutes = map[string][]preTxRoute{}
	}
	this.preTxRoutes[name] = append(this.preTxRoutes[name], preTxRoute{
		method:   method,
		path:     path,
		handlers: handlers,
	})
}

func (this *DefaultAppProvider) preTxDispatchHandler(name string) fiber.Handler {
	return func(c fiber.Ctx) error {
		this.lock.Lock()
		routes := slices.Clone(this.preTxRoutes[name])
		this.lock.Unlock()

		var matched []fiber.Handler
		for _, route := range routes {
			if matchPreTxRoute(route.method, route.path, c.Method(), c.Path()) {
				matched = append(matched, route.handlers...)
			}
		}
		if len(matched) == 0 {
			return c.Next()
		}
		if len(matched) == 1 {
			return matched[0](c)
		}
		return runPreTxChain(c, matched)
	}
}

type preTxCtx struct {
	fiber.Ctx
	next func() error
}

func (this *preTxCtx) Next() error {
	return this.next()
}

func runPreTxChain(c fiber.Ctx, handlers []fiber.Handler) error {
	idx := 0
	var run func() error
	run = func() error {
		if idx >= len(handlers) {
			return c.Next()
		}
		h := handlers[idx]
		idx++
		return h(&preTxCtx{Ctx: c, next: run})
	}
	return run()
}

func matchPreTxRoute(routeMethod, routePath, reqMethod, reqPath string) bool {
	if routeMethod != (RequestMappingMethod{}).New().ALL.String() && routeMethod != reqMethod {
		return false
	}
	return matchRoutePattern(routePath, reqPath)
}

func matchRoutePattern(pattern, reqPath string) bool {
	if pattern == reqPath {
		return true
	}
	pSegs := strings.Split(strings.Trim(pattern, "/"), "/")
	rSegs := strings.Split(strings.Trim(reqPath, "/"), "/")
	for i, p := range pSegs {
		if p == "*" || p == "+" {
			return true
		}
		if i >= len(rSegs) {
			return false
		}
		if strings.HasPrefix(p, ":") {
			continue
		}
		if p != rSegs[i] {
			return false
		}
	}
	return len(pSegs) == len(rSegs)
}

func bindTxContextHandler(txMgr TransactionManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		prev := txMgr.GetContext()
		txMgr.SetContext(c.Context())
		defer txMgr.SetContext(prev)
		return c.Next()
	}
}

func (this *DefaultAppProvider) getInstanceHandlers(instance reflect.Value) []FiberHandler {
	handlers := []FiberHandler{}
	value := instance
	annotedType := instance.Type()
	if annotedType.Kind() == reflect.Pointer {
		value = value.Elem()
		annotedType = annotedType.Elem()
	}
	for i := range annotedType.NumField() {
		field := annotedType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if middleware, ok := value.Field(i).Addr().Interface().(FiberHandler); ok {
				handlers = append(handlers, middleware)
			}
		}
	}
	return handlers
}

func (this *DefaultAppProvider) Run(wg *sync.WaitGroup) {
	for name, app := range this.apps {
		address := this.appsConfig.Address(name)
		listener := this.appsConfig.Listener(name)
		wg.Go(func() {
			err := Try(func() {
				RaiseNilErr(app.Listen(address, listener))
			})
			err.LogError()
		})
	}
}

func (this *DefaultAppProvider) Shutdown() {
	for _, app := range this.apps {
		err := Try(func() {
			RaiseNilErr(app.Shutdown())
		})
		err.LogError()
	}
}

func (this *DefaultAppProvider) PostConstruct() Error {
	this.apps = map[string]*fiber.App{}
	for name := range this.appsConfig.config.servers {
		this.Provide(name)
	}
	return nil
}

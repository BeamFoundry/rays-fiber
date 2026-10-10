package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
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

type DefaultAppProvider struct {
	Service
	appsConfig          *FiberAppServerConfig        `@:"Inject"`
	filters             map[string][]FiberPathFilter `@:"Inject{GroupBy=Server,Default=primary}"`
	logger              Logging                      `@:"Inject"`
	httpSecurity        map[string][]*HttpSecurity   `@:"Inject{GroupBy=Server,Default=primary}"`
	apps                map[string]*fiber.App
	securityMiddlewares Foundry[*FiberSecurityAppHandler] `@:"Inject"`
	securityFilterChain Optional[SecurityFilterChain]     `@:"Inject"`
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

		if err := Try(func() {
			if _, ok := this.httpSecurity[name]; ok {
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

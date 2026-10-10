package controllers

import (
	"net/http"
	"reflect"

	. "github.com/BeamFoundry/rays/pkg/lang"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v3"
)

type fiberMethodFunc struct {
	CANONICAL           fiber.Handler
	CANONICAL_VOID      func(fiber.Ctx)
	EXPR                func(fiber.Req, fiber.Res) error
	EXPR_VOID           func(fiber.Req, fiber.Res)
	EXPR_NT             func(fiber.Req, fiber.Res, func() error) error
	EXPR_NT_VOID        func(fiber.Req, fiber.Res, func() error)
	EXPR_NT_NO          func(fiber.Req, fiber.Res, func()) error
	EXPR_NT_NO_VOID     func(fiber.Req, fiber.Res, func())
	EXPR_NT_ERR         func(fiber.Req, fiber.Res, func(error) error) error
	EXPR_NT_ERR_VOID    func(fiber.Req, fiber.Res, func(error) error)
	EXPR_NT_ERR_NO      func(fiber.Req, fiber.Res, func(error)) error
	EXPR_NT_ERR_NO_VOID func(fiber.Req, fiber.Res, func(error))
	HTTP_HANDLER_FUNC   func() http.HandlerFunc
	HTTP_HANDLER        func() http.Handler
	HTTP_FUNC           func(http.ResponseWriter, *http.Request)
	FASTHTTP            func(*fasthttp.RequestCtx)
	FASTHTTP_ERR        func(*fasthttp.RequestCtx) error
}

type ControllerMethods struct {
	CANONICAL,
	CANONICAL_VOID,
	EXPR,
	EXPR_VOID,
	EXPR_NT,
	EXPR_NT_VOID,
	EXPR_NT_NO,
	EXPR_NT_NO_VOID,
	EXPR_NT_ERR,
	EXPR_NT_ERR_VOID,
	EXPR_NT_ERR_NO,
	EXPR_NT_ERR_NO_VOID,
	HTTP_HANDLER_FUNC,
	HTTP_HANDLER,
	HTTP_FUNC,
	FASTHTTP,
	FASTHTTP_ERR Enum[ControllerMethods]
}

type ControllerMethod struct {
	method *reflect.Method
	Enum[ControllerMethods]
}

func (this *ControllerMethod) Handler(defaults []reflect.Value, instance any) any {
	switch this.Enum {
	case Enum[ControllerMethods]{}.New().CANONICAL:
		return func(ctx fiber.Ctx) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(ctx))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().CANONICAL_VOID:
		return func(ctx fiber.Ctx) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(ctx))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().EXPR:
		return func(req fiber.Req, res fiber.Res) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().EXPR_VOID:
		return func(req fiber.Req, res fiber.Res) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT:
		return func(req fiber.Req, res fiber.Res, next func() error) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_VOID:
		return func(req fiber.Req, res fiber.Res, next func() error) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_NO:
		return func(req fiber.Req, res fiber.Res, next func()) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_NO_VOID:
		return func(req fiber.Req, res fiber.Res, next func()) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_ERR:
		return func(req fiber.Req, res fiber.Res, next func(error) error) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_ERR_VOID:
		return func(req fiber.Req, res fiber.Res, next func(error) error) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_ERR_NO:
		return func(req fiber.Req, res fiber.Res, next func(error)) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	case Enum[ControllerMethods]{}.New().EXPR_NT_ERR_NO_VOID:
		return func(req fiber.Req, res fiber.Res, next func(error)) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(req), reflect.ValueOf(res), reflect.ValueOf(next))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().HTTP_HANDLER,
		Enum[ControllerMethods]{}.New().HTTP_HANDLER_FUNC:

		in := []reflect.Value{}
		in = append(in, defaults...)
		out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		return out[0].Interface()
	case Enum[ControllerMethods]{}.New().HTTP_FUNC:
		return func(w http.ResponseWriter, req *http.Request) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(w), reflect.ValueOf(req))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		}
	case Enum[ControllerMethods]{}.New().FASTHTTP:
		return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(ctx))
			reflect.ValueOf(instance).Method(this.method.Index).Call(in)
		})
	// case Enum[ControllerMethods]{}.New().FASTHTTP_ERR:
	default:
		return func(ctx *fasthttp.RequestCtx) error {
			in := []reflect.Value{}
			in = append(in, defaults...)
			in = append(in, reflect.ValueOf(ctx))
			out := reflect.ValueOf(instance).Method(this.method.Index).Call(in)
			ret, _ := out[0].Interface().(error)
			return ret
		}
	}
}

func (this *ControllerMethod) isAnonymousStruct(m reflect.Type) bool {
	switch {
	case m.Name() != "",
		m.Kind() == reflect.Pointer && m.Elem().Kind() != reflect.Struct,
		m.Kind() == reflect.Pointer && m.Elem().PkgPath() != "":
		return false
	// case m.Kind() == reflect.Struct || m.Kind() == reflect.Pointer:
	default:
		return true
	}
}

func (this ControllerMethod) FromMethod(method *reflect.Method) (*ControllerMethod, bool) {
	m := method.Type
	fmf := reflect.TypeFor[fiberMethodFunc]()

	skip := 1
	//nolint:gocritic // Keeping it
	// GetLogging().Debug("Testing for anonymous method", "name", method.Name)
	if m.NumIn() > skip && this.isAnonymousStruct(m.In(skip)) {
		skip++
	}
	this.method = method

field:
	for field := range fmf.Fields() {
		f := field.Type
		if f.NumOut() != m.NumOut() {
			continue
		}
		if f.NumOut() > 1 || (f.NumOut() == 1 && f.Out(0) != m.Out(0)) {
			continue
		}
		if f.NumIn()+skip != m.NumIn() {
			continue
		}
		for j := 0; j < f.NumIn(); j++ {
			if m.NumIn()-skip <= j || f.In(j) != m.In(j+skip) {
				continue field
			}
		}
		if err := this.From(field.Name); err == nil {
			return &this, true
		}
	}
	return nil, false
}

func (this *ControllerMethod) Method() *reflect.Method {
	return this.method
}

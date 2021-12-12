package fly

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// HTTP methods
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

var (
	methods = [...]string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}
)

type Fly struct {
	HttpServer *http.Server

	// debug mode : print some info
	IsDebug bool

	/**
	{
	 [GET]{
	    "/A/B":ControllerHandler,
	    "/C/D":ControllerHandler,
	 },
	[POST]{
	    "/Y":ControllerHandler,
	    "/Z/X":ControllerHandler,
	 },
	 ...
	*/
	router map[string]map[string]HandlerFunc

	// pre global middlewares
	preMiddleware []MiddlewareFunc

	// global middlewares
	middleware []MiddlewareFunc

	// 错误处理的handler
	ErrorDefaultHandler HandlerFunc
}

func NewFly(isDebug bool) *Fly {
	fly := &Fly{
		router:        make(map[string]map[string]HandlerFunc),
		preMiddleware: make([]MiddlewareFunc, 0),
		middleware:    make([]MiddlewareFunc, 0),
		IsDebug:       isDebug,
	}

	for _, method := range methods {
		fly.router[method] = make(map[string]HandlerFunc)
	}
	return fly
}

// group middleware -> controller
func (fly *Fly) add(method string, path string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	path = strings.ToUpper(path)

	f := func(ctx Context) error {
		h := handlerFunc
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h(ctx)
	}

	fly.router[method][path] = f
}

func (fly *Fly) FindRouter(request *http.Request) HandlerFunc {
	method := request.Method
	path := strings.ToUpper(request.RequestURI)
	if handler, hit := fly.router[method][path]; hit {
		return handler
	}
	return nil
}

func (fly *Fly) Pre(middleware ...MiddlewareFunc) {
	fly.preMiddleware = append(fly.preMiddleware, middleware...)
}

func (fly *Fly) Use(middleware ...MiddlewareFunc) {
	fly.middleware = append(fly.middleware, middleware...)
}

func (fly *Fly) Group(prefix string) *Group {
	return NewGroup(prefix, fly)
}

func (fly *Fly) HEAD(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(HEAD, path, handler, middleware...)
}

func (fly *Fly) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(GET, path, handler, middleware...)
}

func (fly *Fly) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(POST, path, handler, middleware...)
}

func (fly *Fly) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(PUT, path, handler, middleware...)
}

func (fly *Fly) OPTIONS(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(OPTIONS, path, handler, middleware...)
}

func (fly *Fly) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(DELETE, path, handler, middleware...)
}

func (fly *Fly) CONNECT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(CONNECT, path, handler, middleware...)
}

func (fly *Fly) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(PATCH, path, handler, middleware...)
}

func (fly *Fly) TRACE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	fly.add(TRACE, path, handler, middleware...)
}

func (fly *Fly) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range methods {
		fly.add(method, path, handler, middleware...)
	}
}

// global pre middleware -> global middleware -> group middleware -> method middleware -> controller
func (fly *Fly) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// global middleware + controller
	handler := func(ctx Context) error {
		h := fly.FindRouter(request)
		if h == nil {
			return ctx.JSON(http.StatusInternalServerError, "inner err")
		}

		for i := len(fly.middleware) - 1; i >= 0; i-- {
			h = fly.middleware[i](h)
		}

		return h(ctx)
	}

	// pre middleware
	for i := len(fly.preMiddleware) - 1; i >= 0; i-- {
		handler = fly.preMiddleware[i](handler)
	}

	ctx := NewContext(request, response)
	handler(ctx)
}

func (fly *Fly) StartServer(server *http.Server) {
	// debug
	fly.printRouter()

	fly.HttpServer = server
	fly.HttpServer.ListenAndServe()
}

func (fly *Fly) printRouter() {
	if !fly.IsDebug {
		return
	}

	fmt.Println("fly->router:")
	for method, router := range fly.router {
		fmt.Printf("    %s:\n", method)
		for path, handler := range router {
			fmt.Printf("        %s:%s\n", path, reflect.ValueOf(handler).String())
		}
		fmt.Println()
	}
}

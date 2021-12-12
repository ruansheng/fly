package fly

type Group struct {
	prefix     string
	middleware []MiddlewareFunc
	fly        *Fly
}

func NewGroup(prefix string, fly *Fly) *Group {
	return &Group{
		prefix:     prefix,
		middleware: make([]MiddlewareFunc, 0),
		fly:        fly,
	}
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) Group(prefix string) *Group {
	group := g.fly.Group(g.prefix + prefix)

	// extends parent's middleware
	group.middleware = append(group.middleware, g.middleware...)
	return group
}

func (g *Group) add(method string, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	m := g.middleware
	m = append(m, middleware...)
	g.fly.add(method, g.prefix+path, handler, m...)
}

func (g *Group) HEAD(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(HEAD, path, handler, middleware...)
}

func (g *Group) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(GET, path, handler, middleware...)
}

func (g *Group) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(POST, path, handler, middleware...)
}

func (g *Group) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(PUT, path, handler, middleware...)
}

func (g *Group) OPTIONS(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(OPTIONS, path, handler, middleware...)
}

func (g *Group) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(DELETE, path, handler, middleware...)
}

func (g *Group) CONNECT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(CONNECT, path, handler, middleware...)
}

func (g *Group) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(PATCH, path, handler, middleware...)
}

func (g *Group) TRACE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.add(TRACE, path, handler, middleware...)
}

func (g *Group) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range methods {
		g.add(method, path, handler, middleware...)
	}
}

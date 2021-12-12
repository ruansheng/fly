package fly

type MiddlewareFunc func(HandlerFunc) HandlerFunc

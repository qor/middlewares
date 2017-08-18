package middlewares

import "net/http"

// MiddlewareHandler HTTP middleware
type MiddlewareHandler func(http.Handler) http.Handler

// Middlewares middlewares stack
type Middlewares struct {
	middlewares []*Middleware
}

// Use use middleware
func (middlewares *Middlewares) Use(name string, handler MiddlewareHandler) {
	middlewares.middlewares = append(middlewares.middlewares, &Middleware{
		middlewares: middlewares,
		Name:        name,
		Handler:     handler,
	})
}

// Remove remove middleware by name
func (middlewares *Middlewares) Remove(name string) {
	registeredMiddlewares := middlewares.middlewares
	for idx, middleware := range registeredMiddlewares {
		if middleware.Name == name {
			if idx > 0 {
				middlewares.middlewares = middlewares.middlewares[0 : idx-1]
			} else {
				middlewares.middlewares = []*Middleware{}
			}

			if idx < len(registeredMiddlewares)-1 {
				middlewares.middlewares = append(middlewares.middlewares, registeredMiddlewares[idx+1:]...)
			}
		}
	}
}

// Before insert middleware before name
func (middlewares *Middlewares) Before(name ...string) Middleware {
	return Middleware{
		Before:      name,
		middlewares: middlewares,
	}
}

// After insert middleware after name
func (middlewares *Middlewares) After(name ...string) Middleware {
	return Middleware{
		After:       name,
		middlewares: middlewares,
	}
}

func (middlewares *Middlewares) String() string {
	// TODO sort, compile middlewares, print its name in order
	return ""
}

// Apply apply middlewares to handler
func (middlewares *Middlewares) Apply(handler http.Handler) http.Handler {
	// TODO sort, compile middlewares, wrap current handler
	return handler
}

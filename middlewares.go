package middlewares

import (
	"fmt"
	"net/http"
)

// MiddlewareHandler HTTP middleware
type MiddlewareHandler func(http.Handler) http.Handler

// Middlewares middlewares stack
type Middlewares struct {
	middlewares []*Middleware
}

// Use use middleware
func (middlewares *Middlewares) Use(middleware Middleware) {
	middlewares.middlewares = append(middlewares.middlewares, &middleware)
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

func (middlewares *Middlewares) Compile() error {
	var (
		errs           []error
		middlewaresMap = map[string]*Middleware{}
	)

	for _, middleware := range middlewares.middlewares {
		middlewaresMap[middleware.Name] = middleware
	}

	for _, middleware := range middlewaresMap {
		for _, require := range middleware.Requires {
			if _, ok := middlewaresMap[require]; !ok {
				errs = append(errs, fmt.Errorf("middleware %v requires %v, but it doesn't exist", middleware.Name, require))
			}
		}
	}

	return nil
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

package middlewares

// Middleware middleware struct
type Middleware struct {
	Name     string
	Handler  MiddlewareHandler
	Before   []string
	After    []string
	Requires []string

	middlewares *Middlewares
}

// Use use middleware
func (middleware Middleware) Use(name string, handler MiddlewareHandler) {
	middleware.Name = name
	middleware.Handler = handler
	if middleware.middlewares != nil {
		middleware.middlewares.middlewares = append(middleware.middlewares.middlewares, &middleware)
	}
}

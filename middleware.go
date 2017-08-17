package middlewares

// Middleware middleware struct
type Middleware struct {
	Name    string
	Handler MiddlewareHandler

	middlewares *Middlewares
	before      string
	after       string
}

// Use use middleware
func (middleware Middleware) Use(name string, handler MiddlewareHandler) {
	middleware.Name = name
	middleware.Handler = handler
	if middleware.middlewares != nil {
		middleware.middlewares.middlewares = append(middleware.middlewares.middlewares, &middleware)
	}
}

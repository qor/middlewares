package middlewares

// Middleware middleware struct
type Middleware struct {
	Name     string
	Handler  MiddlewareHandler
	Before   []string
	After    []string
	Requires []string
}

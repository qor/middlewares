package middlewares

import "net/http"

// DefaultMiddlewareStack default middleware stack
var DefaultMiddlewareStack = &Middlewares{}

// Use register middleware from DefaultMiddlewareStack
func Use(name string, handler MiddlewareHandler) {
	DefaultMiddlewareStack.Use(name, handler)
}

// Remove remove middleware by name from DefaultMiddlewareStack
func Remove(name string) {
	DefaultMiddlewareStack.Remove(name)
}

// Before insert middleware before name into DefaultMiddlewareStack
func Before(name string) Middleware {
	return DefaultMiddlewareStack.Before(name)
}

// After insert middleware after name into DefaultMiddlewareStack
func After(name string) Middleware {
	return DefaultMiddlewareStack.After(name)
}

// Apply apply middlewares to handler
func Apply(handler http.Handler) http.Handler {
	return DefaultMiddlewareStack.Apply(handler)
}

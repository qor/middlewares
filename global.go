package middlewares

// GlobalMiddlewareStack default middleware stack
var GlobalMiddlewareStack = &Middlewares{}

// Use register middleware from GlobalMiddlewareStack
func Use(name string, handler MiddlewareHandler) {
	GlobalMiddlewareStack.Use(name, handler)
}

// Remove remove middleware by name from GlobalMiddlewareStack
func Remove(name string) {
	GlobalMiddlewareStack.Remove(name)
}

// Before insert middleware before name into GlobalMiddlewareStack
func Before(name string) Middleware {
	return GlobalMiddlewareStack.Before(name)
}

// After insert middleware after name into GlobalMiddlewareStack
func After(name string) Middleware {
	return GlobalMiddlewareStack.After(name)
}

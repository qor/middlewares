# Middleware Stack

Manage Golang HTTP middlewares

## Usage

```go
func main() {
	Stack := &MiddlewareStack{}

	// Add middleware `auth` to stack
	Stack.Use(&middlewares.Middleware{
		Name: "auth",
		// Insert middleware `auth` after middleware `session` if it exists
		InsertAfter: []string{"session"},
		// Insert middleware `auth` before middleare `authorization` if it exists
		InsertBefore: []string{"authorization"},
	})

	// Remove middleware `cookie` from stack
	Stack.Remove("cookie")

	mux := http.NewServeMux()
	http.ListenAndServe(":9000", Stack.Apply(mux))
}
```

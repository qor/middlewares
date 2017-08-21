# Middlewares

Manage Golang HTTP middlewares

## Usage

```go
func main() {
	MiddlewaresStack := &Middlewares{}

	// Adding a Middleware
	MiddlewaresStack.Use(&middlewares.Middleware{
		Name: "auth",
		// Insert middleware `auth` after middleware `session` if it exists
		InsertAfter: []string{"session"},
		// Insert middleware `auth` before middleare `authorization` if it exists
		InsertBefore: []string{"authorization"},
	})

	// Removing middleware by name
	MiddlewaresStack.Remove("cookie")

	mux := http.NewServeMux()
	http.ListenAndServe(":9000", MiddlewaresStack.Apply(mux))
}
```

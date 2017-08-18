package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

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

// sortMiddlewares sort middlewares
func (middlewares *Middlewares) sortMiddlewares() []*Middleware {
	var (
		errs                         []error
		middlewareNames, sortedNames []string
		middlewaresMap               = map[string]*Middleware{}
		sortMiddleware               func(m *Middleware)
	)

	for _, middleware := range middlewares.middlewares {
		middlewaresMap[middleware.Name] = middleware
		middlewareNames = append(middlewareNames, middleware.Name)
	}

	for _, middleware := range middlewares.middlewares {
		for _, require := range middleware.Requires {
			if _, ok := middlewaresMap[require]; !ok {
				errs = append(errs, fmt.Errorf("middleware %v requires %v, but it doesn't exist", middleware.Name, require))
			}
		}
	}

	if len(errs) > 0 {
		panic(fmt.Sprint(errs))
	}

	sortMiddleware = func(m *Middleware) {
		if getRIndex(sortedNames, m.Name) == -1 { // if not sorted
			// sort by before
			var maxBeforeIndex int
			for _, before := range m.Before {
				idx := getRIndex(sortedNames, before)
				if idx != -1 {
					if idx > maxBeforeIndex {
						maxBeforeIndex = idx
					}
				} else if idx := getRIndex(middlewareNames, before); idx != -1 {
					sortedNames = append(sortedNames, m.Name)
					sortMiddleware(middlewares.middlewares[idx])
					// update middlewares after
					middlewares.middlewares[idx].After = uniqueAppend(middlewares.middlewares[idx].After, m.Name)
					maxBeforeIndex = len(sortedNames)
				}
			}

			// FIXME
			if maxBeforeIndex > 0 {
				sortedNames = append(sortedNames[:maxBeforeIndex+1], append([]string{m.Name}, sortedNames[maxBeforeIndex+1:]...)...)
			}

			// sort by after
			var minAfterIndex int
			for _, after := range m.After {
				idx := getRIndex(sortedNames, after)
				if idx != -1 {
					if idx < minAfterIndex {
						minAfterIndex = idx
					}
				} else if idx := getRIndex(middlewareNames, after); idx != -1 {
					middlewares.middlewares[idx].Before = uniqueAppend(middlewares.middlewares[idx].Before, m.Name)
				}
			}

			if minAfterIndex > 0 {
				sortedNames = append(sortedNames[:minAfterIndex+1], append([]string{m.Name}, sortedNames[minAfterIndex+1:]...)...)
			}

			// if current callback haven't been sorted, append it to last
			if getRIndex(sortedNames, m.Name) == -1 {
				sortedNames = append(sortedNames, m.Name)
			}
		}
	}

	for _, middleware := range middlewares.middlewares {
		sortMiddleware(middleware)
	}

	var sortedMiddlewares []*Middleware
	for _, name := range sortedNames {
		sortedMiddlewares = append(sortedMiddlewares, middlewaresMap[name])
	}

	return sortedMiddlewares
}

func (middlewares *Middlewares) String() string {
	var (
		sortedNames       []string
		sortedMiddlewares = middlewares.sortMiddlewares()
	)

	for _, middleware := range sortedMiddlewares {
		sortedNames = append(sortedNames, middleware.Name)
	}

	return fmt.Sprintf("Middlewares: %v", strings.Join(sortedNames, ", "))
}

// Apply apply middlewares to handler
func (middlewares *Middlewares) Apply(handler http.Handler) http.Handler {
	var (
		compiledHandler   http.Handler
		sortedMiddlewares = middlewares.sortMiddlewares()
	)

	for idx := len(sortedMiddlewares) - 1; idx >= 0; idx-- {
		middleware := sortedMiddlewares[idx]

		if compiledHandler == nil {
			compiledHandler = middleware.Handler(handler)
		} else {
			compiledHandler = middleware.Handler(compiledHandler)
		}
	}

	return compiledHandler
}

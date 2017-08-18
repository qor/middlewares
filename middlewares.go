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

		for _, insertBefore := range middleware.InsertBefore {
			if m, ok := middlewaresMap[insertBefore]; ok {
				m.InsertAfter = uniqueAppend(m.InsertAfter, middleware.Name)
			}
		}

		for _, insertAfter := range middleware.InsertAfter {
			if m, ok := middlewaresMap[insertAfter]; ok {
				m.InsertBefore = uniqueAppend(m.InsertBefore, middleware.Name)
			}
		}
	}

	if len(errs) > 0 {
		panic(fmt.Sprint(errs))
	}

	sortMiddleware = func(m *Middleware) {
		if _, found := getRIndex(sortedNames, m.Name); !found { // if not sorted
			var minIndex = -1

			// sort by InsertAfter
			for _, insertAfter := range m.InsertAfter {
				idx, found := getRIndex(sortedNames, insertAfter)
				if !found {
					if middleware, ok := middlewaresMap[insertAfter]; ok {
						sortMiddleware(middleware)
						idx, found = getRIndex(sortedNames, insertAfter)
					}
				}

				if found && idx > minIndex {
					minIndex = idx
				}
			}

			// sort by InsertBefore
			for _, insertBefore := range m.InsertBefore {
				if idx, found := getRIndex(sortedNames, insertBefore); found {
					if idx < minIndex {
						sortedNames = append(sortedNames[:idx], sortedNames[idx+1:]...)
						sortMiddleware(middlewaresMap[insertBefore])
						return
					}
				}
			}

			if minIndex >= 0 {
				sortedNames = append(sortedNames[:minIndex+1], append([]string{m.Name}, sortedNames[minIndex+1:]...)...)
			} else if _, has := getRIndex(sortedNames, m.Name); !has {
				// if current callback haven't been sorted, append it to last
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

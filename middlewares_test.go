package middlewares

import "testing"

func TestCompileMiddlewares(t *testing.T) {
	middlewares := &Middlewares{}
	availableMiddlewares := []Middleware{{Name: "cookie"}, {Name: "flash", After: []string{"cookie"}}}

	for _, m := range availableMiddlewares {
		middlewares.Use(m)
	}

	testSortedMiddlewares(middlewares, []string{"cookie", "flash"}, t)
}

func testSortedMiddlewares(middlewares *Middlewares, sortedNames []string, t *testing.T) {
	sortedMiddlewares := middlewares.sortMiddlewares()

	if len(sortedMiddlewares) != len(sortedNames) {
		t.Errorf("Length should be same, but got %v, expect %v", middlewares.String(), sortedNames)
	}

	for idx, middleware := range sortedMiddlewares {
		if sortedNames[idx] != middleware.Name {
			t.Errorf("Expected sorted middleware is %v, but got %v", sortedNames, middlewares.String())
		}
	}
}

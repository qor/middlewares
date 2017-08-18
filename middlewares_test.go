package middlewares

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func registerMiddlewareRandomly(registeredMiddlewares []Middleware) *Middlewares {
	middlewares := &Middlewares{}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	sort.Slice(registeredMiddlewares, func(i, j int) bool {
		return r.Intn(100)%2 == 1
	})

	for _, m := range registeredMiddlewares {
		middlewares.Use(m)
	}

	return middlewares
}

func registerMiddleware(registeredMiddlewares []Middleware) *Middlewares {
	middlewares := &Middlewares{}

	for _, m := range registeredMiddlewares {
		middlewares.Use(m)
	}

	return middlewares
}

func checkSortedMiddlewares(middlewares *Middlewares, sortedNames []string, t *testing.T) {
	sortedMiddlewares := middlewares.sortMiddlewares()

	if len(sortedMiddlewares) != len(sortedNames) {
		t.Errorf("Length should be same, but got %v, expect %v", middlewares.String(), strings.Join(sortedNames, ", "))
		return
	}

	for idx, middleware := range sortedMiddlewares {
		if sortedNames[idx] != middleware.Name {
			t.Errorf("Expected sorted middleware is %v, but got %v", strings.Join(sortedNames, ", "), middlewares.String())
		}
	}
}

func TestCompileMiddlewares(t *testing.T) {
	availableMiddlewares := []Middleware{{Name: "cookie"}, {Name: "flash", After: []string{"cookie"}}, {Name: "auth", Before: []string{"flash"}}}

	middlewares := registerMiddlewareRandomly(availableMiddlewares)
	checkSortedMiddlewares(middlewares, []string{"cookie", "flash", "auth"}, t)
}

func TestConflictingMiddlewares(t *testing.T) {
	t.Skipf("conflicting middlewares")
}

func TestMiddlewaresWithRequires(t *testing.T) {
	t.Skipf("conflicting middlewares")
}

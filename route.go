package tango

import "strings"

type RoutePath struct {
	Name  string
	Match bool
}

// Route struct is the glue between a RoutePath and the Component available on that path
type Route struct {
	// TODO: add guards etc.
	Path  []RoutePath
	scope *Scope
	root  Component
}

// NewRoute will make a Route struct out of the path (using makeRoute) and the Component (Controller)
func NewRoute(path string, root Component) Route {
	return Route{
		Path: makeRoute(path),
		root: root,
	}
}

// makeRoute parses a path with placeholder syntax
// eg '/path/to/:id' where :id can be anything
func makeRoute(r string) (path []RoutePath) {
	if len(r) > 0 && r[0] == '/' {

		s := strings.Split(r, "/")

		for _, a := range s {
			if len(a) > 0 && a[0] == ':' {
				path = append(path, RoutePath{
					Name:  a[1:],
					Match: false,
				})
			} else {
				path = append(path, RoutePath{
					Name:  a,
					Match: true,
				})
			}
		}

	}

	return path
}

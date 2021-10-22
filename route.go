package tango

import "strings"

type RoutePath struct {
	Name  string
	Match bool
}

type Route struct {
	// TODO: add guards etc.
	Path  []RoutePath
	scope *Scope
	root  Component
}

func NewRoute(path string, root Component) Route {
	return Route{
		Path: makeRoute(path),
		root: root,
	}
}

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

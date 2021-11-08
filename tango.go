package tango

import (
	"crypto/rand"
	"fmt"
	"strings"
	"syscall/js"
)

type Queue struct {
	Render []func()
	Post   []func()
}

type Tango struct {
	scope      *Scope
	components []Component
	routes     []Route
	Root       js.Value
}

func New() *Tango {
	return &Tango{
		scope: NewScope(nil),
	}
}

func (t *Tango) GenId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (t *Tango) AddRoutes(routes ...Route) {
	t.routes = routes
}

func (t *Tango) AddComponents(components ...Component) {
	t.components = components
}

func (t *Tango) Bootstrap() {
	js.Global().Get("window").Call("addEventListener", "hashchange", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			hash := js.Global().Get("window").Get("location").Get("hash").String()
			defer t.navigate(hash[2:])
			return nil
		},
	))
	t.finish(t.scope, js.Global().Get("document").Call("getElementsByTagName", "body").Index(0))
}

func (t *Tango) matchRoute(path string) (route *Route, params map[string]string) {
	params = make(map[string]string)
	splt := strings.Split(path, "/")
	for _, r := range t.routes {
		if len(splt) == len(r.Path) {
			for i, s := range splt {
				if r.Path[i].Match && r.Path[i].Name == s {
					route = &r
				} else if r.Path[i].Match && r.Path[i].Name != s {
					route = nil
					break
				} else {
					params[r.Path[i].Name] = s
				}
			}
			if route != nil {
				break
			}
		} else {
			continue
		}
	}
	return
}

func (t *Tango) navigate(path string) {
	route, attrs := t.matchRoute(path)
	if route != nil {
		construct := false
		if route.scope == nil {
			route.scope = NewScope(t.scope)
			construct = true
		}
		hook := Hook{
			Self:  t,
			Scope: route.scope,
			Attrs: attrs,
			Node:  t.Root,
			Queue: nil,
		}
		route.root.BeforeRender(hook)
		t.render(route.root, construct, hook)
		route.root.AfterRender(hook)
		t.finish(route.scope, t.Root)
	} else {
		println("route not found: " + path)
	}
}

func (t Tango) Nav(path string) {
	js.Global().Get("window").Get("location").Set("hash", "!"+path) // unsafe
}

func (t *Tango) finish(scope *Scope, node js.Value) {
	var queue Queue
	t.Compile(scope, node, &queue)
	scope.Digest()
	// process post queue when everything is finished
	// prevent directives to pickup pre-processed content (like router)
	for _, p := range queue.Post {
		p()
	}
}

func (t *Tango) Compile(scope *Scope, node js.Value, queue *Queue) {
	stop := false
	if !node.Equal(t.Root) {
		stop = t.exec(scope, node, queue)
	}
	if !stop {
		children := node.Get("children")
		for i := 0; i < children.Length(); i++ {
			t.Compile(scope, children.Index(i), queue)
		}
	}
}

func (t *Tango) exec(scope *Scope, node js.Value, queue *Queue) bool {
	stop := false
	m := make(map[string]string)

	// collect all attributes in a single map
	p := node.Get("attributes")
	for j := 0; j < p.Length(); j++ {
		name := p.Index(j).Get("nodeName")
		val := p.Index(j).Get("nodeValue")
		m[name.String()] = val.String()
	}

	// check for unique tagName directive matches...
	var components []Component = nil
	tn := node.Get("tagName").String()
	for _, c := range t.components {
		if c.Config().Kind == Tag && strings.ToLower(c.Config().Name) == strings.ToLower(tn) {
			components = append(components, c)
			break
		}
	}

	// if tagName is a known HTML5 tag, process attribute directives
	if len(components) == 0 {
		for _, c := range t.components {
			if _, e := m[c.Config().Name]; e {
				if c.Config().Kind == Attribute {
					components = append(components, c)
				}
			}
		}
	}

	// render the component
	if len(components) > 0 {
		// retrieve the node element's id
		construct := false
		var id string
		if n, e := m["tng-id"]; !e {
			id = t.GenId()
			node.Call("setAttribute", "tng-id", id)
			construct = true
		} else {
			id = n
		}

		for _, component := range components {
			var local *Scope
			if component.Config().Scoped {
				if c, e := scope.children[id]; !e {
					local = NewScope(scope)
					scope.children[id] = local
				} else {
					local = c
				}
			} else {
				local = scope
			}
			hook := Hook{
				Self:  t,
				Scope: local,
				Attrs: m,
				Node:  node,
				Queue: queue,
			}
			stop = t.render(component, construct, hook)
		}
	}
	return stop
}

func (t *Tango) render(component Component, construct bool, hook Hook) (stop bool) {
	if construct {
		stop = !component.Constructor(hook)
	}
	if component.Config().Kind != Attribute {
		hook.Node.Set("innerHTML", component.Render())
	}
	return
}

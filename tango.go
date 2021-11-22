package tango

import (
	"crypto/rand"
	"fmt"
	"strings"
	"syscall/js"
)

// Tango is an opinionated go(lang) WASM/SPA framework based on React & Angular JS.
// See the tango-example for a full batteries included example of what is possible only using go & html/css.
type Tango struct {
	scope      *Scope
	components []Component
	routes     []Route
	Root       js.Value
}

// New returns a fresh tango instance with an empty root Scope
func New() *Tango {
	return &Tango{
		scope: NewScope(nil),
	}
}

// Queue is able to collect functions that needs to be run after the recursion is done.
// This means Queue is not nil (when available in a Hook) if the hook is part of a recursion.
type Queue struct {
	Post []func()
}

// Hook is a passable struct used to give access to the most common internals
// the Hook's are available during different phases of the Component lifecycle
type Hook struct {
	Self  *Tango
	Scope *Scope            // local scope where the Node resides
	Attrs map[string]string // attrs can be mapped from a RoutePath or Attributes from an HTMLElement
	Node  js.Value          // Element or value from the view that has been 'hooked'
	Queue *Queue            // not nil when recursive
}

func (h *Hook) Run(attr string) {
	h.Scope.Exec(h.Attrs[attr], h)
}

// Get is an alias for Scope.Get
func (h *Hook) Get(name string) (js.Value, bool) {
	return h.Scope.Get(name)
}

// Absorb is an alias for Scope.Absorb
func (h *Hook) Absorb(i interface{}) {
	h.Scope.Absorb(i)
}

// GenId creates a random ID used to differentiate objects within the DOM
func (t *Tango) GenId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// AddRoutes is used to add your Route's to Tango
func (t *Tango) AddRoutes(routes ...Route) {
	t.routes = routes
}

// AddComponents can be used to add your Component's to Tango
func (t *Tango) AddComponents(components ...Component) {
	t.components = components
}

// Bootstrap has to be called to initialize the Tango SPA when everything has been set up
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

// matchRoute can be used to retrieve the Route and a map of parameters based on the path input
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

// navigate will matchRoute the path and render the Route's Controller Component
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
		t.finish(route.scope, t.Root)
		route.root.AfterRender(hook)
	} else {
		println("route not found: " + path)
	}
}

// Nav will trigger the onChange event that will trigger navigate based on the path
func (t Tango) Nav(path string) {
	js.Global().Get("window").Get("location").Set("hash", "!"+path) // unsafe
}

// finish is called when BeforeRender and render are done.
// it will compile the scope and run the Queue when done.
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

// Compile runs exec on each Component found within the DOM matching the criteria of the Config
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

// exec is called recursively from Compile to touch all DOM elements
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

// render is a helper func to set the innerHTML of a template provided by a Component
func (t *Tango) render(component Component, construct bool, hook Hook) (stop bool) {
	if construct {
		stop = !component.Constructor(hook)
	}
	if component.Config().Kind != Attribute {
		hook.Node.Set("innerHTML", component.Render())
	}
	return
}

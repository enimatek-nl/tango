package tango

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"syscall/js"
)

type Queue struct {
	Render []func()
	Post   []func()
}

type Route struct {
	// TODO: add guards etc.
	scope *Scope
	root  Component
}

type Tango struct {
	scope      *Scope
	components []Component
	routes     map[string]Route
	Root       js.Value
}

func New() *Tango {
	return &Tango{
		scope:  NewScope(nil),
		routes: make(map[string]Route),
	}
}

func (t *Tango) GenId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (t *Tango) AddRoute(path string, component Component) {
	t.routes[path] = Route{
		root: component,
	}
}

func (t *Tango) AddComponents(components ...Component) {
	t.components = components
}

func (t *Tango) Bootstrap() {
	js.Global().Get("window").Call("addEventListener", "hashchange", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			hash := js.Global().Get("window").Get("location").Get("hash").String()
			defer t.Navigate(hash[2:])
			return nil
		},
	))
	t.finish(t.scope, js.Global().Get("document").Call("getElementsByTagName", "body").Index(0))
}

func (t *Tango) Navigate(path string) {
	if route, exists := t.routes[path]; exists {
		if route.scope == nil {
			route.scope = NewScope(t.scope)
			route.root.Constructor(t, route.scope, t.Root, nil, nil)
		}
		route.root.BeforeRender(route.scope)
		t.Root.Set("innerHTML", route.root.Render())
		t.finish(route.scope, t.Root)
		route.root.AfterRender(route.scope)
	} else {
		panic("route not found")
	}
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
	if !node.Equal(t.Root) {
		t.exec(scope, node, queue)
	}
	children := node.Get("children")
	for i := 0; i < children.Length(); i++ {
		t.Compile(scope, children.Index(i), queue)
	}
}

func (t *Tango) exec(scope *Scope, node js.Value, queue *Queue) {
	m := make(map[string]js.Value)

	// collect all attributes in a single map
	p := node.Get("attributes")
	for j := 0; j < p.Length(); j++ {
		name := p.Index(j).Get("nodeName")
		val := p.Index(j).Get("nodeValue")
		m[name.String()] = val
	}

	// retrieve the node element's id
	construct := false
	var id string
	if n, e := m["tng-id"]; !e {
		id = t.GenId()
		node.Call("setAttribute", "tng-id", id)
		construct = true
	} else {
		id = n.String()
	}

	// check for unique tagName directive matches...
	var component Component = nil
	tn := node.Get("tagName").String()
	for _, c := range t.components {
		if c.Kind() == Tag && strings.ToLower(c.Name()) == strings.ToLower(tn) {
			component = c
			break
		}
	}

	// if tagName is a known HTML5 tag, process attribute directives
	if component == nil {
		for _, c := range t.components {
			if _, e := m[c.Name()]; e {
				if c.Kind() == Attribute {
					component = c
				}
			}
		}
	}

	// render the component
	if component != nil {
		var local *Scope
		if component.Scoped() {
			if c, e := scope.children[id]; !e {
				local = NewScope(scope)
				scope.children[id] = local
			} else {
				local = c
			}
		} else {
			local = scope
		}
		if construct {
			component.Constructor(t, local, node, m, queue)
		}
		if component.Kind() == Tag {
			component.BeforeRender(local)
			node.Set("innerHTML", component.Render())
			component.AfterRender(local)
		}
	}
}

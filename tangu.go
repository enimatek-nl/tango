package tango

import (
	"strings"
	"syscall/js"
)

type Queue struct {
	Post []func()
}

type Lifecycle int

const (
	Constructor Lifecycle = iota
	BeforeRender
	AfterRender
)

type Route struct {
	controller *Controller
}

type Tangu struct {
	scope      *Scope
	directives []Directive
	routes     map[string]Route
	Root       js.Value
}

func New() *Tangu {
	return &Tangu{
		scope:  NewScope(nil),
		routes: make(map[string]Route),
	}
}

func (t *Tangu) AddRoute(path string, controller *Controller) {
	t.routes[path] = Route{controller: controller}
}

func (t *Tangu) AddDirective(directives ...Directive) {
	t.directives = directives
}

func (t *Tangu) Bootstrap() {
	js.Global().Get("window").Call("addEventListener", "hashchange", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			hash := js.Global().Get("window").Get("location").Get("hash").String()
			defer t.Navigate(hash[2:])
			return nil
		},
	))
	t.finish(t.scope, js.Global().Get("document").Call("getElementsByTagName", "body").Index(0))
}

func (t *Tangu) Navigate(path string) {
	if route, exists := t.routes[path]; exists {
		if route.controller.Scope == nil {
			route.controller.Scope = NewScope(t.scope)
			route.controller.Work(route.controller.Scope, Constructor)
		}
		route.controller.Work(route.controller.Scope, BeforeRender)
		t.Root.Set("innerHTML", route.controller.Template())
		t.finish(route.controller.Scope, t.Root)
		route.controller.Work(route.controller.Scope, AfterRender)
	} else {
		panic("route not found")
	}
}

func (t *Tangu) finish(scope *Scope, node js.Value) {
	var queue Queue
	t.Compile(scope, node, &queue)
	scope.Digest()
	// process post queue when everything is finished
	// prevent directives to pickup pre-processed content (like router)
	for _, p := range queue.Post {
		p()
	}
}

func (t *Tangu) Compile(scope *Scope, node js.Value, queue *Queue) {
	if !node.Equal(t.Root) {
		t.exec(scope, node, queue)
	}
	children := node.Get("children")
	println(node.Get("tagName").String() + ": " + string(children.Length()))
	for i := 0; i < children.Length(); i++ {
		t.Compile(scope, children.Index(i), queue)
	}
}

func (t *Tangu) exec(scope *Scope, node js.Value, queue *Queue) {
	m := make(map[string]js.Value)

	// collect all attributes in a single map
	p := node.Get("attributes")
	for j := 0; j < p.Length(); j++ {
		name := p.Index(j).Get("nodeName")
		val := p.Index(j).Get("nodeValue")
		m[name.String()] = val
	}

	// check for unique tagName directive matches...
	tn := node.Get("tagName").String()
	for _, d := range t.directives {
		if d.Kind() == Tag && strings.ToLower(d.Name()) == strings.ToLower(tn) {
			d.Callback(t, scope, node, m, queue)
			return
		}
	}

	// if tagName is a known HTML5 tag, process attribute directives
	for _, d := range t.directives {
		if _, e := m[d.Name()]; e {
			if d.Kind() == Attribute {
				d.Callback(t, scope, node, m, queue)
			}
		}
	}

}

package tango

import (
	"strings"
	"syscall/js"
)

type Subscription struct {
	name     string
	callback func(scope *Scope, value js.Value)
	previous js.Value
}

type SModel struct {
	values    map[string]js.Value
	functions map[string]func(value js.Value, scope *Scope)
}

type Scope struct {
	model         SModel
	parent        *Scope
	subscriptions []*Subscription
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		model: SModel{
			values:    make(map[string]js.Value),
			functions: make(map[string]func(value js.Value, scope *Scope)),
		},
		parent: parent,
	}
}

func (s *Scope) AddFunc(name string, f func(value js.Value, scope *Scope)) {
	s.model.functions[name] = f
}

func (s *Scope) Add(name string, value js.Value) {
	parts := strings.Split(name, ".")
	last := len(parts) - 1
	if len(parts) == 1 {
		s.model.values[name] = value
	} else {
		exists := false
		var id string
		var obj js.Value
		for i, p := range parts {
			if i == 0 {
				id = p
				obj, exists = s.model.values[id]
			} else {
				if i != last {
					if exists {
						obj = obj.Get(p)
					} else {
						obj = js.ValueOf(make(map[string]interface{}))
						s.model.values[name] = obj
					}
				} else {
					obj.Set(p, value)
				}
			}
		}
	}
}

func (s *Scope) Get(name string) (js.Value, bool) {
	parts := strings.Split(name, ".")
	exists := false
	var id string
	var obj js.Value
	c := 0

	for i, p := range parts {
		c++
		if i == 0 {
			id = p
			if o, e := s.model.values[id]; e {
				obj = o
				exists = true
			} else {
				if s.parent != nil {
					return s.parent.Get(name)
				} else {
					println("'" + name + "' not found in scope model")
					return obj, false
				}
			}
		} else {
			if s.model.values[id].Call("hasOwnProperty", p).Bool() {
				obj = obj.Get(p)
			}
		}
	}
	return obj, exists
}

func (s *Scope) Digest() {
	for _, sub := range s.subscriptions {
		if v, e := s.Get(sub.name); e {
			if !v.Equal(sub.previous) {
				sub.previous = v
				sub.callback(s, v)
			}
		} else {
			println("subscription name not found in scope model")
		}
	}
	if s.parent != nil {
		s.parent.Digest()
	}
}

func (s *Scope) Clone() *Scope {
	clone := &Scope{
		model: SModel{
			values:    make(map[string]js.Value),
			functions: make(map[string]func(value js.Value, scope *Scope)),
		},
		parent: s,
	}
	return clone
}

func (s *Scope) AddSubscription(name string, f func(scope *Scope, value js.Value)) {
	s.subscriptions = append(s.subscriptions, &Subscription{
		name:     name,
		callback: f,
		previous: js.Value{},
	})
	println("subscription added for", name)
}

func (s *Scope) Destroy() {
	// TODO: clean ?
}

func (s *Scope) Exec(node js.Value, scope *Scope, valueOf js.Value) {
	id := valueOf.String()
	if f, exists := s.model.functions[id]; exists {
		f(node, scope)
		return
	} else if s.parent != nil {
		s.parent.Exec(node, scope, valueOf)
	} else {
		println("'" + id + "' function name not found in scope model")
	}
}

func (s *Scope) Root() *Scope {
	if s.parent != nil {
		return s.Root()
	} else {
		return s
	}
}

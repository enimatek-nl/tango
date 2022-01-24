//go:build js && wasm

package tango

import (
	"encoding/json"
	"github.com/enimatek-nl/tango/vert"
	"reflect"
	"strings"
	"syscall/js"
)

// SFunc type defines a loc(al) Hook used within JS functions and callbacks
type SFunc func(self *Tango, this js.Value, local *Scope)

// SModel stands for Scope Model it contains all values (properties) and functions (methods) that will be available in the DOM (view)
type SModel struct {
	values    map[string]js.Value
	functions map[string]SFunc
}

type Subscription struct {
	name     string
	callback func(scope *Scope, value js.Value)
	previous js.Value
}

// Scope is the binding part between the HTML (view) and the go-code (controller)
type Scope struct {
	model         SModel
	parent        *Scope
	subscriptions []*Subscription
	children      map[string]*Scope
}

// NewScope create a new scope based on the parent
func NewScope(parent *Scope) *Scope {
	return &Scope{
		model: SModel{
			values:    make(map[string]js.Value),
			functions: make(map[string]SFunc),
		},
		parent: parent,
	}
}

// Absorb will reflect all tng tags and map them to the SModel
func (s *Scope) Absorb(i interface{}) bool {
	d := false // was something digested?

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		t := reflect.TypeOf(i).Elem()
		rv := reflect.ValueOf(i).Elem()

		num := t.NumField()

		for i := 0; i < num; i++ {
			sf := t.Field(i)
			if tag, ok := sf.Tag.Lookup("tng"); ok {
				if !d {
					d = true
				}
				fv := rv.Field(i)

				if fv.Kind() == reflect.Func {
					s.model.functions[tag] = fv.Interface().(SFunc)
				} else {
					s.Set(tag, fv.Interface())
				}
			}
		}
	}

	return d
}

// Extract will reflect all tng tags and retrieve all the values from the SModel
func (s *Scope) Extract(i interface{}) {
	t := reflect.TypeOf(i).Elem()
	rv := reflect.ValueOf(i).Elem()

	num := t.NumField()
	for i := 0; i < num; i++ {
		sf := t.Field(i)
		if tag, ok := sf.Tag.Lookup("tng"); ok {
			fv := rv.Field(i)
			if fv.Kind() != reflect.Func {

				if v, o := s.Get(tag); o {
					// get js.Value .. check against kind (struct? etc.)
					// 'vert' it into the correct value somehow and set it :)
					//fv.Set(v.JSValue())
					println(v.String())
				}

			}
		}
	}
}

// SetFunc will add a function shared by name with the DOM (view)
func (s *Scope) SetFunc(name string, f SFunc) {
	s.model.functions[name] = f
}

// Get the name index as js.Value
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
					return obj, false // name not found = false
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

// Set will add a value/property shared by name with the DOM (view)
func (s *Scope) Set(name string, value interface{}) {
	parts := strings.Split(name, ".")
	last := len(parts) - 1
	if len(parts) == 1 {
		s.model.values[name] = vert.ValueOf(value)
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

// Parent retrieves the parent of the Scope
func (s *Scope) Parent() *Scope {
	return s.parent
}

// GetJSON converts the object behind the name index to a JSON.stringify string
func (s *Scope) GetJSON(name string) string {
	if v, e := s.Get(name); e {
		return js.Global().Get("JSON").Call("stringify", v, js.FuncOf(replacer)).String()
	}
	return "{}"
}

// replacer is a function used by JSON.stringify
func replacer(this js.Value, inputs []js.Value) interface{} {
	if !inputs[1].Equal(js.Undefined()) || !inputs[1].Equal(js.Null()) {
		return inputs[1]
	} else {
		return js.Null()
	}
}

// Decode the name index into the struct interface
func (s *Scope) Decode(name string, i interface{}) error {
	return json.NewDecoder(strings.NewReader(s.GetJSON(name))).Decode(i)
}

// Digest all Subscription's made on the Scope
func (s *Scope) Digest() {
	for _, sub := range s.subscriptions {
		if v, e := s.Get(sub.name); e {
			if !v.Equal(sub.previous) {
				sub.previous = v
				sub.callback(s, v)
			}
		}
	}
	if s.parent != nil {
		s.parent.Digest()
	}
}

// Clone current Scope as a child of current Scope
func (s *Scope) Clone() *Scope {
	return NewScope(s)
}

// Subscribe on changes of the name index and receive the changed value
func (s *Scope) Subscribe(name string, f func(scope *Scope, value js.Value)) {
	s.subscriptions = append(s.subscriptions, &Subscription{
		name:     name,
		callback: f,
		previous: js.Value{},
	})
}

func (s *Scope) Destroy() {
	// TODO: clean ?
}

// Exec the underlying func of the name index
func (s *Scope) Exec(name string, hook *Hook) {
	if f, exists := s.model.functions[name]; exists {
		f(hook.Self, hook.Node, hook.Scope)
		return
	} else if s.parent != nil {
		s.parent.Exec(name, hook)
	} else {
		println("'" + name + "' function name not found in scope model")
	}
}

// Root will recursively retrieve the root Scope
func (s *Scope) Root() *Scope {
	if s.parent != nil {
		return s.Root()
	} else {
		return s
	}
}

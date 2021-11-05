package std

import (
	"github.com/enimatek-nl/tango"
	"strings"
	"syscall/js"
)

type Attr struct{}

func (a Attr) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-attr",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (a Attr) Constructor(hook tango.Hook) bool {
	if valueOf, e := hook.Attrs[a.Config().Name]; e {
		onlyWhen := true
		parts := strings.Split(valueOf, " when ")
		if len(parts) == 1 {
			parts = strings.Split(valueOf, " is ")
			onlyWhen = false
		}

		if len(parts) == 2 {
			if _, e := hook.Scope.Get(parts[1]); e {
				handle := func(v js.Value) {
					if onlyWhen {
						if v.Bool() {
							hook.Node.Call("setAttribute", parts[0], js.ValueOf(""))
						} else {
							hook.Node.Call("removeAttribute", parts[0])
						}
					} else {
						hook.Node.Call("setAttribute", parts[0], v)
					}
				}
				hook.Scope.Subscribe(parts[1], func(scope *tango.Scope, value js.Value) {
					handle(value)
				})
			}
		} else {
			panic("can't parse '" + valueOf + "'")
		}
	} else {
		panic(a.Config().Name + " attribute not set")
	}
	return true
}

func (a Attr) BeforeRender(hook tango.Hook) {}

func (a Attr) AfterRender(hook tango.Hook) {}

func (a Attr) Render() string { return "" }

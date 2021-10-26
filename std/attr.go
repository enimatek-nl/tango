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

func (a Attr) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		if valueOf, e := attrs[a.Config().Name]; e {
			onlyWhen := true
			parts := strings.Split(valueOf, " when ")
			if len(parts) == 1 {
				parts = strings.Split(valueOf, " is ")
				onlyWhen = false
			}

			if len(parts) == 2 {
				if _, e := scope.Get(parts[1]); e {
					handle := func(v js.Value) {
						if onlyWhen {
							if v.Bool() {
								node.Call("setAttribute", parts[0], js.ValueOf(""))
							} else {
								node.Call("removeAttribute", parts[0])
							}
						} else {
							node.Call("setAttribute", parts[0], v)
						}
					}
					scope.AddSubscription(parts[1], func(scope *tango.Scope, value js.Value) {
						handle(value)
					})
				}
			} else {
				panic("can't parse '" + valueOf + "'")
			}
		} else {
			panic(a.Config().Name + " attribute not set")
		}
	}
	return true
}

func (a Attr) Render() string { return "" }

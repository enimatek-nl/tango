package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Change struct{}

func (c Change) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-change",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (c Change) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		if valueOf, e := attrs[c.Config().Name]; e {
			node.Call("addEventListener", "change", js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					scope.Exec(node, scope, js.ValueOf(valueOf))
					return nil
				}),
			)
		} else {
			panic(c.Config().Name + " attribute not set")
		}
	}
	return true
}

func (c Change) Render() string { return "" }

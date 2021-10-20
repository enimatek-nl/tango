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

func (c Change) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
	if valueOf, e := attrs[c.Config().Name]; e {
		node.Call("addEventListener", "change", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				scope.Exec(node, scope, valueOf)
				return nil
			}),
		)
	} else {
		panic(c.Config().Name + " attribute not set")
	}
	return true
}

func (c Change) Hook(scope *tango.Scope, hook tango.ComponentHook) {}

func (c Change) Render() string { return "" }

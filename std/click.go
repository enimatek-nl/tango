package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Click struct{}

func (c Click) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-click",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (c Click) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
	node.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			args[0].Call("stopPropagation")
			args[0].Call("preventDefault")
			println("clicked: " + c.Config().Name)
			scope.Exec(node, scope, attrs[c.Config().Name])
			scope.Digest()
			return nil
		}),
	)
	return true
}

func (c Click) Hook(scope *tango.Scope, hook tango.ComponentHook) {}

func (c Click) Render() string { return "" }

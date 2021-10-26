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

func (c Click) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		node.Call("addEventListener", "click", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				args[0].Call("stopPropagation")
				args[0].Call("preventDefault")
				println("clicked: " + c.Config().Name)
				scope.Exec(node, scope, js.ValueOf(attrs[c.Config().Name]))
				scope.Digest()
				return nil
			}),
		)
	}
	return true
}

func (c Click) Render() string { return "" }

package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Click struct{}

func (c Click) Name() string {
	return "tng-click"
}

func (c Click) Kind() tango.Kind {
	return tango.Attribute
}

func (c Click) Scoped() bool {
	return false
}

func (c Click) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	node.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			args[0].Call("stopPropagation")
			args[0].Call("preventDefault")
			println("clicked: " + c.Name())
			scope.Exec(node, scope, attrs[c.Name()])
			scope.Digest()
			return nil
		}),
	)
}

func (c Click) BeforeRender(scope *tango.Scope) {}

func (c Click) Render() string { return "" }

func (c Click) AfterRender(scope *tango.Scope) {}

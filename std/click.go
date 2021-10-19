package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Click struct{}

func (c Click) Kind() tango.Kind {
	return tango.Attribute
}

func (c Click) Name() string {
	return "tng-click"
}

func (c Click) Callback(self *tango.Tangu, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	node.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			args[0].Call("stopPropagation")
			args[0].Call("preventDefault")
			scope.Exec(node, scope, attrs[c.Name()])
			scope.Digest()
			return nil
		}),
	)
}

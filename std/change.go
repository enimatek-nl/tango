package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Change struct{}

func (c Change) Name() string {
	return "tng-change"
}

func (c Change) Kind() tango.Kind {
	return tango.Attribute
}

func (c Change) Scoped() bool {
	return false
}

func (c Change) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	if valueOf, e := attrs[c.Name()]; e {
		node.Call("addEventListener", "change", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				scope.Exec(node, scope, valueOf)
				return nil
			}),
		)
	} else {
		panic(c.Name() + " attribute not set")
	}
}

func (c Change) BeforeRender(scope *tango.Scope) {}

func (c Change) Render() string { return "" }

func (c Change) AfterRender(scope *tango.Scope) {}

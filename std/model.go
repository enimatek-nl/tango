package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Model struct{}

func (m Model) Kind() tango.Kind {
	return tango.Attribute
}

func (m Model) Name() string {
	return "tng-model"
}

func (m Model) Callback(self *tango.Tangu, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	if valueOf, e := attrs[m.Name()]; e {
		act := "keyup"
		// TODO: more exceptions needed?
		if node.Get("nodeName").String() == "SELECT" {
			act = "change"
		}
		node.Call("addEventListener", act, js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				scope.Add(valueOf.String(), node.Get("value"))
				scope.Digest()
				return nil
			}),
		)
	} else {
		panic(m.Name() + " attribute not set")
	}
}

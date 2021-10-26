package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Model struct{}

func (m Model) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-model",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (m Model) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		if valueOf, e := attrs[m.Config().Name]; e {
			act := "keyup"
			// TODO: more exceptions needed?
			if node.Get("nodeName").String() == "SELECT" {
				act = "change"
			}
			node.Call("addEventListener", act, js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					scope.Add(valueOf, node.Get("value"))
					scope.Digest()
					return nil
				}),
			)
		} else {
			panic(m.Config().Name + " attribute not set")
		}
	}
	return true
}

func (m Model) Render() string { return "" }

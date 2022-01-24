//go:build js && wasm

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

func (m Model) Constructor(hook tango.Hook) bool {
	if valueOf, e := hook.Attrs[m.Config().Name]; e {
		act := "keyup"
		// TODO: more exceptions needed?
		if hook.Node.Get("nodeName").String() == "SELECT" {
			act = "change"
		}
		hook.Node.Call("addEventListener", act, js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				hook.Scope.Set(valueOf, hook.Node.Get("value"))
				hook.Scope.Digest()
				return nil
			}),
		)
	} else {
		panic(m.Config().Name + " attribute not set")
	}
	return true
}

func (m Model) Render() string { return "" }

func (m Model) AfterRender(hook tango.Hook) bool {
	return false
}

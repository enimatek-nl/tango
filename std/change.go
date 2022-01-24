//go:build js && wasm

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

func (c Change) Constructor(hook tango.Hook) bool {
	if valueOf, e := hook.Attrs[c.Config().Name]; e {
		hook.Node.Call("addEventListener", "change", js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				hook.Run(valueOf)
				return nil
			}),
		)
	} else {
		panic(c.Config().Name + " attribute not set")
	}
	return true
}

func (c Change) Render() string { return "" }

func (c Change) AfterRender(hook tango.Hook) bool {
	return false
}

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

func (c Click) Constructor(hook tango.Hook) bool {
	hook.Node.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			args[0].Call("stopPropagation")
			args[0].Call("preventDefault")
			hook.Run(c.Config().Name)
			hook.Scope.Digest()
			return nil
		}),
	)
	return true
}

func (c Click) BeforeRender(hook tango.Hook) {}

func (c Click) AfterRender(hook tango.Hook) {}

func (c Click) Render() string { return "" }

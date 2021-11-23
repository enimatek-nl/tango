package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Bind struct{}

func (b Bind) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-bind",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (b Bind) Constructor(hook tango.Hook) bool {
	if valueOf, e := hook.Attrs[b.Config().Name]; e {
		hook.Scope.Subscribe(valueOf, func(scope *tango.Scope, value js.Value) {
			// TODO: add more exceptions on elements.
			if hook.Node.Get("nodeName").String() == "INPUT" {
				hook.Node.Set("value", value)
			} else {
				hook.Node.Set("innerHTML", value)
			}
		})
	} else {
		panic(b.Config().Name + " attribute not set")
	}
	return true
}

func (b Bind) Render() string { return "" }

func (b Bind) AfterRender(hook tango.Hook) bool {
	return false
}

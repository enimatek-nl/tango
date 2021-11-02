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

func (b Bind) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		if valueOf, e := attrs[b.Config().Name]; e {
			if _, e := scope.Get(valueOf); e {
				scope.Subscribe(valueOf, func(scope *tango.Scope, value js.Value) {
					// TODO: based on element type
					node.Set("innerHTML", value)
					node.Set("value", value)
				})
			}
		} else {
			panic(b.Config().Name + " attribute not set")
		}
	}
	return true
}

func (b Bind) Render() string { return "" }

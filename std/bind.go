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

func (b Bind) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
	if valueOf, e := attrs[b.Config().Name]; e {
		if _, e := scope.Get(valueOf.String()); e {
			scope.AddSubscription(valueOf.String(), func(scope *tango.Scope, value js.Value) {
				// TODO: based on element type
				node.Set("innerHTML", value)
				node.Set("value", value)
			})
		}
	} else {
		panic(b.Config().Name + " attribute not set")
	}
	return true
}

func (b Bind) Hook(scope *tango.Scope, attrs map[string]string, hook tango.ComponentHook) {}

func (b Bind) Render() string { return "" }

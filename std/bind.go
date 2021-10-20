package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Bind struct{}

func (b Bind) Name() string {
	return "tng-bind"
}

func (b Bind) Kind() tango.Kind {
	return tango.Attribute
}

func (b Bind) Scoped() bool {
	return false
}

func (b Bind) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	if valueOf, e := attrs[b.Name()]; e {
		if _, e := scope.Get(valueOf.String()); e {
			scope.AddSubscription(valueOf.String(), func(scope *tango.Scope, value js.Value) {
				// TODO: based on element type
				node.Set("innerHTML", value)
				node.Set("value", value)
			})
		}
	} else {
		panic(b.Name() + " attribute not set")
	}
}

func (b Bind) BeforeRender(scope *tango.Scope) {}

func (b Bind) Render() string { return "" }

func (b Bind) AfterRender(scope *tango.Scope) {}

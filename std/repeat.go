package std

import (
	"github.com/enimatek-nl/tango"
	"strings"
	"syscall/js"
)

type Repeat struct{}

func (r Repeat) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "tng-repeat",
		Kind:   tango.Attribute,
		Scoped: false,
	}
}

func (r Repeat) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
	if valueOf, e := attrs[r.Config().Name]; e {
		id := self.GenId()
		parts := strings.Split(valueOf.String(), " in ")
		parentNode := node.Get("parentNode")
		document := js.Global().Get("document")

		placeholder := document.Call("createElement", "template")
		placeholder.Call("setAttribute", "tng-repeat-id", id)
		parentNode.Call("replaceChild", placeholder, node)

		if _, e := scope.Get(parts[1]); e {
			scope.AddSubscription(parts[1], func(scope *tango.Scope, value js.Value) {
				var nominated []js.Value
				children := parentNode.Get("children")
				for j := 0; j < children.Length(); j++ {
					child := children.Index(j)
					if child.Call("hasAttribute", "tng-repeat").Bool() ||
						child.Call("hasAttribute", "tng-repeat-item").Bool() {
						if child.Equal(node) || child.Call("getAttribute", "tng-repeat-id").String() == id {
							nominated = append(nominated, child)
						}
					}
				}

				for _, n := range nominated {
					parentNode.Call("removeChild", n)
				}

				for i := 0; i < value.Length(); i++ {
					clone := node.Call("cloneNode", js.ValueOf(true))
					clone.Call("removeAttribute", "tng-repeat")
					clone.Call("setAttribute", "tng-repeat-item", "")
					clone.Call("setAttribute", "tng-repeat-id", id)
					parentNode.Call("insertBefore", clone, placeholder)

					childScope := scope.Clone()
					item := value.Index(i)
					childScope.Add(parts[0], item)
					var childQueue tango.Queue

					self.Compile(childScope, clone, &childQueue)
					childScope.Digest()
				}
			})
		}
	} else {
		panic(r.Config().Name + " attribute not set")
	}
	return false
}

func (r Repeat) Hook(scope *tango.Scope, hook tango.ComponentHook) {}

func (r Repeat) Render() string { return "" }

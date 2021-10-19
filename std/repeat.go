package std

import (
	"github.com/enimatek-nl/tango"
	"strings"
	"syscall/js"
)

type Repeat struct{}

func (r Repeat) Kind() tango.Kind {
	return tango.Attribute
}

func (r Repeat) Name() string {
	return "tng-repeat"
}

func (r Repeat) Callback(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	if valueOf, e := attrs[r.Name()]; e {
		id := "random" // TODO: gen some id
		parts := strings.Split(valueOf.String(), " in ")
		parentNode := node.Get("parentNode")
		document := js.Global().Get("document")

		placeholder := document.Call("createElement", "template")
		placeholder.Call("setAttribute", "tng-id", id)
		parentNode.Call("replaceChild", placeholder, node)

		if _, e := scope.Get(parts[1]); e {
			scope.AddSubscription(parts[1], func(scope *tango.Scope, value js.Value) {
				var nominated []js.Value
				children := parentNode.Get("children")
				for j := 0; j < children.Length(); j++ {
					child := children.Index(j)
					if child.Call("hasAttribute", "tng-repeat").Bool() ||
						child.Call("hasAttribute", "tng-repeat-item").Bool() {
						if child.Equal(node) || child.Call("getAttribute", "tng-id").String() == id {
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
					clone.Call("setAttribute", "tng-repeat-item", js.ValueOf(""))
					clone.Call("setAttribute", "tng-id", js.ValueOf(id))
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
		panic(r.Name() + " attribute not set")
	}
}

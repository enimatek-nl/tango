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

func (r Repeat) Constructor(hook tango.Hook) bool {
	if valueOf, e := hook.Attrs[r.Config().Name]; e {
		id := hook.Self.GenId()
		parts := strings.Split(valueOf, " in ")
		parentNode := hook.Node.Get("parentNode")
		document := js.Global().Get("document")

		placeholder := document.Call("createElement", "template")
		placeholder.Call("setAttribute", "tng-repeat-id", id)
		parentNode.Call("replaceChild", placeholder, hook.Node)

		if _, e := hook.Scope.Get(parts[1]); e {
			hook.Scope.Subscribe(parts[1], func(scope *tango.Scope, value js.Value) {
				var nominated []js.Value
				children := parentNode.Get("children")
				for j := 0; j < children.Length(); j++ {
					child := children.Index(j)
					if child.Call("hasAttribute", "tng-repeat").Bool() ||
						child.Call("hasAttribute", "tng-repeat-item").Bool() {
						if child.Equal(hook.Node) || child.Call("getAttribute", "tng-repeat-id").String() == id {
							nominated = append(nominated, child)
						}
					}
				}

				for _, n := range nominated {
					parentNode.Call("removeChild", n)
				}

				for i := 0; i < value.Length(); i++ {
					clone := hook.Node.Call("cloneNode", js.ValueOf(true))
					clone.Call("removeAttribute", "tng-repeat")
					clone.Call("setAttribute", "tng-repeat-item", "")
					clone.Call("setAttribute", "tng-repeat-id", id)
					parentNode.Call("insertBefore", clone, placeholder)

					childScope := scope.Clone()
					item := value.Index(i)
					childScope.Set(parts[0], item)
					var childQueue tango.Queue

					hook.Self.Compile(childScope, clone, &childQueue)
					childScope.Digest()
				}
			})
		}
	} else {
		panic(r.Config().Name + " attribute not set")
	}
	return false
}

func (r Repeat) Render() string { return "" }

func (r Repeat) AfterRender(hook tango.Hook) bool {
	return false
}

package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Router struct{}

const PATH = "path"

func (r Router) Config() tango.ComponentConfig {
	return tango.ComponentConfig{
		Name:   "Router",
		Kind:   tango.Tag,
		Scoped: false,
	}
}

func (r Router) Constructor(hook tango.Hook) bool {
	hook.Self.Root = hook.Node
	if p, e := hook.Attrs[PATH]; e {
		js.Global().Get("window").Get("location").Set("hash", "") // clear current hash
		hook.Queue.Post = append(hook.Queue.Post, func() {
			hook.Self.Nav(p) // nav to defined default hash (picked up by 'hash change event')
		})
	} else {
		panic("don't forget to set a 'path=' attribute")
	}
	return true
}

func (r Router) Render() string { return "" }

func (r Router) AfterRender(hook tango.Hook) bool {
	return false
}

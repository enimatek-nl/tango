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

func (r Router) Hook(self *tango.Tango, scope *tango.Scope, hook tango.ComponentHook, attrs map[string]string, node js.Value, queue *tango.Queue) bool {
	switch hook {
	case tango.Construct:
		self.Root = node
		if p, e := attrs[PATH]; e {
			js.Global().Get("window").Get("location").Set("hash", "#!"+p)
			queue.Post = append(queue.Post, func() {
				self.Navigate(p)
			})
		} else {
			panic("don't forget to set a 'path=' attribute")
		}
	}
	return true
}

func (r Router) Render() string { return "" }

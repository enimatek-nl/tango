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

func (r Router) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
	self.Root = node
	if p, e := attrs[PATH]; e {
		js.Global().Get("window").Get("location").Set("hash", "#!"+p.String())
		queue.Post = append(queue.Post, func() {
			self.Navigate(p.String())
		})
	} else {
		panic("don't forget to set a 'path=' attribute")
	}
	return true
}

func (r Router) Hook(scope *tango.Scope, attrs map[string]string, hook tango.ComponentHook) {}

func (r Router) Render() string { return "" }

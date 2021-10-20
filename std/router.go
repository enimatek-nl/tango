package std

import (
	"github.com/enimatek-nl/tango"
	"syscall/js"
)

type Router struct{}

const PATH = "path"

func (r Router) Kind() tango.Kind {
	return tango.Tag
}

func (r Router) Name() string {
	return "Router"
}

func (r Router) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) {
	self.Root = node
	if p, e := attrs[PATH]; e {
		js.Global().Get("window").Get("location").Set("hash", "#!"+p.String())
		queue.Post = append(queue.Post, func() {
			self.Navigate(p.String())
		})
	} else {
		panic("don't forget to set a 'path=' attribute")
	}
}

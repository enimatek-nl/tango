package tango

import (
	"syscall/js"
)

type Kind int

const (
	Attribute Kind = iota
	Tag
)

type Directive interface {
	Kind() Kind
	Name() string
	Callback(self *Tango, scope *Scope, node js.Value, attrs map[string]js.Value, queue *Queue)
}

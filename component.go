package tango

import "syscall/js"

type Kind int

const (
	Attribute Kind = iota
	Tag
)

type Component interface {
	Name() string
	Kind() Kind
	Scoped() bool
	Constructor(self *Tango, scope *Scope, node js.Value, attrs map[string]js.Value, queue *Queue)
	BeforeRender(scope *Scope)
	Render() string
	AfterRender(scope *Scope)
}

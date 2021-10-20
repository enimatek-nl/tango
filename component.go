package tango

import "syscall/js"

type Kind int

const (
	Attribute Kind = iota
	Tag
)

type ComponentHook int

const (
	BeforeRender ComponentHook = iota
	AfterRender
)

type ComponentConfig struct {
	Name   string
	Kind   Kind
	Scoped bool
}

type Component interface {
	Config() ComponentConfig
	Constructor(self *Tango, scope *Scope, node js.Value, attrs map[string]js.Value, queue *Queue) bool
	Hook(scope *Scope, hook ComponentHook)
	Render() string
}

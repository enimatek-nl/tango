package tango

import "syscall/js"

type Kind int

const (
	Controller Kind = iota // Controller has ComponentHook hooks
	Attribute              // Attribute matches html-attributes and does not Render()
	Tag                    // Tag adds a custom tag to html and calls Render() without ComponentHook
)

type ComponentHook int

const (
	Construct ComponentHook = iota
	BeforeRender
	AfterRender
)

type ComponentConfig struct {
	Name   string
	Kind   Kind
	Scoped bool
}

type Hook struct {
	Self  *Tango
	Scope *Scope
	Attrs map[string]string
	Node  js.Value
	Queue *Queue
}

type Component interface {
	Config() ComponentConfig // ComponentConfig describes the details of the component
	Constructor(hook Hook) bool
	BeforeRender(hook Hook)
	Render() string // return a template of innerHTML for Controller and Tag
	AfterRender(hook Hook)
	// Hook(self *Tango, scope *Scope, hook ComponentHook, attrs map[string]string, node js.Value, queue *Queue) bool
}

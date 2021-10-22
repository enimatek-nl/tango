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
	BeforeRender ComponentHook = iota
	AfterRender
)

type ComponentConfig struct {
	Name   string
	Kind   Kind
	Scoped bool
}

type Component interface {
	Config() ComponentConfig                                                                            // ComponentConfig describes the details of the component
	Constructor(self *Tango, scope *Scope, node js.Value, attrs map[string]js.Value, queue *Queue) bool // return true to continue with construct on children
	Hook(scope *Scope, hook ComponentHook)                                                              // used with a Controller Kind
	Render() string                                                                                     // return a template of innerHTML for Controller and Tag
}

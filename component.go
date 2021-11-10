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

// ComponentConfig describes the Component's Kind and Name (details)
type ComponentConfig struct {
	Name   string
	Kind   Kind
	Scoped bool
}

// Hook is a passable struct used to give access to the most common internals
// the Hook's are available during different phases of the Component lifecycle
type Hook struct {
	Self  *Tango
	Scope *Scope
	Attrs map[string]string
	Node  js.Value
	Queue *Queue
}

// Component can be of different Kind's
// In all cases a ComponentConfig is mandatory
// Only the Tag and Controller Kind can use the Render func to return an HTML template
type Component interface {
	Config() ComponentConfig
	Constructor(hook Hook) bool
	BeforeRender(hook Hook)
	Render() string
	AfterRender(hook Hook)
}

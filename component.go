package tango

type Kind int

const (
	Controller Kind = iota // Controller has ComponentHook hooks
	Attribute              // Attribute matches html-attributes and does not Render()
	Tag                    // Tag adds a custom tag to html and calls Render() without ComponentHook
)

// ComponentConfig describes the Component's Kind and Name (details)
type ComponentConfig struct {
	Name   string
	Kind   Kind
	Scoped bool
}

// Component can be of different Kind's
// In all cases a ComponentConfig is mandatory
// Only the Tag and Controller Kind can use the Render func to return an HTML template
type Component interface {
	Config() ComponentConfig
	Constructor(tng Hook) bool // Scope values can be setup once first run of the Component with a new scope
	Render() string
	AfterRender(tng Hook) bool // Set SModel values and Scope.Digest them each time the route is hit
}

package tango

//type Controller interface {
//	Scope() *Scope
//	Work(scope *Scope, lifecycle Lifecycle)
//	Template() string
//}

type Controller struct {
	Scope    *Scope
	Work     func(scope *Scope, lifecycle Lifecycle)
	Template func() string
}

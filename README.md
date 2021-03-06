# Tango
_Lightweight WASM HTML / Javascript Framework_

### Intro

WebAssembly is nice, Go on the web is nice, so I ported [Tangu](https://github.com/enimatek-nl/tangu/) to Go and
WebAssembly.

Tangu is an AngularJS inspired project I started a while back to explore [nim](https://nim-lang.org).

Where Tangu is nim 100% transpiled Javascript, Tango is golang compiled to WASM.

### Usage

`GOOS=js GOARCH=wasm go get github.com/enimatek-nl/tango`

### Get Started

```go
func main() {

    tg := tango.New()
    
    tg.AddComponents(
        std.Router{},
        std.Repeat{},
        std.Click{},
        std.Bind{},
        std.Change{},
        std.Model{},
        std.Attr{})
    
    tg.AddRoute(tango.NewRoute("/", &ViewController{Hello: "hello?"}))
    
    tg.Bootstrap()

}
```

```go
type ViewController struct {
	Hello   string      `tng:"hello"`
	ClickMe tango.SFunc `tng:"clickme"`
}

func (v ViewController) Config() tango.ComponentConfig {
    return tango.ComponentConfig{
        Name:   "ViewController",
        Kind:   tango.Controller,
        Scoped: false,
    }
}

func (v *ViewController) Constructor(tng tango.Hook) bool {
    v.ClickMe = func(self *tango.Tango, this js.Value, local *tango.Scope) {
        v.Hello = "world!"
        tng.Digest(v)
    }

    return true
}

func (v ViewController) Render() string {
    return `
        <div>
            <p tng-bind="hello"></p>
            <button tng-click="clickme">click me!</button>
        </div>
    `
}

func (v ViewController) AfterRender(tng tango.Hook) bool {
	return false
}
```

### How to use this Web Framework (SPA) & Backend API as a single project?
Check out this [todo example project](https://github.com/enimatek-nl/tango-example) as a reference implementation including Makefile and project structure.

_This project is a WIP..._

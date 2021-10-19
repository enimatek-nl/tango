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
    tg := tango.New()
    // add standard directives
    tg.AddDirective(
        std.Router{},
        std.Repeat{},
        std.Click{},
        std.Bind{},
        std.Change{},
        std.Model{},
        std.Attr{})
    // add a default path
    tg.AddRoute("/", &tango.Controller{
        Template: func () string {
        // return some html 
        },
        Work: func (scope *core.Scope, lifecycle core.Lifecycle) {
        // add logic to the scope
        }
    })
    // bootstrap tango
    tg.Bootstrap()
```

_This project is a WIP..._
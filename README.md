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
    
    tg.AddRoute("/", &ViewController{})
    
    tg.Bootstrap()

}
```

```go
type ViewController struct { }

func (v ViewController) Config() tango.ComponentConfig {
    return tango.ComponentConfig{
        Name:   "ViewController",
        Kind:   tango.Tag,
        Scoped: true,
    }
}

func (v ViewController) Constructor(self *tango.Tango, scope *tango.Scope, node js.Value, attrs map[string]js.Value, queue *tango.Queue) bool {
    scope.AddFunc("clickFunc", func(value js.Value, scope *tango.Scope) {
        println("hello world!")
    })
}

func (v ViewController) Hook(scope *tango.Scope, hook tango.ComponentHook) { }

func (v ViewController) Render() string {
    return `<button tng-click="clickFunc">click me!</button>`   
}
```

_This project is a WIP..._
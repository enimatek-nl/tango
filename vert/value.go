package vert

import (
	"reflect"
	"syscall/js"
	"time"
)

// Modified version of github.com/norunners/vert's value.go
// .. this is mandatory to convert 'complicated' structs 'generic' into a js.Value

var (
	null   = js.ValueOf(nil)
	object = js.Global().Get("Object")
	array  = js.Global().Get("Array")
)

func ValueOf(i interface{}) js.Value {
	//j, _ := json.Marshal(i)
	switch i.(type) {
	case nil, js.Value, js.Wrapper:
		return js.ValueOf(i)
	default:
		v := reflect.ValueOf(i)
		return valueOf(v)
	}
}

func valueOf(v reflect.Value) js.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return valueOfPointerOrInterface(v)
	case reflect.Slice, reflect.Array:
		return valueOfSliceOrArray(v)
	case reflect.Map:
		return valueOfMap(v)
	case reflect.Struct:
		return valueOfStruct(v)
	default:
		return js.ValueOf(v.Interface())
	}
}

// valueOfPointerOrInterface returns a new value.
func valueOfPointerOrInterface(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	return valueOf(v.Elem())
}

// valueOfSliceOrArray returns a new array object value.
func valueOfSliceOrArray(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	a := array.New()
	n := v.Len()
	for i := 0; i < n; i++ {
		e := v.Index(i)
		a.SetIndex(i, valueOf(e))
	}
	return a
}

// valueOfMap returns a new object value.
// Map keys must be of type string.
func valueOfMap(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	m := object.New()
	i := v.MapRange()
	for i.Next() {
		k := i.Key().Interface().(string)
		m.Set(k, valueOf(i.Value()))
	}
	return m
}

// valueOfStruct returns a new object value.
func valueOfStruct(v reflect.Value) js.Value {
	s := object.New()
	deepFields(&s, v)
	return s
}

// deepFields Recursively add fields to an object
// 1. this will ensure Anonymous embedded structs will be flattened
// 2. when a 'Valid bool' is part of the struct (think nullable sql structs) it will be ignored when false
func deepFields(s *js.Value, v reflect.Value) {
	t := v.Type()
	n := v.NumField()
	for i := 0; i < n; i++ {
		if f := v.Field(i); f.CanInterface() {
			sf := t.Field(i)
			k := nameOf(sf)
			if f.Type().Kind() == reflect.Struct { // ignore '!Valid' structs (sql nullable)
				u := f.FieldByName("Valid")
				if u.IsValid() && u.Type().Kind() == reflect.Bool {
					if !u.Interface().(bool) {
						s.Set(k, js.Undefined())
						break
					}
				}
			}
			if sf.Type.PkgPath() == "time" && sf.Type.Name() == "Time" { // parse time.Time into JS JSON standard
				t := f.Interface().(time.Time).UTC().Format("2006-01-02T15:04:05Z")
				s.Set(k, t)
			} else if sf.Anonymous { // flatten embedded structs
				deepFields(s, f)
			} else {
				s.Set(k, valueOf(f))
			}
		}
	}
}

// nameOf returns the JS tag name, otherwise the field name.
func nameOf(sf reflect.StructField) string {
	name := sf.Tag.Get("js")
	if name == "" {
		name = sf.Tag.Get("json")
	}
	if name == "" {
		return sf.Name
	}
	return name
}

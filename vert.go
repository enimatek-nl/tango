package tango

import (
	"reflect"
	"syscall/js"
)

// Heavily based on some pre-work done by: https://github.com/norunners/vert/
// .. this is mandatory to convert 'complicated' structs 'generic' into a js.Value
// I think it should be part of the js.Value std lib

var (
	null   = js.ValueOf(nil)
	object = js.Global().Get("Object")
	array  = js.Global().Get("Array")
)

func ValueOf(in interface{}) js.Value {
	r := reflect.ValueOf(in)
	return valueOf(r)
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
	t := v.Type()
	s := object.New()
	n := v.NumField()
	for i := 0; i < n; i++ {
		if f := v.Field(i); f.CanInterface() {
			k := nameOf(t.Field(i))
			s.Set(k, valueOf(f))
		}
	}
	return s
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

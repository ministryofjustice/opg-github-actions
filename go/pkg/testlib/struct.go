package testlib

import "reflect"

func StructFields(t reflect.Type) []reflect.StructField {
	return reflect.VisibleFields(t)
}

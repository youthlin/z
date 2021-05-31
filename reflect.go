package z

import (
	"fmt"
	"reflect"
)

// copyNoneZeroField copy none zero field. only exported field can be copied
func copyNoneZeroField(from, to interface{}) {
	fValue := reflect.ValueOf(from)
	tValue := reflect.ValueOf(to)
	if fValue.Type() != tValue.Type() {
		panic(fmt.Sprintf("from/to must by same type:from=%v, to=%v", fValue.Type(), tValue.Type()))
	}
	fValue = fValue.Elem()
	tValue = tValue.Elem()
	if !tValue.CanAddr() {
		panic("copy destination must be CanAddr")
	}

	for i := 0; i < fValue.NumField(); i++ {
		field := fValue.Field(i)
		if !field.IsZero() && field.CanSet() {
			tValue.Field(i).Set(field)
		}
	}
}

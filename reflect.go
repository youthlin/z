package z

import (
	"fmt"
	"reflect"
	"runtime"
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

// funcName 获取一个函数/方法的全限定名
func funName(f interface{}) (name string) {
	if f == nil {
		return "<nil>"
	}
	defer func() {
		if e := recover(); e != nil {
			name = fmt.Sprintf("<err: %v>", e)
			return
		}
	}()
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func {
		return "<not-func>"
	}
	if rf.IsNil() {
		return "<nil>"
	}
	return runtime.FuncForPC(rf.Pointer()).Name()
}

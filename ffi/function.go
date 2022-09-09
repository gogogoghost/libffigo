package ffi

import (
	"reflect"
	"unsafe"
)
import "C"

type Function struct {
	cif     *Cif
	outType reflect.Type
}

func NewFunction(cif *Cif, outType reflect.Type) *Function {
	return &Function{
		cif:     cif,
		outType: outType,
	}
}

func (self *Function) Call(args []reflect.Value) []reflect.Value {
	finalArgs := []any{}
	for _, arg := range args {
		if arg.Kind() == reflect.String {
			strPtr := C.CString(arg.String())
			defer FreePtr(unsafe.Pointer(strPtr))
		} else {
			finalArgs = append(finalArgs, arg)
		}
	}
	self.cif.Call(finalArgs...)
	return []reflect.Value{}
}

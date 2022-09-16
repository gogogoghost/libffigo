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

func (self *Function) Call(rawArgs []reflect.Value) []reflect.Value {
	args := []any{}
	for _, arg := range rawArgs {
		var r any
		switch arg.Kind() {
		case reflect.Int8:
			r = int8(arg.Int())
		case reflect.Int16:
			r = int16(arg.Int())
		case reflect.Int32:
			r = int32(arg.Int())
		case reflect.Int:
			r = int32(arg.Int())
		case reflect.Int64:
			r = int64(arg.Int())
			//uint
		case reflect.Uint8:
			r = uint8(arg.Uint())
		case reflect.Uint16:
			r = uint16(arg.Uint())
		case reflect.Uint32:
			r = uint32(arg.Uint())
		case reflect.Uint:
			r = uint32(arg.Uint())
		case reflect.Uint64:
			r = uint64(arg.Uint())
			// float
		case reflect.Float32:
			r = float32(arg.Float())
		case reflect.Float64:
			r = float64(arg.Float())
			// string
		case reflect.String:
			strPtr := C.CString(arg.String())
			defer FreePtr(unsafe.Pointer(strPtr))
			// ptr
		case reflect.Pointer, reflect.UnsafePointer:
			r = arg.Pointer()
			// default
		}
		args = append(args, r)
	}
	res := self.cif.Call(args...)
	if self.outType == nil {
		return []reflect.Value{}
	}
	var r any
	switch self.outType.Kind() {
	// int
	case reflect.Int8:
		r = res.Int8()
	case reflect.Int16:
		r = res.Int16()
	case reflect.Int32:
		r = res.Int32()
	case reflect.Int:
		r = int(res.Int32())
	case reflect.Int64:
		r = res.Int64()
		//uint
	case reflect.Uint8:
		r = res.Uint8()
	case reflect.Uint16:
		r = res.Uint16()
	case reflect.Uint32:
		r = res.Uint32()
	case reflect.Uint:
		r = uint(res.Uint32())
	case reflect.Uint64:
		r = res.Uint64()
		// float
	case reflect.Float32:
		r = res.Float()
	case reflect.Float64:
		r = res.Double()
		// string
	case reflect.String:
		r = res.String()
		// ptr
	case reflect.Pointer, reflect.UnsafePointer:
		r = res.Pointer()
		// default
	default:
		r = res
	}
	return []reflect.Value{
		reflect.ValueOf(r),
	}
}

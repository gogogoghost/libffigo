package ffi

//#cgo LDFLAGS: -lffi
//
//#include "ffi.h"
//#include "stdlib.h"
//#include "string.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var FFI_TYPE_VOID = C.ffi_type_void
var FFI_TYPE_UINT8 = C.ffi_type_uint8
var FFI_TYPE_SINT8 = C.ffi_type_sint8
var FFI_TYPE_UINT16 = C.ffi_type_uint16
var FFI_TYPE_SINT16 = C.ffi_type_sint16
var FFI_TYPE_UINT32 = C.ffi_type_uint32
var FFI_TYPE_SINT32 = C.ffi_type_sint32
var FFI_TYPE_UINT64 = C.ffi_type_uint64
var FFI_TYPE_SINT64 = C.ffi_type_sint64
var FFI_TYPE_FLOAT = C.ffi_type_float
var FFI_TYPE_DOUBLE = C.ffi_type_double
var FFI_TYPE_POINTER = C.ffi_type_pointer

var (
	emptyType = reflect.TypeOf((*interface{})(nil)).Elem()
)

type Cif struct {
	ptr  C.ffi_cif
	fPtr *[0]byte
}

func NewCif(fPtr *[0]byte, rType C.ffi_type, aTypes ...*C.ffi_type) (cif *Cif, err error) {
	cif = &Cif{}
	cif.fPtr = fPtr
	nargs := len(aTypes)
	var argsPtr **C.ffi_type
	if nargs > 0 {

		argsPtr = &aTypes[0]
	}
	ret := C.ffi_prep_cif(
		&cif.ptr,
		C.FFI_DEFAULT_ABI,
		C.uint(nargs),
		&rType,
		argsPtr,
	)
	if ret != C.FFI_OK {
		return nil, errors.New(fmt.Sprintf("prep fail:%d", ret))
	}
	return cif, nil
}

func (cif *Cif) Call(args ...any) {
	count := len(args)
	var argp *unsafe.Pointer

	fmt.Println(count)
	argsRaw := make([]unsafe.Pointer, count)
	for index, arg := range args {
		argsRaw[index] = unsafe.Pointer(&arg)
		// v := reflect.ValueOf(arg)
		// switch v.Kind() {
		// case reflect.String:
		// 	//字符串
		// 	s := C.CString(v.String())
		// 	defer C.free(unsafe.Pointer(s))
		// 	argsRaw[index] = unsafe.Pointer(s)
		// case reflect.Int:
		// 	//整数
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Int()))
		// case reflect.Int8:
		// 	//8位
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Int()))
		// case reflect.Int16:
		// 	//16位
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Int()))
		// case reflect.Int32:
		// 	//32位
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Int()))
		// case reflect.Int64:
		// 	//64位
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Int()))
		// case reflect.Uint:
		// 	//无符号
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// case reflect.Uint8:
		// 	//8位无符号
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// case reflect.Uint16:
		// 	//16位无符号
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// case reflect.Uint32:
		// 	//32位无符号
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// case reflect.Uint64:
		// 	//64位无符号
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// case reflect.Float32:
		// 	//32位浮点
		// 	argsRaw[index] = unsafe.Pointer(uintptr(math.Float32bits(float32(v.Float()))))
		// case reflect.Float64:
		// 	//64位浮点
		// 	argsRaw[index] = unsafe.Pointer(uintptr(math.Float64bits(v.Float())))
		// case reflect.Ptr:
		// 	//指针
		// 	argsRaw[index] = unsafe.Pointer(v.Pointer())
		// case reflect.Slice:
		// 	//切片
		// 	if v.Len() > 0 {
		// 		argsRaw[index] = unsafe.Pointer(v.Index(0).UnsafeAddr())
		// 	}
		// case reflect.Uintptr:
		// 	//无符号整数指针
		// 	argsRaw[index] = unsafe.Pointer(uintptr(v.Uint()))
		// default:
		// 	//其他类型 崩溃
		// 	panic(fmt.Errorf("can't bind value of type %s", v.Type()))
		// }
	}

	if count > 0 {
		argp = (*unsafe.Pointer)(&argsRaw[0])
	}

	ret := C.malloc(8)

	C.ffi_call(
		&cif.ptr,
		cif.fPtr,
		unsafe.Pointer(ret),
		argp,
	)
	fmt.Println(*(*int)(ret))
	C.free(ret)
}

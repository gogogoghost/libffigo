package ffi

//#cgo LDFLAGS: -lffi
//
//#include "ffi.h"
//#include "stdlib.h"
//#include "string.h"
//#include <stdint.h>
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

// var FFI_TYPE_VOID_PTR = &C.ffi_type_void
// var FFI_TYPE_UINT8_PTR = &C.ffi_type_uint8
// var FFI_TYPE_SINT8_PTR = &C.ffi_type_sint8
// var FFI_TYPE_UINT16_PTR = &C.ffi_type_uint16
// var FFI_TYPE_SINT16_PTR = &C.ffi_type_sint16
// var FFI_TYPE_UINT32_PTR = &C.ffi_type_uint32
// var FFI_TYPE_SINT32_PTR = &C.ffi_type_sint32
// var FFI_TYPE_UINT64_PTR = &C.ffi_type_uint64
// var FFI_TYPE_SINT64_PTR = &C.ffi_type_sint64
// var FFI_TYPE_FLOAT_PTR = &C.ffi_type_float
// var FFI_TYPE_DOUBLE_PTR = &C.ffi_type_double
// var FFI_TYPE_POINTER_PTR = &C.ffi_type_pointer

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
		argsPtr = AllocArrayOf(aTypes)
		defer FreePtr(unsafe.Pointer(argsPtr))
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
	//参数指针
	var argp *unsafe.Pointer

	if count > 0 {
		//申请一片数组空间
		arrPtr := AllocArray(count)
		defer FreePtr(arrPtr)
		//转换成指针数组
		arr := (*[1 << 30]unsafe.Pointer)(arrPtr)
		//给数组写入指对应C内存的地址
		for index, arg := range args {
			//把参数全部复制到C内存中 获取指针
			fmt.Println("参数本来的值")
			fmt.Println(arg)
			ptr := AllocValOf(arg)
			fmt.Println("这是参数拷贝后的值")
			fmt.Println((*int)(unsafe.Pointer(ptr)))
			uptr := unsafe.Pointer(ptr)
			fmt.Println("这是参数拷贝后的地址")
			PrintPtr(unsafe.Pointer(&ptr), 8)

			defer FreePtr(uptr)
			(*arr)[index] = uptr
		}
		argp = &((*arr)[0])
		PrintPtr(unsafe.Pointer(argp), 16)
	}
	fmt.Println(argp)

	var resPtrLocal uintptr
	resPtr := unsafe.Pointer(AllocValOf(resPtrLocal))
	defer FreePtr(resPtr)

	C.ffi_call(
		&cif.ptr,
		cif.fPtr,
		resPtr,
		argp,
	)
}

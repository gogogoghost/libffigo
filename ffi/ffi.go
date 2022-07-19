package ffi

//#cgo LDFLAGS: -lffi
//
//#include "ffi.h"
//#include "stdlib.h"
import "C"
import (
	"errors"
	"fmt"
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
	} else {
		argsPtr = nil
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

func (cif *Cif) Call(args ...any) any {
	var ret any
	count := len(args)
	argsRaw := make([]unsafe.Pointer, count)
	for index, arg := range args {
		argsRaw[index] = unsafe.Pointer(&arg)
	}
	argsPtr := C.malloc(C.size_t(unsafe.Sizeof(argsRaw)))
	*(*unsafe.Pointer)(argsPtr) = &argsRaw[0]
	defer C.free(argsPtr)
	C.ffi_call(
		&cif.ptr,
		cif.fPtr,
		unsafe.Pointer(&ret),
		argsPtr,
	)
	return ret
}

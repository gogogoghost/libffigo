package ffi

//#cgo LDFLAGS: -lffi
//
//#include <ffi.h>
/*
#include <stdio.h>
void ffi_call_test(ffi_cif *cif,void(*func)(void),void *result,void **args){
	// int* arg = (int*)(args[0]);
	// printf("%d\n",*arg);
	// int* res = (int*)result;
	// *res=999;
	ffi_call(cif,func,result,args);
}
void ffi_prep_cif_test(ffi_type **types){
	printf("%d\n",types[0]);
}
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
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
	ptr        *C.ffi_cif
	fPtr       unsafe.Pointer
	args_count int
}

func NewCif(fPtr unsafe.Pointer, rType C.ffi_type, aTypes ...*C.ffi_type) (cif *Cif, err error) {
	//申请空间 把cif存到C内存中
	empty_cif := C.ffi_cif{}
	cif = &Cif{
		ptr: (*C.ffi_cif)(AllocValOf(empty_cif)),
	}
	//对象销毁时释放内存
	runtime.SetFinalizer(cif, func(cif *Cif) {
		FreePtr(unsafe.Pointer(cif.ptr))
	})
	cif.fPtr = fPtr
	cif.args_count = len(aTypes)
	var argsPtr **C.ffi_type
	if cif.args_count > 0 {
		argsPtr = AllocArrayOf(aTypes)
		defer FreePtr(unsafe.Pointer(argsPtr))
	}
	C.ffi_prep_cif_test(argsPtr)
	ret := C.ffi_prep_cif(
		cif.ptr,
		C.FFI_DEFAULT_ABI,
		C.uint(cif.args_count),
		&rType,
		argsPtr,
	)

	if ret != C.FFI_OK {
		return nil, errors.New(fmt.Sprintf("prep fail:%d", ret))
	}
	return cif, nil
}

func (cif *Cif) Call(resPtr unsafe.Pointer, args ...any) {
	if len(args) != cif.args_count {
		panic("Wrong args count")
	}

	argp := AllocParams(args)
	defer FreeParams(argp)

	// fmt.Println(C.call_test(

	// 	(*[0]byte)(cif.fPtr),
	// 	argp,
	// ))
	// C.test(argp)
	C.ffi_call_test(
		cif.ptr,
		(*[0]byte)(cif.fPtr),
		resPtr,
		argp,
	)
}

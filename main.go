package main

//#cgo LDFLAGS: -lffi
//
/*
#include<ffi.h>
int add(void **src){
	return *(int*)(src[0])+1;
}
void call_abs(
	ffi_type *res_type,
	ffi_type **arg_type,
	int count,
	void* func,
	void* res,
	void** values
){
	ffi_cif cif;
	// ffi_type *arg_types[1];
	// void *arg_values[1];
	ffi_status status;

	// arg_types[0] = &ffi_type_sint32;

	if ((status = ffi_prep_cif(&cif, FFI_DEFAULT_ABI,
		1, res_type, arg_type)) != FFI_OK)
	{
		// Handle the ffi_status error.
		return;
	}

	// Invoke the function.
	ffi_call(&cif, FFI_FN(func), res, values);
}
*/
import "C"
import (
	"fmt"

	"github.com/gogogoghost/libffigo/dl"
	"github.com/gogogoghost/libffigo/ffi"
)

func callAdd(args ...any) {
	argp := ffi.AllocParams(args)
	defer ffi.FreeParams(argp)
	fmt.Println(
		C.add(
			argp,
		),
	)
}

func main() {

	// var src any = 99
	// srcPtr := unsafe.Pointer(ffi.AllocValOf(src))
	// arrPtr := ffi.AllocArray(1)
	// arr := (*[1 << 10]unsafe.Pointer)(arrPtr)
	// (*arr)[0] = srcPtr
	// callAdd(-101)

	lib, err := dl.Open("libc.so.6", dl.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	f, err := lib.Sym("abs")
	if err != nil {
		panic(err)
	}

	cfi, err := ffi.NewCif(
		f,
		ffi.FFI_TYPE_SINT32,
		&ffi.FFI_TYPE_SINT32,
	)
	if err != nil {
		panic(err)
	}
	resPtr := ffi.Alloc(8)
	ffi.PrintPtr(resPtr, 8)
	defer ffi.FreePtr(resPtr)
	cfi.Call(resPtr, -100)

	ffi.PrintPtr(resPtr, 8)
	fmt.Println(*(*int)(resPtr))
}

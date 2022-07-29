package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/gogogoghost/libffigo/dl"
	"github.com/gogogoghost/libffigo/ffi"
)

func TestAny(*testing.T) {
	a := 123456789
	var b any = a
	// fmt.Println(unsafe.Pointer(&a))
	ptr := &a
	ffi.PrintPtr(unsafe.Pointer(ptr), 8)
	fmt.Println("变量地址==========")
	ffi.PrintPtr(unsafe.Pointer(&ptr), 8)
	fmt.Println((*int)(ffi.AllocValOf(b)))
}

func BenchmarkFFI(*testing.B) {
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
	cfi.Call(-100)
}

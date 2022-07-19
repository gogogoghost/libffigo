package main

import "C"

import (
	"fmt"

	"github.com/gogogoghost/libffigo/dl"
	"github.com/gogogoghost/libffigo/ffi"
)

func main() {
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
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(
		cfi.Call(C.int(-10)),
	)
}

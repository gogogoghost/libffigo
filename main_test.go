package main

import (
	"fmt"
	"testing"
	"unsafe"

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
	fmt.Println(*(*int)(ffi.AllocValOf(b)))
}

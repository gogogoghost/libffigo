package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/gogogoghost/libffigo/ffi"
)

func TestSize(t *testing.T) {
	a := 10
	fmt.Println(unsafe.Sizeof(a))
	ptr := ffi.PtrGo2C(a)
	fmt.Println(*(*int)(ptr))
}

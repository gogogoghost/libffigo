package test

import (
	"fmt"
	"reflect"
	"testing"

	ffi "github.com/gogogoghost/libffigo"
)

func TestAbsRaw(t *testing.T) {
	lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	f := lib.SymRawMust("abs", ffi.SINT32, ffi.SINT32)
	fmt.Println(f.Call(-100).Int())
}

func TestAbs(t *testing.T) {
	lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	var abs func(int) int
	lib.SymMust("abs", &abs, ffi.SINT32, ffi.SINT32)
	fmt.Println(abs(-100))
}

func TestType(t *testing.T) {
	var out reflect.Type
	fmt.Println(out)
}

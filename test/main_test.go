package main_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gogogoghost/libffigo/ffi"
)

func TestAny(*testing.T) {
	lib, err := ffi.Open("libudev.so", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	Udev_new := lib.SymMust("udev_new", ffi.PTR)
	Udev_enumerate_new := lib.SymMust("udev_enumerate_new", ffi.PTR, ffi.PTR)
	Udev_enumerate_scan_devices := lib.SymMust("udev_enumerate_scan_devices", ffi.SINT64, ffi.PTR)

	ctx := Udev_new.Call().Pointer()
	println(ctx)
	println("start call=============")
	enumer := Udev_enumerate_new.Call(ctx).Pointer()
	println(enumer)
	println("start call2=============")
	Udev_enumerate_scan_devices.Call(enumer)
}

func TestAbs(t *testing.T) {
	lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	f := lib.SymMust("abs", ffi.SINT32, ffi.SINT32)
	res := f.Call(-100)
	fmt.Println(res.Int32())
}

func TestType(t *testing.T) {
	var out reflect.Type
	fmt.Println(out)
}

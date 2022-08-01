package main_test

import (
	"fmt"
	"testing"

	"github.com/gogogoghost/libffigo/ffi"
)

func TestAny(*testing.T) {
	// lib, err := ffi.Open("libudev.so", ffi.RTLD_LAZY)
	// if err != nil {
	// 	panic(err)
	// }
	// Udev_new, err := lib.Sym("udev_new", ffi.FFI_TYPE_POINTER)
	// if err != nil {
	// 	panic(err)
	// }
	// Udev_enumerate_new, err := lib.Sym("udev_enumerate_new", ffi.FFI_TYPE_POINTER, ffi.FFI_TYPE_POINTER)
	// if err != nil {
	// 	panic(err)
	// }
	// Udev_enumerate_scan_devices, err := lib.Sym("udev_enumerate_scan_devices", ffi.FFI_TYPE_VOID, ffi.FFI_TYPE_POINTER)

	// ctx := Udev_new.Call()
	// enumer := Udev_enumerate_new.Call(ctx)
	// Udev_enumerate_scan_devices.Call(enumer)
}

func TestAbs(t *testing.T) {
	lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	f, err := lib.Sym("abs", &ffi.TYPE_SINT32, &ffi.TYPE_SINT32)
	if err != nil {
		panic(err)
	}
	resPtr := f.Call(-100)
	fmt.Println(*(*int)(resPtr))
}

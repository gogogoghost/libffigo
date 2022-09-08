package main_test

import (
	"fmt"
	"testing"
	"unsafe"

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
	enumer := Udev_enumerate_new.Call(ctx).Pointer()
	println(enumer)
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

func TestPtr(t *testing.T) {
	num := 99999999999
	numPtr := &num
	tmpArr := (*[1 << 30]byte)(unsafe.Pointer(numPtr))
	arr := make([]byte, ffi.PtrSize)
	for i := 0; i < int(ffi.PtrSize); i++ {
		arr[i] = tmpArr[i]
	}
	newPtr := (*int)(unsafe.Pointer(&arr[0]))
	println(*newPtr)
}

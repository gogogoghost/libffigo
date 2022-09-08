package main_test

import (
	"fmt"
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

func TestPtr(t *testing.T) {
	num := 64

	numPtr := &num
	println(numPtr)
	numPtr2 := ffi.AllocValOf(numPtr)
	println(*(**int)(numPtr2))
	println(*(*(**int)(numPtr2)))
}

func TestMyLib(t *testing.T) {
	lib, err := ffi.Open("/home/ghost/tmp/libtest.so", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	num := 64
	numOnC := ffi.AllocValOf(num)
	f1 := lib.SymMust("addOne", ffi.SINT32, ffi.SINT32)
	fmt.Println(f1.Call(num).Int32())
	f := lib.SymMust("getNum", ffi.SINT32, ffi.PTR)
	res := f.Call(numOnC)
	fmt.Println(res.Int32())
	f2 := lib.SymMust("getStr", ffi.PTR)
	fmt.Println(f2.Call().String())
}

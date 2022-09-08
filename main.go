package main

import "github.com/gogogoghost/libffigo/ffi"

func main() {
	lib, err := ffi.Open("libudev.so", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	Udev_new := lib.SymMust("udev_new", ffi.PTR)
	Udev_enumerate_new := lib.SymMust("udev_enumerate_new", ffi.PTR, ffi.PTR)
	Udev_enumerate_scan_devices := lib.SymMust("udev_enumerate_scan_devices", ffi.SINT32, ffi.PTR)

	ctx := Udev_new.Call().Pointer()
	println(ctx)
	enumer := Udev_enumerate_new.Call(ctx).Pointer()
	println(enumer)
	Udev_enumerate_scan_devices.Call(enumer)
}

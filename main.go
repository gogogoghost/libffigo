package main

/*
int add(void **src){
	return *(int*)(src[0])+1;
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/gogogoghost/libffigo/ffi"
)

func main() {

	var src any = 99
	srcPtr := unsafe.Pointer(ffi.AllocValOf(src))
	arrPtr := ffi.AllocArray(1)
	arr := (*[1 << 10]unsafe.Pointer)(arrPtr)
	(*arr)[0] = srcPtr
	fmt.Println(
		C.add(&(*arr)[0]),
	)
	// lib, err := dl.Open("libc.so.6", dl.RTLD_LAZY)
	// if err != nil {
	// 	panic(err)
	// }
	// f, err := lib.Sym("abs")
	// if err != nil {
	// 	panic(err)
	// }
	// cfi, err := ffi.NewCif(
	// 	f,
	// 	ffi.FFI_TYPE_SINT32,
	// 	&ffi.FFI_TYPE_SINT32,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// cfi.Call(-100)
}

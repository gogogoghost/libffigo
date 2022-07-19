package ffi

//#include "stdlib.h"
//#include "string.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func PtrGo2C(src any) unsafe.Pointer {

	size := unsafe.Sizeof(reflect.TypeOf(src).Size())
	fmt.Println(size)
	ptr := C.malloc(C.size_t(size))
	C.memcpy(ptr, unsafe.Pointer(&src), C.size_t(size))
	return unsafe.Pointer(ptr)
}

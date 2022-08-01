package main

import (
	"fmt"

	"github.com/gogogoghost/libffigo/ffi"
)

func work() {

}

func main() {
	//加载libc库
	lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
	if err != nil {
		panic(err)
	}
	//查找函数
	f, err := lib.Sym("abs", ffi.FFI_TYPE_SINT32, &ffi.FFI_TYPE_SINT32)
	if err != nil {
		panic(err)
	}
	//申请8字节内存
	resPtr := ffi.Alloc(8)
	defer ffi.FreePtr(resPtr)
	//调用函数
	f.Call(resPtr, -10086)
	fmt.Println(*(*int)(resPtr))
}

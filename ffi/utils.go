package ffi

//#include <stdlib.h>
import "C"
import (
	"encoding/hex"
	"fmt"
	"reflect"
	"unsafe"
)

// C语言空指针
const NilPtr uintptr = 0

// 一个指针长度
var PtrSize = unsafe.Sizeof(NilPtr)

// any类型
type AnyStruct struct {
	typePtr uintptr
	dataPtr uintptr
}

// []T of go -> *T of C 仅支持转换指针
func AllocArrayOf[T any](src []T) *T {
	length := len(src)
	ptr := AllocArray(length)
	//转换成数组
	arr := (*[1 << 30]T)(ptr)
	//把数组内容拷贝过去
	for index, value := range src {
		arr[index] = value
	}
	return &arr[0]
}

// 申请一个size大小的C指针数组空间 什么都不拷贝
func AllocArray(size int) unsafe.Pointer {
	//生成一个数组 长度指针字节*长度
	ptr := C.malloc(C.size_t(int(PtrSize) * (size + 1)))
	//数组最后一位填0
	arr := (*[1 << 30]uint64)(ptr)
	(*arr)[size] = 0
	return ptr
}

// 从指针部分取出内容并生成新指针
func GetPtrFromAny(ptr *any) unsafe.Pointer {
	anyPtr := (*AnyStruct)(unsafe.Pointer(ptr))
	return unsafe.Pointer(anyPtr.dataPtr)
}

// 转换一个指针
func AllocValOf(src any) unsafe.Pointer {
	//获取实际指向的数据
	dataPtr := GetPtrFromAny(&src)
	//获取反射
	val := reflect.ValueOf(src)
	//获取src实际大小
	realSize := val.Type().Size()
	//申请空间
	destPtr := C.malloc(C.size_t(realSize))
	destArr := (*[1 << 30]byte)(destPtr)
	//获取src的指针 转换成数组
	srcArr := (*[1 << 30]byte)(dataPtr)
	//按字节加上偏移量拷贝
	copy(destArr[0:realSize], srcArr[0:realSize])
	return destPtr
}

// 申请一片指定大小的空间
func Alloc(size int) unsafe.Pointer {
	ptr := C.malloc(C.size_t(size))
	return ptr
}

// 释放一个指针
func FreePtr(ptr unsafe.Pointer) {
	C.free(ptr)
}

// 把一个指针指向的东西按字节数量
func PrintPtr(ptr unsafe.Pointer, size int) {
	arr := (*[1 << 30]byte)(ptr)
	var buf = []byte{}
	for i := 0; i < size; i++ {
		buf = append(buf, (*arr)[i])
	}
	fmt.Println(hex.EncodeToString(buf))
}

// 把any[] 转void*
func AllocParams(args []any) *unsafe.Pointer {
	count := len(args)
	var argp *unsafe.Pointer
	//申请一片数组空间
	arrPtr := AllocArray(count)
	// defer FreePtr(arrPtr)
	//转换成指针数组
	arr := (*[1 << 30]unsafe.Pointer)(arrPtr)
	//给数组写入指对应C内存的地址
	for index, arg := range args {
		ptr := AllocValOf(arg)
		// defer FreePtr(ptr)
		(*arr)[index] = ptr
	}
	argp = &((*arr)[0])
	return argp
}

// 释放void** 并将数组内的所有内存释放
func FreeParams(ptr *unsafe.Pointer) {
	arrPtr := unsafe.Pointer(ptr)
	arr := (*[1 << 30]unsafe.Pointer)(arrPtr)
	for _, p := range arr {
		if uintptr(p) == 0 {
			break
		}
		C.free(p)
	}
	C.free(arrPtr)
}

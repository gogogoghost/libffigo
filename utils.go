package ffi

//#include <stdlib.h>
import "C"
import (
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
	//转换成对应类型的数组
	arr := Ptr2Arr[T](ptr, length)
	//把数组内容拷贝过去
	copy(arr[:], src)
	return &arr[0]
}

// 申请一个size大小的C指针数组空间 什么都不拷贝
func AllocArray(size int) unsafe.Pointer {
	//生成一个数组 长度指针字节*长度
	ptr := C.malloc(C.size_t(int(PtrSize) * (size + 1)))
	//数组最后一位填0
	arr := Ptr2Arr[uintptr](ptr, size)
	arr[size] = uintptr(0)
	return ptr
}

// 从指针部分取出内容并生成新指针
func GetPtrFromAny(ptr *any) unsafe.Pointer {
	anyPtr := (*AnyStruct)(unsafe.Pointer(ptr))
	return unsafe.Pointer(anyPtr.dataPtr)
}

// 转换一个指针 传入any 整个复制到C去 返回指针
func AllocValOf(src any) unsafe.Pointer {
	//获取实际指向的数据
	dataPtr := GetPtrFromAny(&src)
	//获取反射
	val := reflect.ValueOf(src)
	if val.Kind() == reflect.UnsafePointer || val.Kind() == reflect.Pointer {
		//如果类型为指针，dataPtr取出来其实就是指针内容
		ptrValue := uintptr(dataPtr)
		dataPtr = unsafe.Pointer(&ptrValue)
	}
	//获取src实际大小
	realSize := val.Type().Size()
	//申请空间
	destPtr := C.malloc(C.size_t(realSize))
	destArr := Ptr2Arr[byte](destPtr, int(realSize))
	//获取src的指针 转换成数组
	srcArr := Ptr2Arr[byte](dataPtr, int(realSize))
	//按字节加上偏移量拷贝
	/*
		由于数据本来不是数组，转成数组好拷贝，所以内存布局并没有末尾的0
		必须强制指定要拷贝的长度，否则容易拷贝越界
	*/
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

// 把any[] 转void*
func AllocParams(args []any) *unsafe.Pointer {
	count := len(args)
	var argp *unsafe.Pointer
	//申请一片数组空间
	arrPtr := AllocArray(count)
	//转换成指针数组
	arr := Ptr2Arr[unsafe.Pointer](arrPtr, count)
	//给数组写入指对应C内存的地址
	for index, arg := range args {
		// 给每个变量单独申请空间
		ptr := AllocValOf(arg)
		arr[index] = ptr
	}
	argp = &(arr[0])
	return argp
}

// 释放void** 并将数组内的所有内存释放
func FreeParams(ptr *unsafe.Pointer) {
	arrPtr := unsafe.Pointer(ptr)
	ptrAddr := uintptr(arrPtr)
	for {
		//取出指针指向数据
		dataPtr := *(*unsafe.Pointer)(unsafe.Pointer(ptrAddr))
		//0表示数组末尾了
		if uintptr(dataPtr) == 0 {
			break
		}
		C.free(dataPtr)
		ptrAddr += PtrSize
	}
	C.free(arrPtr)
}

// ptr->[]T
func Ptr2Arr[T any](ptr unsafe.Pointer, length int) []T {
	sliceHeader := struct {
		p   unsafe.Pointer
		len int
		cap int
	}{ptr, length + 1, length + 1}
	return *(*[]T)(unsafe.Pointer(&sliceHeader))
}

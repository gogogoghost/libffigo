package ffi

//#cgo LDFLAGS: -lffi -ldl
//
//#include <dlfcn.h>
//#include <stdlib.h>
//#include <ffi.h>
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

//所有FFI参数类型
type FFI_TYPE struct {
	typePtr *C.ffi_type
	size    int
}

//变量类型及长度
var (
	TYPE_VOID = FFI_TYPE{
		typePtr: &C.ffi_type_void,
		size:    0,
	}
	TYPE_UINT8 = FFI_TYPE{
		typePtr: &C.ffi_type_uint8,
		size:    1,
	}
	TYPE_SINT8 = FFI_TYPE{
		typePtr: &C.ffi_type_sint8,
		size:    1,
	}
	TYPE_UINT16 = FFI_TYPE{
		typePtr: &C.ffi_type_uint16,
		size:    2,
	}
	TYPE_SINT16 = FFI_TYPE{
		typePtr: &C.ffi_type_sint16,
		size:    2,
	}
	TYPE_UINT32 = FFI_TYPE{
		typePtr: &C.ffi_type_uint32,
		size:    4,
	}
	TYPE_SINT32 = FFI_TYPE{
		typePtr: &C.ffi_type_sint32,
		size:    4,
	}
	TYPE_UINT64 = FFI_TYPE{
		typePtr: &C.ffi_type_uint64,
		size:    8,
	}
	TYPE_SINT64 = FFI_TYPE{
		typePtr: &C.ffi_type_sint64,
		size:    8,
	}
	TYPE_FLOAT = FFI_TYPE{
		typePtr: &C.ffi_type_float,
		size:    4,
	}
	TYPE_DOUBLE = FFI_TYPE{
		typePtr: &C.ffi_type_double,
		size:    8,
	}
	TYPE_POINTER = FFI_TYPE{
		typePtr: &C.ffi_type_pointer,
		size:    int(PtrSize),
	}
)

//dlopen flag
const (
	RTLD_LAZY     = int(C.RTLD_LAZY)
	RTLD_NOW      = int(C.RTLD_NOW)
	RTLD_GLOBAL   = int(C.RTLD_GLOBAL)
	RTLD_LOCAL    = int(C.RTLD_LOCAL)
	RTLD_NODELETE = int(C.RTLD_NODELETE)
	RTLD_NOLOAD   = int(C.RTLD_NOLOAD)
)

//描述一个dlopen 的 library
type Lib struct {
	ptr unsafe.Pointer
}

//描述一个Cif
type Cif struct {
	ptr       *C.ffi_cif
	fPtr      unsafe.Pointer
	argsCount int
	resType   *FFI_TYPE
}

//构造一个cif
func NewCif(fPtr unsafe.Pointer, rType *FFI_TYPE, aTypes ...*FFI_TYPE) (cif *Cif, err error) {
	//申请空间 把cif存到C内存中
	empty_cif := C.ffi_cif{}
	cif = &Cif{
		ptr: (*C.ffi_cif)(AllocValOf(empty_cif)),
	}
	cif.fPtr = fPtr
	cif.argsCount = len(aTypes)
	var argsPtr **C.ffi_type
	if cif.argsCount > 0 {
		//这片参数空间在对象销毁时释放
		typesArr := make([]*C.ffi_type, cif.argsCount)
		for index, aType := range aTypes {
			typesArr[index] = aType.typePtr
		}
		argsPtr = AllocArrayOf(typesArr)
	}
	//对象销毁时释放内存
	runtime.SetFinalizer(cif, func(cif *Cif) {
		//销毁cif内存
		FreePtr(unsafe.Pointer(cif.ptr))
		//销毁参数数组内存
		if argsPtr != nil {
			FreePtr(unsafe.Pointer(argsPtr))
		}
	})
	//保存一下返回类型
	cif.resType = rType
	ret := C.ffi_prep_cif(
		cif.ptr,
		C.FFI_DEFAULT_ABI,
		C.uint(cif.argsCount),
		rType.typePtr,
		argsPtr,
	)
	if ret != C.FFI_OK {
		return nil, errors.New(fmt.Sprintf("prep fail:%d", ret))
	}
	return cif, nil
}

//调用函数
func (cif *Cif) Call(args ...any) unsafe.Pointer {
	if len(args) != cif.argsCount {
		panic("Wrong args count")
	}
	//内存复制到C空间
	argp := AllocParams(args)
	defer FreeParams(argp)
	//返回类型默认为nil
	var resPtr unsafe.Pointer
	//返回类型不为void时，申请内存
	resSize := cif.resType.size
	//go空间的字节数组 用于拷贝返回数据到此
	var resArr []byte
	if resSize > 0 {
		//申请用于存放临时返回值的空间
		resPtr = Alloc(resSize)
		defer FreePtr(resPtr)
		//再申请一片go空间 存放拷贝后的数据
		resArr = make([]byte, resSize)
	}
	//发起调用
	C.ffi_call(
		cif.ptr,
		(*[0]byte)(cif.fPtr),
		resPtr,
		argp,
	)
	if resSize > 0 {
		//返回复制后的地址
		tmpArr := (*[1 << 30]byte)(resPtr)
		copy(resArr[:], (*tmpArr)[0:resSize])
		return unsafe.Pointer(&resArr[0])
	} else {
		//返回nil
		return resPtr
	}
}

//获取dl错误
func dlerror() error {
	s := C.dlerror()
	return errors.New(C.GoString(s))
}

//dlopen
func Open(name string, flag int) (lib *Lib, err error) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	ptr := C.dlopen(str, C.int(flag))
	if ptr == nil {
		return nil, dlerror()
	}
	return &Lib{
		ptr: ptr,
	}, nil
}

//dlsym
func (lib *Lib) Sym(name string, rType *FFI_TYPE, aTypes ...*FFI_TYPE) (*Cif, error) {
	//查找函数指针
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	ptr := C.dlsym(lib.ptr, str)
	if ptr == nil {
		return nil, dlerror()
	}
	//将函数指针使用ffi初始化
	cif, err := NewCif(ptr, rType, aTypes...)
	if err != nil {
		return nil, err
	}
	return cif, nil
}

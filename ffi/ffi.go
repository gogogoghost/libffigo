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
var (
	FFI_TYPE_VOID    = C.ffi_type_void
	FFI_TYPE_UINT8   = C.ffi_type_uint8
	FFI_TYPE_SINT8   = C.ffi_type_sint8
	FFI_TYPE_UINT16  = C.ffi_type_uint16
	FFI_TYPE_SINT16  = C.ffi_type_sint16
	FFI_TYPE_UINT32  = C.ffi_type_uint32
	FFI_TYPE_SINT32  = C.ffi_type_sint32
	FFI_TYPE_UINT64  = C.ffi_type_uint64
	FFI_TYPE_SINT64  = C.ffi_type_sint64
	FFI_TYPE_FLOAT   = C.ffi_type_float
	FFI_TYPE_DOUBLE  = C.ffi_type_double
	FFI_TYPE_POINTER = C.ffi_type_pointer
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
	ptr        *C.ffi_cif
	fPtr       unsafe.Pointer
	args_count int
}

func NewCif(fPtr unsafe.Pointer, rType C.ffi_type, aTypes ...*C.ffi_type) (cif *Cif, err error) {
	//申请空间 把cif存到C内存中
	empty_cif := C.ffi_cif{}
	cif = &Cif{
		ptr: (*C.ffi_cif)(AllocValOf(empty_cif)),
	}
	cif.fPtr = fPtr
	cif.args_count = len(aTypes)
	var argsPtr **C.ffi_type
	if cif.args_count > 0 {
		//这片参数空间在对象销毁时释放
		argsPtr = AllocArrayOf(aTypes)
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
	ret := C.ffi_prep_cif(
		cif.ptr,
		C.FFI_DEFAULT_ABI,
		C.uint(cif.args_count),
		&rType,
		argsPtr,
	)
	if ret != C.FFI_OK {
		return nil, errors.New(fmt.Sprintf("prep fail:%d", ret))
	}
	return cif, nil
}

func (cif *Cif) Call(resPtr unsafe.Pointer, args ...any) {
	if len(args) != cif.args_count {
		panic("Wrong args count")
	}

	argp := AllocParams(args)
	defer FreeParams(argp)

	C.ffi_call(
		cif.ptr,
		(*[0]byte)(cif.fPtr),
		resPtr,
		argp,
	)
}

func dlerror() error {
	s := C.dlerror()
	return errors.New(C.GoString(s))
}

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

func (lib *Lib) Sym(name string, rType C.ffi_type, aTypes ...*C.ffi_type) (*Cif, error) {
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

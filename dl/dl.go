package dl

import (
	"errors"
	"unsafe"
)

//#cgo LDFLAGS: -ldl
//
//#include <dlfcn.h>
//#include <stdlib.h>
import "C"

const (
	RTLD_LAZY     = int(C.RTLD_LAZY)
	RTLD_NOW      = int(C.RTLD_NOW)
	RTLD_GLOBAL   = int(C.RTLD_GLOBAL)
	RTLD_LOCAL    = int(C.RTLD_LOCAL)
	RTLD_NODELETE = int(C.RTLD_NODELETE)
	RTLD_NOLOAD   = int(C.RTLD_NOLOAD)
)

type Lib struct {
	ptr unsafe.Pointer
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

func (lib *Lib) Sym(name string) (unsafe.Pointer, error) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	ptr := C.dlsym(lib.ptr, str)
	if ptr == nil {
		return nil, dlerror()
	}
	return ptr, nil
}

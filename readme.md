## libffigo

dynamic loading library with libffi support go 1.19

### usage

get module

```sh
go get github.com/gogogoghost/libffigo
```

```go
// load library
lib, err := ffi.Open("libc.so.6", ffi.RTLD_LAZY)
if err != nil {
    panic(err)
}
// declare a same function with library
var abs func(int) int
// find it
lib.SymMust(
    //function name
    "abs",
    //local funtion pointer
    &abs,
    //return type
    ffi.SINT32,
    //parameters type
    ffi.SINT32,
)
//use it
fmt.Println(abs(-100))
```
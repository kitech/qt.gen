package main

/*
#include <stdio.h>

int print_go_array(unsigned char *a0, int len)
{
    int res = 0;
    for (int i = 0; i < len; i++) {
        printf("e%d=%d\n", i, a0[i]);
        res += a0[i];
    }
    printf("s=%s\n", a0);
    return res;
}

*/
import "C"
import "unsafe"
import "reflect"
import "fmt"

// 一维数组到指针的转换
func Byte2Charp(a0 interface{}) *C.uchar {
	var p0 = a0.([]byte)

	// https://coderwall.com/p/m_ma7q/pass-go-slices-as-c-array-parameters
	Byte2CharpLong := func(a0 interface{}) *C.uchar {
		var p0 = a0.([]byte)
		var plen = len(p0)
		ref := reflect.ValueOf(p0)
		ty := reflect.TypeOf(p0)
		// var c int = 1
		addr := uint64(ref.Pointer()) // + (ty.Align() * c)

		if false {
			fmt.Println(plen)
			fmt.Printf("the slice addr: %p\n", &p0)
			fmt.Printf("the first element of the underlying array: %p\n", &p0[0])
			fmt.Println("type []byte align:", ty.Align())
			fmt.Printf("Addr of the underlying data: %d\n", addr)
		}

		retp0 := (*C.uchar)((unsafe.Pointer)(uintptr(addr)))
		// fmt.Println(retp0)
		return retp0
	}

	retp0 := Byte2CharpLong(a0)
	retp1 := (*C.uchar)((unsafe.Pointer)(reflect.ValueOf(p0).Pointer())) // OK
	if retp0 != retp1 {
		panic("should equal")
	}
	return retp1
	// return (*C.uchar)((unsafe.Pointer)(reflect.ValueOf(p0).UnsafeAddr())) // crash
}

func main() {
	var ba = make([]byte, 3)
	ba[0] = 'a'
	ba[1] = 'b'
	ba[2] = 'e'
	ba = []byte{'e', 'f', 'g', 'y'}
	ba = []byte("iop")

	sum := C.print_go_array(Byte2Charp(ba), C.int(len(ba)))
	fmt.Println(ba, sum)
}

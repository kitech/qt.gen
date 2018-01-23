package main

/*
#include <stdio.h>

float print_go_array(float **a0, int len)
{
    float res = 0;
    for (int i = 0; i < len; i++) {
        printf("e%d=%f\n", i, a0[1][i]);
        res += a0[0][i];
    }
    // printf("s=%s\n", a0);
    return res;
}

float print_go_array1(float *a0, int len)
{
    float res = 0;
    for (int i = 0; i < len; i++) {
        for (int j = 0; j < len; j++) {
            printf("e%d=%f\n", i*len + j, a0[i*len + j]);
            res += a0[len*i + j];
        }
    }
    // printf("s=%s\n", a0);
    return res;
}

*/
import "C"
import "unsafe"
import "reflect"
import "fmt"

// 二维数组到二级指针的转换
func DDArr2PP(a0 interface{}) **C.float {
	var p0 = a0.([][]float32)

	pp0 := []unsafe.Pointer{}
	for i := 0; i < len(p0); i++ {
		pp0 = append(pp0, unsafe.Pointer(reflect.ValueOf(p0[i]).Pointer()))
	}

	retp1 := (**C.float)((unsafe.Pointer)(reflect.ValueOf(pp0).Pointer())) // OK
	return retp1
}

func DDArr2PP1(a0 interface{}) *C.float {
	var p0 = a0.([]float32)

	// https://coderwall.com/p/m_ma7q/pass-go-slices-as-c-array-parameters
	Byte2CharpLong := func(a0 interface{}) *C.float {
		var p0 = a0.([]float32)
		var plen = len(p0)
		ref := reflect.ValueOf(p0)
		ty := reflect.TypeOf(p0)
		// var c int = 1
		addr := uint64(ref.Pointer()) // + (ty.Align() * c)

		if true {
			fmt.Printf("slice len: %d\n", plen)
			fmt.Printf("the slice addr: %p\n", &p0)
			fmt.Printf("the first element of the underlying array: %p\n", &p0[0])
			fmt.Println("type []float32 align:", ty.Align())
			fmt.Printf("Addr of the underlying data: %d\n", addr)
		}

		retp0 := (*C.float)((unsafe.Pointer)(uintptr(addr)))
		// fmt.Println(retp0)
		return retp0
	}

	retp0 := Byte2CharpLong(a0)
	retp1 := (*C.float)((unsafe.Pointer)(reflect.ValueOf(p0).Pointer())) // OK
	if retp0 != retp1 {
		panic("should equal")
	}
	return retp1
	// return (*C.uchar)((unsafe.Pointer)(reflect.ValueOf(p0).UnsafeAddr())) // crash
}

// https://golang.org/doc/effective_go.html  2D array
func main() {
	var ba1 = make([]float32, 3*3)
	ba1 = []float32{3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 1.0, 1.1, 1.2}
	sum1 := C.print_go_array1(DDArr2PP1(ba1), C.int(len(ba1)/3))
	fmt.Println(ba1, sum1)

	var ba = make([][]float32, 3)
	ba[0] = append(ba[0], 3.4)
	ba[0] = append(ba[0], 3.5)
	ba[0] = append(ba[0], 3.6)
	// 不再是连续内存了。
	ba[1] = append(ba[1], 3.7)
	ba[1] = append(ba[1], 3.8)
	ba[1] = append(ba[1], 3.9)
	// 不再是连续内存了。
	ba[2] = append(ba[1], 1.0)
	ba[2] = append(ba[1], 1.1)
	ba[2] = append(ba[1], 1.2)
	fmt.Println(ba, ba[0], len(ba[0]))

	sum := C.print_go_array(DDArr2PP(ba), C.int(len(ba)))
	fmt.Println(ba, sum)

}

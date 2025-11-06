package main

// dont modify gopp


import "unsafe"

// // begin cgo
//
// /*
// #include <stdint.h>
// */
// import "C"
//
// end cgo

type i32 = int32
type i64 = int64
type u32 = uint32
type u64 = uint64
type f32 = float32
type f64 = float64
type sizet = int
type isize = int
type usize = uintptr
type voidptr = unsafe.Pointer
type u128 = struct{ H, L u64 }
type i128 = struct{ H, L i64 }
type fatptr = struct{ H, L usize }
type quadptr = struct{ H0, H1, L0, L1 usize }

// c&go conflict???
// type byteptr = *byte
// type charptr = *uint8

// // begin cgo
// type cuptr = C.uintptr_t
// type csizet = C.size_t
// type cvptr = *C.void
// type charptr = *C.char
// type ucharptr = *C.uchar
// type byteptr = *C.uchar
// type scharptr = *C.schar
// type wcharptr = *C.uint16_t
// type cbool = C.int  // = go.int32
// type cint = C.int   // = go.int32
// type cuint = C.uint // = go.uint32
// type cshort = C.short
// type clong = C.long
// type culong = C.ulong
// type clonglong = C.longlong
// type culonglong = C.ulonglong
// type cfloat = C.float
// type cdouble = C.double
// type cuintptr = C.uintptr_t
// type ci64 = C.int64_t
// type cu64 = C.uint64_t
// type ci32 = C.int32_t
// type cu32 = C.uint32_t
// type ci16 = C.int16_t
// type cu16 = C.uint16_t
//
// end cgo

func anyptr2uptr[T any](p *T) usize {
	var pp = usize(voidptr(p))
	return pp
}

// // begin cgo
//
// func anyptr2uptrc[T any](p *T) cuptr {
// 	var pp = uintptr(voidptr(p))
// 	return cuptr(pp)
// }
//
// end cgo

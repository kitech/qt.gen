package main

import (
	"fmt"
	"gopp"
	"log"
	"reflect"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

func init() {
	if false {
		reflect.TypeOf("heheh")
	}
}

// 需要考虑的目标类型转换，还是挺多的
// 转换的源类型为CPP类型
const (
	// 在go源代码中使用的类型转换
	ArgDesc_GO_SIGNATURE = iota + 8
	AsGoSignature        // 转换到go函数签名相应的类型
	AsGoITF              // 转换到go函数签名中需要用到interface的相应的类型
	AsGoReturn           // 转换 go函数返回值相应的类型，可能与签名中的不一样
	ArgTyDesc_GO_INVOKE_GO
	PrmTyDesc_GO_INVOKE_GO  // 有时需要做*或者&操作或者强制类型转换
	AsGoCallGoFConv         // 同上, FConv指的是force convert, reference/dereference
	AsGoCallCGOFConv        // 当在一个纯go的函数中调用一个接收cgo类型参数时使用
	PrmTyDesc_GO_INVOKE_CGO // 同上

	// 在cgo源代码中使用的类型转换
	AsCGOSignature // 转换到cgo函数签名相应的类型，指的是//export的go函数，需要接收cgo数据类型
	AsCGOReturn    // 同上
	ArgDesc_CGO_SIGNATURE
	PrmTyDesc_CGO_INVOKE_CGO
	ArgTyDesc_CGO_SIGNATURE
	ArgTyDesc_CGO_INVOKE_CGO
	ArgTyDesc_CGO_INVOKE_GO
	RetTyDesc_CGO
	PrmTyDesc_CGO_INVOKE_GO

	// 在cpp源代码中使用的类型转换
	AsCPPSignature
	AsCPPReturn
	ArgDesc_CPP_SIGNATURE
	PrmTyDesc_CPP_INVOKE_CPP
	ArgTyDesc_CPP_SIGNAUTE
	ArgTyDesc_CPP_INVOKE_CPP
	RetTyDesc_CPP
	PrmTyDesc_CPP_INVOKE_C

	// 在c源代码中使用的类型转换
	AsCSignature
	AsCReturn
	ArgDesc_C_SIGNATURE
	PrmTyDesc_C_INVOKE_C
	ArgTyDesc_C_SIGNATURE
	ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN
	ArgTyDesc_C_INVOKE_C
	RetTyDesc_C
	PrmTyDesc_C_INVOKE_CPP

	// 在ffi源代码中使用的类型转换
	ArgDesc_FFI_SIGNATURE
	PrmTyDesc_FFI_INVOKE_FFI
	ArgTyDesc_FFI_SIGNATURE
	ArgTyDesc_FFI_INVOKE_FFI
	RetTyDesc_FFI
	PrmTyDesc_FFI_INVOKE_C

	// 在rust中
	ArgDesc_RS_SIGNATURE
	AsRsSignature // 转换到go函数签名相应的类型
	AsRsITF       // 转换到go函数签名中需要用到interface的相应的类型
	AsRsReturn    // 转换 go函数返回值相应的类型，可能与签名中的不一样
	ArgTyDesc_Rs_INVOKE_Rs
	PrmTyDesc_Rs_INVOKE_Rs // 有时需要做*或者&操作或者强制类型转换
	AsRsCallRsFConv        // 同上, FConv指的是force convert, reference/dereference
	AsRsCallCRsFConv       // 当在一个纯go的函数中调用一个接收cgo类型参数时使用
	AsRsCallFFITy
	PrmTyDesc_Rs_INVOKE_CRs // 同上

	// 在crystal中
	ArgDesc_CR_SIGNATURE
	AsCrSignature // 转换到go函数签名相应的类型
	AsCrITF       // 转换到go函数签名中需要用到interface的相应的类型
	AsCrReturn    // 转换 go函数返回值相应的类型，可能与签名中的不一样
	ArgTyDesc_Cr_INVOKE_Cr
	PrmTyDesc_Cr_INVOKE_Cr // 有时需要做*或者&操作或者强制类型转换
	AsCrCallRsFConv        // 同上, FConv指的是force convert, reference/dereference
	AsCrCallCRsFConv       // 当在一个纯go的函数中调用一个接收cgo类型参数时使用
	AsCrCallFFITy
	PrmTyDesc_Cr_INVOKE_CRs // 同上

	// 在dart中
	ArgDesc_DT_SIGNATURE
	AsDtSignature // 转换到go函数签名相应的类型
	AsDtITF       // 转换到go函数签名中需要用到interface的相应的类型
	AsDtReturn    // 转换 go函数返回值相应的类型，可能与签名中的不一样
	ArgTyDesc_Dt_INVOKE_Dt
	PrmTyDesc_Dt_INVOKE_Dt // 有时需要做*或者&操作或者强制类型转换
	AsDtCallRsFConv        // 同上, FConv指的是force convert, reference/dereference
	AsDtCallCRsFConv       // 当在一个纯go的函数中调用一个接收cgo类型参数时使用
	AsDtCallFFITy
	PrmTyDesc_Dt_INVOKE_Dts // 同上

	// 在vlang中
	ArgDesc_V_SIGNATURE
	AsVSignature // 转换到go函数签名相应的类型
	AsVITF       // 转换到go函数签名中需要用到interface的相应的类型
	AsVReturn    // 转换 go函数返回值相应的类型，可能与签名中的不一样
	ArgTyDesc_V_INVOKE_Rs
	PrmTyDesc_V_INVOKE_Rs // 有时需要做*或者&操作或者强制类型转换
	AsVCallVFConv         // 同上, FConv指的是force convert, reference/dereference
	AsVCallCRsFConv       // 当在一个纯go的函数中调用一个接收cgo类型参数时使用
	AsVCallFFITy
	PrmTyDesc_V_INVOKE_CRs // 同上
	ArgTyDesc_CV_SIGNATURE
)

// 参数与返回值中的类型转换暂存
// 1 key clang.Type表示的是? 可以是
// 2 key int 表示的是转换的方式标识
// 最终的值为转换的结果的字符串描述
var tycvCache = map[clang.Type]map[int]string{}
var argcvCache = map[string]string{}

// cusecs 当前类型引用位置，用于定位模块
func getTyDesc(ty clang.Type, usecat int, usecs clang.Cursor) string {
	che, ok := tycvCache[ty]
	if !ok {
		che = map[int]string{}
		tycvCache[ty] = che
	}
	// 存在的话从暂存里拿，不能在暂存里拿，因为包名前缀可能不同
	if dstr, ok := che[usecat]; ok {
		if false {
			return dstr
		}
	}

	// 类继承的处理，继承某类型，并做特殊处理
	// che[AsGoITF] = getTyDesc(ty, AsGoSignature, usecs)

	// 重新计算
	switch ty.Kind() {
	case clang.Type_Int:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int"
		che[AsGoReturn] = "int"
		che[AsGoSignature] = "int"

		che[AsVReturn] = "int"
		che[AsVSignature] = "int"
		che[ArgTyDesc_CV_SIGNATURE] = "int"

		che[ArgDesc_RS_SIGNATURE] = "i32"
		che[AsRsCallFFITy] = "INT"

		che[ArgDesc_CR_SIGNATURE] = "Int"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_UInt:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "unsigned int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "unsigned int"
		che[AsGoReturn] = "uint"
		che[AsGoSignature] = "uint"

		che[AsVReturn] = "int"
		che[AsVSignature] = "int"
		che[ArgTyDesc_CV_SIGNATURE] = "int"

		che[ArgDesc_RS_SIGNATURE] = "u32"
		che[AsRsCallFFITy] = "UINT32"

		che[ArgDesc_CR_SIGNATURE] = "UInt"
		che[ArgDesc_DT_SIGNATURE] = "int"
	case clang.Type_LongLong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int64_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int64_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int64_t"
		che[AsGoReturn] = "int64"
		che[AsGoSignature] = "int64"

		che[AsVReturn] = "i64"
		che[AsVSignature] = "i64"
		che[ArgTyDesc_CV_SIGNATURE] = "i64"

		che[ArgDesc_RS_SIGNATURE] = "i64"
		che[AsRsCallFFITy] = "SINT64"

		che[ArgDesc_CR_SIGNATURE] = "Int64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_ULongLong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint64_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint64_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "uint64_t"
		che[AsGoReturn] = "uint64"
		che[AsGoSignature] = "uint64"

		che[AsVReturn] = "u64"
		che[AsVSignature] = "u64"
		che[ArgTyDesc_CV_SIGNATURE] = "u64"

		che[ArgDesc_RS_SIGNATURE] = "u64"
		che[AsRsCallFFITy] = "UINT64"

		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Short:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int16_t"
		che[AsGoReturn] = "int16"
		che[AsGoSignature] = "int16"

		che[AsVReturn] = "i16"
		che[AsVSignature] = "i16"
		che[ArgTyDesc_CV_SIGNATURE] = "i16"

		che[ArgDesc_RS_SIGNATURE] = "i16"
		che[AsRsCallFFITy] = "SINT16"

		che[ArgDesc_CR_SIGNATURE] = "Int16"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_UShort:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "uint16_t"
		che[AsGoReturn] = "uint16"
		che[AsGoSignature] = "uint16"

		che[AsVReturn] = "u16"
		che[AsVSignature] = "u16"
		che[ArgTyDesc_CV_SIGNATURE] = "u16"

		che[ArgDesc_RS_SIGNATURE] = "u16"
		che[AsRsCallFFITy] = "UINT16"

		che[ArgDesc_CR_SIGNATURE] = "UInt16"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_UChar:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint8_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint8_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "uint8_t"
		che[AsGoReturn] = "byte"
		che[AsGoSignature] = "byte"

		che[AsVReturn] = "byte"
		che[AsVSignature] = "byte"
		che[ArgTyDesc_CV_SIGNATURE] = "byte"

		che[ArgDesc_RS_SIGNATURE] = "u8"
		che[AsRsCallFFITy] = "UINT8"

		che[ArgDesc_CR_SIGNATURE] = "UInt8"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Char_S:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int8_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int8_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int8_t"
		che[AsGoReturn] = "byte"
		che[AsGoSignature] = "byte"

		che[AsVReturn] = "byte"
		che[AsVSignature] = "byte"
		che[ArgTyDesc_CV_SIGNATURE] = "byte"

		che[ArgDesc_RS_SIGNATURE] = "i8"
		che[AsRsCallFFITy] = "SINT8"

		che[ArgDesc_CR_SIGNATURE] = "Int8"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_SChar:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.char"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "char"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "char"
		che[AsGoReturn] = "byte"
		che[AsGoSignature] = "byte"

		che[AsVReturn] = "byte"
		che[AsVSignature] = "byte"
		che[ArgTyDesc_CV_SIGNATURE] = "byte"

		che[ArgDesc_RS_SIGNATURE] = "i8"
		che[AsRsCallFFITy] = "SINT8"

		che[ArgDesc_CR_SIGNATURE] = "Int8"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Long:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.long"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "long"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "long"
		che[AsGoReturn] = "int64"
		che[AsGoSignature] = "int64"

		che[AsVReturn] = "i64"
		che[AsVSignature] = "i64"
		che[ArgTyDesc_CV_SIGNATURE] = "i64"

		che[ArgDesc_RS_SIGNATURE] = "i64"
		che[AsRsCallFFITy] = "SINT64"

		che[ArgDesc_CR_SIGNATURE] = "Int64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_ULong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.ulong"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "ulong"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "ulong"
		che[AsGoReturn] = "uint64"
		che[AsGoSignature] = "uint64"

		che[AsVReturn] = "u64"
		che[AsVSignature] = "u64"
		che[ArgTyDesc_CV_SIGNATURE] = "u64"

		che[ArgDesc_RS_SIGNATURE] = "u64"
		che[AsRsCallFFITy] = "UINT64"

		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Typedef:
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		if TypeIsQFlags(ty) {
			che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
			che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
			che[AsCReturn] = "int"
			che[AsGoReturn] = "int"
			che[AsGoSignature] = "int"
			che[AsGoITF] = "int"

			che[AsVReturn] = "int"
			che[AsVSignature] = "int"
			che[ArgTyDesc_CV_SIGNATURE] = "int"

			che[ArgDesc_RS_SIGNATURE] = "i32"
			che[AsRsCallFFITy] = "INT"
			che[ArgDesc_CR_SIGNATURE] = "Int"
			che[ArgDesc_DT_SIGNATURE] = "int"
			break
			// typedef template classes
		} else if strings.HasPrefix(ty.CanonicalType().Spelling(), "Q") &&
			strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling())
			tmplArgTy := ty.TemplateArgumentAsType(0)
			if tmplArgTy.Kind() == clang.Type_Pointer {
				tmplArgTy = tmplArgTy.PointeeType()
			}
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling(), tmplArgTy.Spelling())
			refmod := get_decl_mod(tmplArgTy.Declaration())
			usemod := get_decl_mod(usecs)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling(), tmplArgTy.Spelling(), refmod, usemod, pkgPref, usecs.DisplayName(), get_decl_loc(usecs))

			che[AsCSignature] = "*" + ty.Spelling()
			che[AsGoSignature] = "*" + pkgPref + ty.Spelling() + "/*9999*/"
			che[AsGoITF] = pkgPref + ty.Spelling() + "_ITF"
			che[AsGoReturn] = "*" + pkgPref + ty.Spelling() + "/*9999*/"

			che[AsVReturn] = "voidptr"
			che[AsVSignature] = "voidptr"
			che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

			che[ArgDesc_RS_SIGNATURE] = "usize"
			che[AsRsCallFFITy] = "POINTER"
			che[ArgDesc_CR_SIGNATURE] = "UInt64"
			che[ArgDesc_DT_SIGNATURE] = "int"

			break
		} else if is_qt_class(ty.CanonicalType()) {
			refmod := get_decl_mod(ty.CanonicalType().Declaration())
			usemod := get_decl_mod(usecs)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")

			che[AsCSignature] = "*" + ty.Spelling()
			che[AsGoSignature] = "*" + pkgPref + ty.Spelling() + "/*888*/"
			che[AsGoITF] = pkgPref + ty.Spelling() + "_ITF"
			che[AsGoReturn] = "*" + pkgPref + ty.Spelling() + "/*7777*/"

			che[AsVReturn] = "voidptr"
			che[AsVSignature] = "voidptr"
			che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

			che[ArgDesc_RS_SIGNATURE] = "usize"
			che[AsRsCallFFITy] = "POINTER"
			che[ArgDesc_CR_SIGNATURE] = "UInt64"
			che[ArgDesc_DT_SIGNATURE] = "int"

			break
		}
		return getTyDesc(ty.CanonicalType(), usecat, usecs)
	case clang.Type_Record: // TODO qt class
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*444*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"

		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

		if is_qt_class(ty) {
			refmod := get_decl_mod(ty.Declaration())
			usemod := get_decl_mod(usecs)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if strings.ContainsAny(get_bare_type(ty).Spelling(), "<>") {
				log.Println(ty.Spelling(), get_bare_type(ty).Spelling())
			} else {
				che[ArgDesc_DT_SIGNATURE] = get_bare_type(ty).Spelling()
			}
			che[AsCSignature] = "*" + ty.Spelling()
			che[AsGoSignature] = "*" + pkgPref + ty.Spelling() + "/*6666*/"
			che[AsGoITF] = pkgPref + ty.Spelling() + "_ITF"
			che[AsGoITF] = "*" + pkgPref + ty.Spelling() + "/*6666*/"
		}

	case clang.Type_Pointer: // TODO qt class
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*666*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer/*666*/"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr/*666*/"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"
		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

		if isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S {
			che[ArgDesc_DT_SIGNATURE] = "List<String>"
		} else if ty.PointeeType().Kind() == clang.Type_Char_S {
			che[ArgDesc_DT_SIGNATURE] = "String"
		} else if is_qt_class(ty.PointeeType()) {
			che[ArgDesc_DT_SIGNATURE] = get_bare_type(ty.PointeeType()).Spelling()
		}

	case clang.Type_LValueReference: // TODO qt class
		if isPrimitiveType(ty.PointeeType()) {
			// return this.toDest(ty.PointeeType(), cursor)
		} else if is_qt_class(ty.PointeeType()) {
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*555*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer/*555*/"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"
		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_RValueReference:
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*333*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer/*333*/"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"
		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Elaborated:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int"
		che[AsGoReturn] = "int"
		che[AsGoSignature] = "int"

		che[AsVReturn] = "int"
		che[AsVSignature] = "int"
		che[ArgTyDesc_CV_SIGNATURE] = "int"

		che[ArgDesc_RS_SIGNATURE] = "i32"
		che[AsRsCallFFITy] = "INT"
		che[ArgDesc_CR_SIGNATURE] = "Int"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Enum:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int"
		che[AsGoReturn] = "int"
		che[AsGoSignature] = "int"

		che[AsVReturn] = "int"
		che[AsVSignature] = "int"
		che[ArgTyDesc_CV_SIGNATURE] = "int"

		che[ArgDesc_RS_SIGNATURE] = "i32"
		che[AsRsCallFFITy] = "INT"
		che[ArgDesc_CR_SIGNATURE] = "Int"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Bool:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.bool"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "bool"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "bool"
		che[AsGoReturn] = "bool"
		che[AsGoSignature] = "bool"

		che[AsVReturn] = "bool"
		che[AsVSignature] = "bool"
		che[ArgTyDesc_CV_SIGNATURE] = "bool"

		che[ArgDesc_RS_SIGNATURE] = "bool"
		che[AsRsCallFFITy] = "SINT8"
		che[ArgDesc_CR_SIGNATURE] = "Bool"
		che[ArgDesc_DT_SIGNATURE] = "bool"

	case clang.Type_Double:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.double"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "double"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "double"
		che[AsGoReturn] = "float64"
		che[AsGoSignature] = "float64"

		che[AsVReturn] = "f64"
		che[AsVSignature] = "f64"
		che[ArgTyDesc_CV_SIGNATURE] = "f64"

		che[ArgDesc_RS_SIGNATURE] = "f64"
		che[AsRsCallFFITy] = "DOUBLE"
		che[ArgDesc_CR_SIGNATURE] = "Float64"
		che[ArgDesc_DT_SIGNATURE] = "double"

	case clang.Type_LongDouble: // TODO?
		che[ArgTyDesc_CGO_SIGNATURE] = "C.double"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "double"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "double"
		che[AsGoReturn] = "float64"
		che[AsGoSignature] = "float64"

		che[AsVReturn] = "f64"
		che[AsVSignature] = "f64"
		che[ArgTyDesc_CV_SIGNATURE] = "f64"

		che[ArgDesc_RS_SIGNATURE] = "f64"
		che[AsRsCallFFITy] = "DOUBLE"
		che[ArgDesc_CR_SIGNATURE] = "Float64"
		che[ArgDesc_DT_SIGNATURE] = "double"

	case clang.Type_Float:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.float"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "float"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "float"
		che[AsGoReturn] = "float32"
		che[AsGoSignature] = "float32"

		che[AsVReturn] = "f32"
		che[AsVSignature] = "f32"
		che[ArgTyDesc_CV_SIGNATURE] = "f32"

		che[ArgDesc_RS_SIGNATURE] = "f32"
		che[AsRsCallFFITy] = "FLOAT"
		che[ArgDesc_CR_SIGNATURE] = "Float32"
		che[ArgDesc_DT_SIGNATURE] = "double"

	case clang.Type_IncompleteArray:
		// TODO xpm const char *const []
		if TypeIsCharPtr(ty.ElementType()) {
			// return "[]string"
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*222*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer/*222*/"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"
		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_ConstantArray:
		// TODO xpm const char *const []
		if TypeIsCharPtr(ty.ElementType()) {
			// return "[]string"
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*222*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "void*"
		che[AsGoReturn] = "unsafe.Pointer/*222*/"
		che[AsGoSignature] = "unsafe.Pointer"

		che[AsVReturn] = "voidptr"
		che[AsVSignature] = "voidptr"
		che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

		che[ArgDesc_RS_SIGNATURE] = "usize"
		che[AsRsCallFFITy] = "POINTER"
		che[ArgDesc_CR_SIGNATURE] = "UInt64"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Char16:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[AsCReturn] = "int16_t"
		che[AsGoReturn] = "int16"
		che[AsGoSignature] = "int16"

		che[AsVReturn] = "i16"
		che[AsVSignature] = "i16"
		che[ArgTyDesc_CV_SIGNATURE] = "i16"

		che[ArgDesc_RS_SIGNATURE] = "i16"
		che[AsRsCallFFITy] = "SINT16"
		che[ArgDesc_CR_SIGNATURE] = "Int16"
		che[ArgDesc_DT_SIGNATURE] = "int"

	case clang.Type_Void:
		che[ArgTyDesc_CGO_SIGNATURE] = "/*wtf*/"
		che[RetTyDesc_CGO] = "/*void*/"
		che[RetTyDesc_C] = "void"
		che[RetTyDesc_CPP] = "void"
		che[ArgTyDesc_CPP_SIGNAUTE] = "void"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void"
		che[AsCReturn] = "void"
		che[AsGoReturn] = "/*void*/"
		che[AsGoSignature] = "/*void*/"

		che[AsVReturn] = "/*void*/"
		che[AsVSignature] = "/*void*/"

		che[ArgDesc_RS_SIGNATURE] = "(/*void*/)"
		che[AsRsCallFFITy] = "()"
		che[ArgDesc_CR_SIGNATURE] = "Void"
		che[ArgDesc_DT_SIGNATURE] = "void"

	case clang.Type_Unexposed:
		if ty.CanonicalType().Kind() != clang.Type_Unexposed {
			return getTyDesc(ty.CanonicalType(), usecat, usecs)
		} else {
			che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*666*/"
			che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
			che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
			che[AsCReturn] = "void*"
			che[AsGoReturn] = "unsafe.Pointer/*666*/"
			che[AsGoSignature] = "unsafe.Pointer"

			che[AsVReturn] = "voidptr"
			che[AsVSignature] = "voidptr"
			che[ArgTyDesc_CV_SIGNATURE] = "voidptr"

			che[ArgDesc_RS_SIGNATURE] = "usize"
			che[AsRsCallFFITy] = "POINTER"
			che[ArgDesc_CR_SIGNATURE] = "UInt64"
			che[ArgDesc_DT_SIGNATURE] = "int"

			if isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S {
				// return "[]string"
			} else if ty.PointeeType().Kind() == clang.Type_Char_S {
				// return "string"
			} else if is_qt_class(ty.PointeeType()) {
			}
		}
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	// 从经过一次计算的缓存中拿
	if dstr, ok := che[usecat]; ok {
		return dstr
	}
	return fmt.Sprintf("Unknown_Type_%s_%s", ty.Kind().String(), ty.Spelling())
}

type TypeConvertor interface {
	// 目标语言传递参数进来时的类型
	toDest(clang.Type, clang.Cursor) string // dest language
	// 绑定的时候extern的类型
	toBind(clang.Type, clang.Cursor) string // binding type
	// 调用对应C函数时的类型
	toCall(clang.Type, clang.Cursor) string // call C.xxx type

}

// ???
type ValueConvertor interface {
}

type TypeConvertBase struct {
}

func (this *TypeConvertBase) IsQtClass(ty clang.Type) bool {
	return false
}

///
type TypeConvertGo struct {
	TypeConvertBase
}

func NewTypeConvertGo() *TypeConvertGo {
	this := &TypeConvertGo{}
	return this
}

// 把C/C++类型转换为Go的类型表示法
func (this *TypeConvertGo) toDest(ty clang.Type, cursor clang.Cursor) string {
	if strings.Contains(ty.Spelling(), "::Flags") {
		log.Println(ty.Spelling(), ty.Kind().String(), ty.CanonicalType().Spelling(), ty.CanonicalType().Kind().String())
	}
	switch ty.Kind() {
	case clang.Type_Int:
		return "int"
	case clang.Type_UInt:
		return "uint"
	case clang.Type_LongLong:
		return "int64"
	case clang.Type_ULongLong:
		return "uint64"
	case clang.Type_Short:
		return "int16"
	case clang.Type_UShort:
		return "uint16"
	case clang.Type_UChar:
		return "byte"
	case clang.Type_Char_S:
		return "byte"
	case clang.Type_SChar:
		return "byte"
	case clang.Type_Long:
		return "int"
	case clang.Type_ULong:
		return "uint"
	case clang.Type_Typedef:
		if TypeIsQFlags(ty) {
			return "int"
		} else if strings.HasPrefix(ty.CanonicalType().Spelling(), "Q") &&
			strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
			tmplArgTy := ty.TemplateArgumentAsType(0)
			if tmplArgTy.Kind() == clang.Type_Pointer {
				tmplArgTy = tmplArgTy.PointeeType()
			}
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling(), tmplArgTy.Spelling())
			refmod := get_decl_mod(tmplArgTy.Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			return "*" + pkgPref + ty.Spelling() + "/*667*/"
		}
		return this.toDest(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		if strings.HasSuffix(ty.Spelling(), "QList<QVariant>") {
			return "*qtcore.QVariantList" + "/*687*/"
		} else if is_qt_class(ty) {
			refmod := get_decl_mod(ty.Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if strings.ContainsAny(get_bare_type(ty).Spelling(), "<>") {
				log.Println(ty.Spelling(), get_bare_type(ty).Spelling())
			}
			if is_qstring_cls(ty.Spelling()) {
				return "string"
			}
			return "*" + pkgPref + get_bare_type(ty).Spelling() + "/*123*/"
		}
		return "unsafe.Pointer /*444*/"
	case clang.Type_Pointer:
		if isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S {
			return "[]string"
		} else if ty.PointeeType().Kind() == clang.Type_Char_S {
			return "string"
		} else if is_qt_class(ty.PointeeType()) {
			refmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if _, ok := privClasses[ty.PointeeType().Spelling()]; ok {
			} else if usemod == "core" && refmod == "widgets" {
			} else if usemod == "gui" && refmod == "widgets" {
			} else {
				return "*" + pkgPref + get_bare_type(ty).Spelling() +
					fmt.Sprintf("/*777 %s*/", ty.Spelling())
			}
		} else if ty.PointeeType().Kind() == clang.Type_Bool {
			return "*bool"
		}
		return "unsafe.Pointer /*666*/"
	case clang.Type_LValueReference:
		if isPrimitiveType(ty.PointeeType()) {
			return this.toDest(ty.PointeeType(), cursor)
		} else if is_qt_class(ty.PointeeType()) {
			refmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			return "*" + pkgPref + get_bare_type(ty).Spelling()
		}
		return "unsafe.Pointer /*555*/"
	case clang.Type_RValueReference:
		return "unsafe.Pointer /*333*/"
	case clang.Type_Elaborated:
		return "int"
	case clang.Type_Enum:
		return "int"
	case clang.Type_Bool:
		return "bool"
	case clang.Type_Double:
		return "float64"
	case clang.Type_LongDouble:
		return "float64"
	case clang.Type_Float:
		return "float32"
	case clang.Type_IncompleteArray:
		// TODO xpm const char *const []
		if TypeIsCharPtr(ty.ElementType()) {
			return "[]string"
		}
		return "[]interface{}"
	case clang.Type_ConstantArray:
		return "unsafe.Pointer"
	case clang.Type_Char16:
		return "int16"
	case clang.Type_Void:
		return "void"
	case clang.Type_Unexposed:
		if strings.HasPrefix(ty.Spelling(), "QList<") {
			// QList<Qxxx> => QxxxList
			defmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			if strings.Contains(ty.Spelling(), "QCameraInfo") {
				defmod = "multimedia"
			} else if strings.Contains(ty.Spelling(), "QGraphicsItem") {
				defmod = "widgets"
			} else if strings.Contains(ty.Spelling(), "QQuickItem") {
				defmod = "quick"
			}
			pkgPref := gopp.IfElseStr(defmod != usemod, fmt.Sprintf("qt%s.", defmod), "")
			return fmt.Sprintf("*%s%sList /*lll*/", pkgPref, strings.TrimRight(ty.Spelling()[6:], " *>"))
		} else {
			log.Fatalln(ty.Spelling(), ty.Kind().Spelling(),
				cursor.SemanticParent().DisplayName(), cursor.DisplayName())
		}
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling(),
			cursor.SemanticParent().DisplayName(), cursor.DisplayName())
	}
	return fmt.Sprintf("Unknown_%s_%s", ty.Spelling(), ty.Kind().String())
}

func (this *TypeConvertGo) toBind(ty clang.Type, cursor clang.Cursor) string {
	return ""
}

func (this *TypeConvertGo) toCall(ty clang.Type, cursor clang.Cursor) string {
	switch ty.Kind() {
	case clang.Type_Int:
		return "C.int"
	case clang.Type_UInt:
		return "C.uint"
	case clang.Type_LongLong:
		return "C.longlong"
	case clang.Type_ULongLong:
		return "C.ulonglong"
	case clang.Type_Short:
		return "C.short"
	case clang.Type_UShort:
		return "C.ushort"
	case clang.Type_UChar:
		return "C.uchar"
	case clang.Type_Char_S:
		return "C.char"
	case clang.Type_Long:
		return "C.long"
	case clang.Type_ULong:
		return "C.ulong"
	case clang.Type_Typedef:
		return this.toCall(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		return "unsafe.Pointer"
	case clang.Type_Pointer:
		return "unsafe.Pointer"
	case clang.Type_LValueReference:
		return "unsafe.Pointer"
	case clang.Type_RValueReference:
		return "unsafe.Pointer"
	case clang.Type_Elaborated:
		return "C.int"
	case clang.Type_Enum:
		return "C.int"
	case clang.Type_Bool:
		return "C.int"
	case clang.Type_Double:
		return "C.double"
	case clang.Type_Float:
		return "C.float"
	case clang.Type_IncompleteArray:
		return "unsafe.Pointer"
	case clang.Type_Char16:
		return "C.short"
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	return fmt.Sprintf("C.unknown_%s_%s", ty.Spelling(), ty.Kind().String())
}

func (this *TypeConvertGo) toDestMetaType(ty clang.Type, cursor clang.Cursor) string {
	switch ty.Kind() {
	case clang.Type_Int:
		return "qtrt.Int32Ty(false)"
	case clang.Type_UInt:
		return "qtrt.Int32Ty(false)"
	case clang.Type_LongLong:
		return "qtrt.Int64Ty(false)"
	case clang.Type_ULongLong:
		return "qtrt.Int64Ty(false)"
	case clang.Type_Short:
		return "qtrt.Int16Ty(false)"
	case clang.Type_UShort:
		return "qtrt.Int16Ty(false)"
	case clang.Type_UChar:
		return "qtrt.ByteTy(false)"
	case clang.Type_Char_S:
		return "qtrt.ByteTy(false)"
	case clang.Type_Long:
		return "qtrt.Int32Ty(false)"
	case clang.Type_ULong:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Typedef:
		return this.toDestMetaType(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		cmod := get_decl_mod(cursor)
		tmod := get_decl_mod(cursor.Definition())
		if cmod != tmod {
			return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 5", tmod, ty.Spelling())
		} else {
			if strings.Contains(ty.Spelling(), "::") {
				return fmt.Sprintf(" qtrt.VoidpTy() // 16")
			}
			return fmt.Sprintf("reflect.TypeOf(%s{}) // 6", ty.Spelling())
		}
		// return "reflect.TypeOf(qt%s.%s{}) // 1"
	case clang.Type_Pointer:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 3", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 4", canty.Spelling())
			}
		} else {
		recalc:
			switch canty.Kind() {
			case clang.Type_Int:
				return "qtrt.Int32Ty(true)"
			case clang.Type_UInt:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_Char_S:
				return "qtrt.ByteTy(true)"
			case clang.Type_Char_U:
				return "qtrt.ByteTy(true)"
			case clang.Type_UChar:
				return "qtrt.ByteTy(true)"
			case clang.Type_Short:
				return "qtrt.Int16Ty(true)"
			case clang.Type_UShort:
				return "qtrt.UInt16Ty(true)"
			case clang.Type_Char16:
				return "qtrt.UInt16Ty(true)"
			case clang.Type_Char32:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_WChar:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_LongLong:
				return "qtrt.Int64Ty(true)"
			case clang.Type_Float:
				return "qtrt.FloatTy(true)"
			case clang.Type_Double:
				return "qtrt.DoubleTy(true)"
			case clang.Type_Pointer: // for char ** => []string
				return "reflect.TypeOf([]string{})"
			case clang.Type_Typedef:
				// log.Fatalln(canty.Spelling(), canty.CanonicalType().Spelling())
				canty = canty.CanonicalType()
				goto recalc
			case clang.Type_Bool:
				return "qtrt.BoolTy(true)"
			case clang.Type_Void:
				return "qtrt.VoidpTy()"
			case clang.Type_FunctionProto:
				return "qtrt.VoidpTy()"
			case clang.Type_Enum:
				return "qtrt.Int32Ty(true)"
			case clang.Type_Record:
				return "qtrt.VoidpTy()"
			default:
				log.Println("unsupported type:", ty.Spelling(), canty.Spelling(), canty.Kind().String())
			}
		}
		// return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 2", get_decl_mod(cursor), "%")
	case clang.Type_LValueReference:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 3", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 4", canty.Spelling())
			}
		} else {
		recalc2:
			switch canty.Kind() {
			case clang.Type_Int:
				return "qtrt.Int32Ty(false)"
			case clang.Type_UInt:
				return "qtrt.UInt32Ty(false)"
			case clang.Type_Char_S:
				return "qtrt.ByteTy(false)"
			case clang.Type_Short:
				return "qtrt.Int16Ty(false)"
			case clang.Type_Float:
				return "qtrt.FloatTy(false)"
			case clang.Type_Typedef:
				canty = canty.CanonicalType()
				goto recalc2
			case clang.Type_LongLong:
				return "qtrt.Int64Ty(false)"
			case clang.Type_Pointer:
				return this.toDestMetaType(canty, cursor)
			case clang.Type_Record:
				return "qtrt.VoidpTy()"
			default:
				log.Println("unsupported type:", ty.Spelling(), canty.Spelling(), canty.Kind().String())
			}
		}

	case clang.Type_RValueReference:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 8", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 9", canty.Spelling())
			}
		} else {
		}
		// return "reflect.TypeOf(qt%s.%s{}) // 7"
	case clang.Type_Elaborated:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Enum:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Bool:
		return "qtrt.BoolTy(false)"
	case clang.Type_Double:
		return "qtrt.DoubleTy(false)"
	case clang.Type_Float:
		return "qtrt.FloatTy(false)"
	case clang.Type_IncompleteArray:
		return "qtrt.VoidpTy()"
	case clang.Type_Char16:
		return "qtrt.Int16Ty(false)"
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	return fmt.Sprintf("C.unkown_%s_%s", ty.Spelling(), ty.Kind().String())
}

/// vvv
type TypeConvertV struct {
	TypeConvertBase
}

func NewTypeConvertV() *TypeConvertV {
	this := &TypeConvertV{}
	return this
}

// 把C/C++类型转换为Go的类型表示法
func (this *TypeConvertV) toDest(ty clang.Type, cursor clang.Cursor) string {
	if strings.Contains(ty.Spelling(), "::Flags") {
		log.Println(ty.Spelling(), ty.Kind().String(), ty.CanonicalType().Spelling(), ty.CanonicalType().Kind().String())
	}
	switch ty.Kind() {
	case clang.Type_Int:
		return "int"
	case clang.Type_UInt:
		return "int"
	case clang.Type_LongLong:
		return "i64"
	case clang.Type_ULongLong:
		return "u64"
	case clang.Type_Short:
		return "i16"
	case clang.Type_UShort:
		return "u16"
	case clang.Type_UChar:
		return "byte"
	case clang.Type_Char_S:
		return "byte"
	case clang.Type_SChar:
		return "byte"
	case clang.Type_Long:
		return "int"
	case clang.Type_ULong:
		return "int"
	case clang.Type_Typedef:
		if TypeIsQFlags(ty) {
			return "int"
		} else if strings.HasPrefix(ty.CanonicalType().Spelling(), "Q") &&
			strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
			tmplArgTy := ty.TemplateArgumentAsType(0)
			if tmplArgTy.Kind() == clang.Type_Pointer {
				tmplArgTy = tmplArgTy.PointeeType()
			}
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling(), tmplArgTy.Spelling())
			refmod := get_decl_mod(tmplArgTy.Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			return " " + pkgPref + ty.Spelling() + "/*667*/"
		}
		return this.toDest(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		if strings.HasSuffix(ty.Spelling(), "QList<QVariant>") {
			return "qtcore.QVariantList" + "/*687*/"
		} else if is_qt_class(ty) {
			refmod := get_decl_mod(ty.Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if strings.ContainsAny(get_bare_type(ty).Spelling(), "<>") {
				log.Println(ty.Spelling(), get_bare_type(ty).Spelling())
			}
			if is_qstring_cls(ty.Spelling()) {
				return "string"
			}
			return " " + pkgPref + get_bare_type(ty).Spelling() + "/*123*/"
		}
		return "voidptr /*444*/"
	case clang.Type_Pointer:
		if isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S {
			return "[]string"
		} else if ty.PointeeType().Kind() == clang.Type_Char_S {
			return "string"
		} else if is_qt_class(ty.PointeeType()) {
			refmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if _, ok := privClasses[ty.PointeeType().Spelling()]; ok {
			} else if usemod == "core" && refmod == "widgets" {
			} else if usemod == "gui" && refmod == "widgets" {
			} else {
				return " " + pkgPref + get_bare_type(ty).Spelling() +
					fmt.Sprintf("/*777 %s*/", ty.Spelling())
			}
		} else if ty.PointeeType().Kind() == clang.Type_Bool {
			return "&bool"
		}
		return "voidptr /*666*/"
	case clang.Type_LValueReference:
		if isPrimitiveType(ty.PointeeType()) {
			return this.toDest(ty.PointeeType(), cursor)
		} else if is_qt_class(ty.PointeeType()) {
			refmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			pkgPref := gopp.IfElseStr(refmod != usemod, fmt.Sprintf("qt%s.", refmod), "")
			if is_qstring_cls(ty.Spelling()) {
				return "string"
			}
			return " " + pkgPref + get_bare_type(ty).Spelling()
		}
		return "voidptr /*555*/"
	case clang.Type_RValueReference:
		return "voidptr /*333*/"
	case clang.Type_Elaborated:
		return "int"
	case clang.Type_Enum:
		return "int"
	case clang.Type_Bool:
		return "bool"
	case clang.Type_Double:
		return "f64"
	case clang.Type_LongDouble:
		return "f64"
	case clang.Type_Float:
		return "f32"
	case clang.Type_IncompleteArray:
		// TODO xpm const char *const []
		if TypeIsCharPtr(ty.ElementType()) {
			return "[]string"
		}
		return "[]interface{}"
	case clang.Type_ConstantArray:
		return "voidptr"
	case clang.Type_Char16:
		return "i16"
	case clang.Type_Void:
		return "void"
	case clang.Type_Unexposed:
		if strings.HasPrefix(ty.Spelling(), "QList<") {
			// QList<Qxxx> => QxxxList
			defmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			if strings.Contains(ty.Spelling(), "QCameraInfo") {
				defmod = "multimedia"
			} else if strings.Contains(ty.Spelling(), "QGraphicsItem") {
				defmod = "widgets"
			} else if strings.Contains(ty.Spelling(), "QQuickItem") {
				defmod = "quick"
			}
			pkgPref := gopp.IfElseStr(defmod != usemod, fmt.Sprintf("qt%s.", defmod), "")
			return fmt.Sprintf("%s%sList /*lll*/", pkgPref, strings.TrimRight(ty.Spelling()[6:], " *>"))
		} else {
			log.Fatalln(ty.Spelling(), ty.Kind().Spelling(),
				cursor.SemanticParent().DisplayName(), cursor.DisplayName())
		}
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling(),
			cursor.SemanticParent().DisplayName(), cursor.DisplayName())
	}
	return fmt.Sprintf("Unknown_%s_%s", ty.Spelling(), ty.Kind().String())
}

func (this *TypeConvertV) toBind(ty clang.Type, cursor clang.Cursor) string {
	return ""
}

func (this *TypeConvertV) toCall(ty clang.Type, cursor clang.Cursor) string {
	switch ty.Kind() {
	case clang.Type_Int:
		return "C.int"
	case clang.Type_UInt:
		return "C.uint"
	case clang.Type_LongLong:
		return "C.longlong"
	case clang.Type_ULongLong:
		return "C.ulonglong"
	case clang.Type_Short:
		return "C.short"
	case clang.Type_UShort:
		return "C.ushort"
	case clang.Type_UChar:
		return "C.uchar"
	case clang.Type_Char_S:
		return "C.char"
	case clang.Type_Long:
		return "C.long"
	case clang.Type_ULong:
		return "C.ulong"
	case clang.Type_Typedef:
		return this.toCall(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		return "unsafe.Pointer"
	case clang.Type_Pointer:
		return "unsafe.Pointer"
	case clang.Type_LValueReference:
		return "unsafe.Pointer"
	case clang.Type_RValueReference:
		return "unsafe.Pointer"
	case clang.Type_Elaborated:
		return "C.int"
	case clang.Type_Enum:
		return "C.int"
	case clang.Type_Bool:
		return "C.int"
	case clang.Type_Double:
		return "C.double"
	case clang.Type_Float:
		return "C.float"
	case clang.Type_IncompleteArray:
		return "unsafe.Pointer"
	case clang.Type_Char16:
		return "C.short"
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	return fmt.Sprintf("C.unknown_%s_%s", ty.Spelling(), ty.Kind().String())
}

func (this *TypeConvertV) toDestMetaType(ty clang.Type, cursor clang.Cursor) string {
	switch ty.Kind() {
	case clang.Type_Int:
		return "qtrt.Int32Ty(false)"
	case clang.Type_UInt:
		return "qtrt.Int32Ty(false)"
	case clang.Type_LongLong:
		return "qtrt.Int64Ty(false)"
	case clang.Type_ULongLong:
		return "qtrt.Int64Ty(false)"
	case clang.Type_Short:
		return "qtrt.Int16Ty(false)"
	case clang.Type_UShort:
		return "qtrt.Int16Ty(false)"
	case clang.Type_UChar:
		return "qtrt.ByteTy(false)"
	case clang.Type_Char_S:
		return "qtrt.ByteTy(false)"
	case clang.Type_Long:
		return "qtrt.Int32Ty(false)"
	case clang.Type_ULong:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Typedef:
		return this.toDestMetaType(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		cmod := get_decl_mod(cursor)
		tmod := get_decl_mod(cursor.Definition())
		if cmod != tmod {
			return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 5", tmod, ty.Spelling())
		} else {
			if strings.Contains(ty.Spelling(), "::") {
				return fmt.Sprintf(" qtrt.VoidpTy() // 16")
			}
			return fmt.Sprintf("reflect.TypeOf(%s{}) // 6", ty.Spelling())
		}
		// return "reflect.TypeOf(qt%s.%s{}) // 1"
	case clang.Type_Pointer:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 3", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 4", canty.Spelling())
			}
		} else {
		recalc:
			switch canty.Kind() {
			case clang.Type_Int:
				return "qtrt.Int32Ty(true)"
			case clang.Type_UInt:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_Char_S:
				return "qtrt.ByteTy(true)"
			case clang.Type_Char_U:
				return "qtrt.ByteTy(true)"
			case clang.Type_UChar:
				return "qtrt.ByteTy(true)"
			case clang.Type_Short:
				return "qtrt.Int16Ty(true)"
			case clang.Type_UShort:
				return "qtrt.UInt16Ty(true)"
			case clang.Type_Char16:
				return "qtrt.UInt16Ty(true)"
			case clang.Type_Char32:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_WChar:
				return "qtrt.UInt32Ty(true)"
			case clang.Type_LongLong:
				return "qtrt.Int64Ty(true)"
			case clang.Type_Float:
				return "qtrt.FloatTy(true)"
			case clang.Type_Double:
				return "qtrt.DoubleTy(true)"
			case clang.Type_Pointer: // for char ** => []string
				return "reflect.TypeOf([]string{})"
			case clang.Type_Typedef:
				// log.Fatalln(canty.Spelling(), canty.CanonicalType().Spelling())
				canty = canty.CanonicalType()
				goto recalc
			case clang.Type_Bool:
				return "qtrt.BoolTy(true)"
			case clang.Type_Void:
				return "qtrt.VoidpTy()"
			case clang.Type_FunctionProto:
				return "qtrt.VoidpTy()"
			case clang.Type_Enum:
				return "qtrt.Int32Ty(true)"
			case clang.Type_Record:
				return "qtrt.VoidpTy()"
			default:
				log.Println("unsupported type:", ty.Spelling(), canty.Spelling(), canty.Kind().String())
			}
		}
		// return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 2", get_decl_mod(cursor), "%")
	case clang.Type_LValueReference:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 3", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 4", canty.Spelling())
			}
		} else {
		recalc2:
			switch canty.Kind() {
			case clang.Type_Int:
				return "qtrt.Int32Ty(false)"
			case clang.Type_UInt:
				return "qtrt.UInt32Ty(false)"
			case clang.Type_Char_S:
				return "qtrt.ByteTy(false)"
			case clang.Type_Short:
				return "qtrt.Int16Ty(false)"
			case clang.Type_Float:
				return "qtrt.FloatTy(false)"
			case clang.Type_Typedef:
				canty = canty.CanonicalType()
				goto recalc2
			case clang.Type_LongLong:
				return "qtrt.Int64Ty(false)"
			case clang.Type_Pointer:
				return this.toDestMetaType(canty, cursor)
			case clang.Type_Record:
				return "qtrt.VoidpTy()"
			default:
				log.Println("unsupported type:", ty.Spelling(), canty.Spelling(), canty.Kind().String())
			}
		}

	case clang.Type_RValueReference:
		var canty clang.Type
		if ty.PointeeType().Declaration().Type().Spelling() == "" {
			canty = ty.PointeeType()
		} else {
			canty = ty.PointeeType().Declaration().Type()
		}
		if is_qt_class(canty) {
			cmod := get_decl_mod(cursor)
			tmod := get_decl_mod(cursor.Definition())
			if cmod != tmod {
				return fmt.Sprintf("reflect.TypeOf(qt%s.%s{}) // 8", tmod, ty.CanonicalType().Spelling())
			} else {
				return fmt.Sprintf("reflect.TypeOf(%s{}) // 9", canty.Spelling())
			}
		} else {
		}
		// return "reflect.TypeOf(qt%s.%s{}) // 7"
	case clang.Type_Elaborated:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Enum:
		return "qtrt.Int32Ty(false)"
	case clang.Type_Bool:
		return "qtrt.BoolTy(false)"
	case clang.Type_Double:
		return "qtrt.DoubleTy(false)"
	case clang.Type_Float:
		return "qtrt.FloatTy(false)"
	case clang.Type_IncompleteArray:
		return "qtrt.VoidpTy()"
	case clang.Type_Char16:
		return "qtrt.Int16Ty(false)"
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	return fmt.Sprintf("C.unkown_%s_%s", ty.Spelling(), ty.Kind().String())
}

// 是否是基本数据类型的指针的指针
// 像char**
func isPrimitivePPType(ty clang.Type) bool {
	if ty.Kind() == clang.Type_Pointer &&
		ty.PointeeType().Kind() == clang.Type_Pointer &&
		isPrimitiveType(ty.PointeeType().PointeeType()) {
		return true
	}
	return false
}

func isPrimitiveType(ty clang.Type) bool {
	switch ty.Kind() {
	case clang.Type_Int:
		return true
	case clang.Type_UInt:
		return true
	case clang.Type_LongLong:
		return true
	case clang.Type_ULongLong:
		return true
	case clang.Type_Short:
		return true
	case clang.Type_UShort:
		return true
	case clang.Type_UChar:
		return true
	case clang.Type_Char_S:
		return true
	case clang.Type_Long:
		return true
	case clang.Type_ULong:
		return true
	case clang.Type_Typedef:
		return isPrimitiveType(ty.CanonicalType())
	case clang.Type_Record:
		return false
	case clang.Type_Pointer:
		return false
	case clang.Type_LValueReference:
		return false
	case clang.Type_RValueReference:
		return false
	case clang.Type_Elaborated:
		if ty.Spelling() == "QMetaObject::Connection" {
			return false
		}
		return true
	case clang.Type_Enum:
		return true
	case clang.Type_Bool:
		return true
	case clang.Type_Double:
		return true
	case clang.Type_Float:
		return true
	case clang.Type_IncompleteArray:
		return false
	case clang.Type_Char16:
		return true
	case clang.Type_Void:
	default:
		log.Println(ty.Spelling(), ty.Kind().Spelling())
	}
	return false
}

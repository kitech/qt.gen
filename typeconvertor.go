package main

import (
	"fmt"
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
	ArgTyDesc_GO_SIGNATURE
	ArgTyDesc_GO_INVOKE_GO
	RTY_GO
	PrmTyDesc_GO_INVOKE_GO // 有时需要做*或者&操作或者强制类型转换
	PrmTyDesc_GO_INVOKE_CGO
	ATY_GO

	// 在cgo源代码中使用的类型转换
	ArgDesc_CGO_SIGNATURE
	PrmTyDesc_CGO_INVOKE_CGO
	ArgTyDesc_CGO_SIGNATURE
	ArgTyDesc_CGO_INVOKE_CGO
	ArgTyDesc_CGO_INVOKE_GO
	RetTyDesc_CGO
	PrmTyDesc_CGO_INVOKE_GO

	// 在cpp源代码中使用的类型转换
	ArgDesc_CPP_SIGNATURE
	PrmTyDesc_CPP_INVOKE_CPP
	ArgTyDesc_CPP_SIGNAUTE
	ArgTyDesc_CPP_INVOKE_CPP
	RetTyDesc_CPP
	PrmTyDesc_CPP_INVOKE_C

	// 在c源代码中使用的类型转换
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
)

// 参数与返回值中的类型转换暂存
// 1 key clang.Type表示的是? 可以是
// 2 key int 表示的是转换的方式标识
// 最终的值为转换的结果的字符串描述
var tycvCache = map[clang.Type]map[int]string{}
var argcvCache = map[string]string{}

func getTyDesc(ty clang.Type, usecat int) string {
	che, ok := tycvCache[ty]
	if !ok {
		che = map[int]string{}
		tycvCache[ty] = che
	}
	// 存在的话从暂存里拿
	if dstr, ok := che[usecat]; ok {
		return dstr
	}

	// 重新计算
	switch ty.Kind() {
	case clang.Type_Int:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int"
	case clang.Type_UInt:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "unsigned int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "uint"
	case clang.Type_LongLong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int64_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int64_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int64"
	case clang.Type_ULongLong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint64_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint64_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "uint64"
	case clang.Type_Short:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int16"
	case clang.Type_UShort:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "uint16"
	case clang.Type_UChar:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.uint8_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "uint8_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "byte"
	case clang.Type_Char_S:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int8_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int8_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "byte"
	case clang.Type_SChar:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.char"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "char"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "byte"
	case clang.Type_Long:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.long"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "long"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int64"
	case clang.Type_ULong:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.ulong"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "ulong"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "uint64"
	case clang.Type_Typedef:
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		if TypeIsQFlags(ty) {
			che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
			che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
			che[RTY_GO] = "int"
			break
		} else if strings.HasPrefix(ty.CanonicalType().Spelling(), "Q") &&
			strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling())
		}
		return getTyDesc(ty.CanonicalType(), usecat)
	case clang.Type_Record:
		if is_qt_class(ty) {
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*444*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "unsafe.Pointer"
	case clang.Type_Pointer:
		if isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S {
			// return "[]string"
		} else if ty.PointeeType().Kind() == clang.Type_Char_S {
			// return "string"
		} else if is_qt_class(ty.PointeeType()) {
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*666*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "unsafe.Pointer/*666*/"
	case clang.Type_LValueReference:
		if isPrimitiveType(ty.PointeeType()) {
			// return this.toDest(ty.PointeeType(), cursor)
		} else if is_qt_class(ty.PointeeType()) {
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*555*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "unsafe.Pointer/*555*/"
	case clang.Type_RValueReference:
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*333*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "unsafe.Pointer/*333*/"
	case clang.Type_Elaborated:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int"
	case clang.Type_Enum:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int"
	case clang.Type_Bool:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.bool"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "bool"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "bool"
	case clang.Type_Double:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.double"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "double"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "float64"
	case clang.Type_LongDouble: // TODO?
		che[ArgTyDesc_CGO_SIGNATURE] = "C.double"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "double"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "float64"
	case clang.Type_Float:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.float"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "float"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "float32"
	case clang.Type_IncompleteArray:
		// TODO xpm const char *const []
		if TypeIsCharPtr(ty.ElementType()) {
			// return "[]string"
		}
		che[ArgTyDesc_CGO_SIGNATURE] = "unsafe.Pointer  /*222*/"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void*"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "unsafe.Pointer/*222*/"
	case clang.Type_Char16:
		che[ArgTyDesc_CGO_SIGNATURE] = "C.int16_t"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "int16_t"
		che[ArgTyDesc_CPP_SIGNAUTE] = ty.Spelling()
		che[RTY_GO] = "int16"
	case clang.Type_Void:
		che[ArgTyDesc_CGO_SIGNATURE] = "wtf"
		che[RetTyDesc_CGO] = ""
		che[RetTyDesc_C] = "void"
		che[RetTyDesc_CPP] = "void"
		che[ArgTyDesc_CPP_SIGNAUTE] = "void"
		che[ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN] = "void"
		che[RTY_GO] = ""
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	// 从经过一次计算的缓存中拿
	if dstr, ok := che[usecat]; ok {
		return dstr
	}
	return fmt.Sprintf("Unknown_%s_%s", ty.Spelling(), ty.Kind().String())
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

type TypeConvertGo struct {
	TypeConvertBase
}

func NewTypeConvertGo() *TypeConvertGo {
	this := &TypeConvertGo{}
	return this
}

var privClasses = map[string]int{"QV8Engine": 1, "QQmlComponentAttached": 1,
	"QQmlImageProviderBase": 1}

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
			log.Println(ty.Spelling(), ty.CanonicalType().Spelling())
			return "*" + ty.Spelling()
		}
		return this.toDest(ty.CanonicalType(), cursor)
	case clang.Type_Record:
		if is_qt_class(ty) {
			refmod := get_decl_mod(ty.Declaration())
			usemod := get_decl_mod(cursor)
			pkgSuff := ""
			if refmod != usemod {
				pkgSuff = fmt.Sprintf("qt%s.", refmod)
				// log.Println(ty.Spelling(), usemod, refmod)
			}
			if strings.ContainsAny(get_bare_type(ty).Spelling(), "<>") {
				log.Println(ty.Spelling(), get_bare_type(ty).Spelling())
			}
			return "*" + pkgSuff + get_bare_type(ty).Spelling() + "/*123*/"
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
			pkgSuff := ""
			if refmod != usemod {
				pkgSuff = fmt.Sprintf("qt%s.", refmod)
				// log.Println(ty.Spelling(), usemod, refmod)
			}
			if _, ok := privClasses[ty.PointeeType().Spelling()]; ok {
			} else if usemod == "core" && refmod == "widgets" {
			} else if usemod == "gui" && refmod == "widgets" {
			} else {
				return "*" + pkgSuff + get_bare_type(ty).Spelling() +
					fmt.Sprintf("/*777 %s*/", ty.Spelling())
			}
		}
		return "unsafe.Pointer /*666*/"
	case clang.Type_LValueReference:
		if isPrimitiveType(ty.PointeeType()) {
			return this.toDest(ty.PointeeType(), cursor)
		} else if is_qt_class(ty.PointeeType()) {
			refmod := get_decl_mod(get_bare_type(ty).Declaration())
			usemod := get_decl_mod(cursor)
			pkgSuff := ""
			if refmod != usemod {
				pkgSuff = fmt.Sprintf("qt%s.", refmod)
				// log.Println(ty.Spelling(), usemod, refmod)
			}
			return "*" + pkgSuff + get_bare_type(ty).Spelling()
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
	case clang.Type_Char16:
		return "int16"
	case clang.Type_Void:
		return "void"
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
				return "qtrt.UInt32Ty(true)"
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

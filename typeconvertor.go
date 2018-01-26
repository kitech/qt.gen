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

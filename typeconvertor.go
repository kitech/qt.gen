package main

import (
	"fmt"
	"log"

	"github.com/go-clang/v3.9/clang"
)

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

func (this *TypeConvertGo) toDest(ty clang.Type, cursor clang.Cursor) string {

	return ""
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
	default:
		log.Fatalln(ty.Spelling(), ty.Kind().Spelling())
	}
	return fmt.Sprintf("C.unkown_%s_%s", ty.Spelling(), ty.Kind().String())
}

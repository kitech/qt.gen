package main

import (
	"flag"
	"log"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

var filterClass string

func init() {
	flag.StringVar(&filterClass, "fclass", filterClass, "set only one class")
}

type GenFilter interface {
	skipClass(cursor, parent clang.Cursor) bool
	skipMethod(cursor, parent clang.Cursor) bool
	skipArg(cursor, parent clang.Cursor) bool
}

type GenFilterBase struct {
}

func (this *GenFilterBase) skipClass(cursor, parent clang.Cursor) bool {
	skip := this.skipClassImpl(cursor, parent)
	if strings.Contains(cursor.Spelling(), "QPaintDevice") && skip > 0 {
		// log.Fatalln("skipped:", skip)
	}
	return skip > 0
}

func (this *GenFilterBase) skipClassImpl(cursor, parent clang.Cursor) int {
	cname := cursor.Spelling()
	prefixes := []string{
		"QMetaTypeId", "QTypeInfo", "QOpenGLFunctions",
		"QOpenGLExtraFunctions", "QOpenGLVersion", "QOpenGL",
		"QAbstract-", "QPrivate",
	}
	equals := []string{
		"QAbstractOpenGLFunctionsPrivate",
		"QOpenGLFunctionsPrivate",
		"QOpenGLExtraFunctionsPrivate",
		"QAnimationGroup",
		"QMetaType",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(cname, prefix) {
			return 1
		}
	}
	for _, equal := range equals {
		if equal == cname {
			return 2
		}
	}

	// 这个也许是因为qt有bug，也许是因为arch上的qt包有问题。QT_OPENGL_ES_2相关。
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "CoreBackend") {
		return 3
	}
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "DeprecatedBackend") {
		return 4
	}

	if !cursor.IsCursorDefinition() {
		// log.Println("filtered by not definition", cursor.Spelling())
		return 5
	}
	// pure virtual class check
	pure_virtual_class := is_pure_virtual_class(cursor)
	if pure_virtual_class {
		// return true
	}

	// if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True
	if cursor.SpecializedCursorTemplate().IsNull() == false {
		return 6
	}
	// inner class
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return 7
	}
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return 8
	}
	// test
	fixclasses := []string{"QDebug", "QNoDebug", "QDebugStateSaver", "QFileDevice",
		"QLibraryInfo", "QInternal", "QAccessibleObject", "QAccessibleActionInterface",
		"QGraphicsObject"}
	for _, c := range fixclasses {
		if cursor.Spelling() == c {
			return 9
		}
	}
	if cname != "QString" {
		// return true
	}
	if cname != "QStringRef" {
		// return true
	}
	if cname != "QSysInfo" {
		// return true
	}
	if cname != "QCoreApplication" {
		// return true
	}
	if len(filterClass) > 0 && cname != filterClass {
		return 10
	}

	return 0
}

func (this *GenFilterBase) skipMethod(cursor, parent clang.Cursor) bool {
	skip := this.skipMethodImpl(cursor, parent)
	if cursor.Spelling() == "QPaintDevice" && skip > 0 {
		// log.Fatalln(cursor.Spelling(), parent.Spelling(), skip)
	}
	return skip > 0
}

func (this *GenFilterBase) skipMethodImpl(cursor, parent clang.Cursor) int {
	if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
		if cursor.AccessSpecifier() != clang.AccessSpecifier_Protected {
			if cursor.Spelling() == "QPaintDevice" {
				return 0
			}
			return 1
		}
	}

	cname := cursor.Spelling()
	metamths := []string{"qt_metacall", "qt_metacast", "qt_check_for_"}
	for _, mm := range metamths {
		if strings.HasPrefix(cname, mm) {
			return 2
		}
	}

	for _, mm := range []string{"tr", "trUtf8", "data_ptr"} {
		if cname == mm {
			return 3
		}
	}

	if strings.HasPrefix(cname, "operator") {
		return 4
	}

	for _, mm := range []string{"rend", "append", "insert", "rbegin", "prepend", "crend", "crbegin"} {
		if cname == mm {
			return 5
		}
	}

	if cursor.IsVariadic() {
		return 6
	}
	// TODO move ctor and copy ctor?
	if cursor.CXXConstructor_IsCopyConstructor() {
		return 7
	}
	if cursor.CXXConstructor_IsMoveConstructor() {
		return 8
	}

	//
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		if this.skipArg(cursor.Argument(uint32(idx)), cursor) {
			return 9
		}
	}

	return 0
}

func (this *GenFilterBase) skipArg(cursor, parent clang.Cursor) bool {
	skip := this.skipArgImpl(cursor, parent)
	if skip > 0 {
		// log.Println(skip)
	}
	return skip > 0
}

func (this *GenFilterBase) skipArgImpl(cursor, parent clang.Cursor) int {
	// C_ZN16QCoreApplication11aboutToQuitENS_14QPrivateSignalE(void *this_, QCoreApplication::QPrivateSignal a0)
	if strings.Contains(cursor.Type().Spelling(), "QPrivate") {
		return 1
	}
	if strings.HasSuffix(cursor.Type().Spelling(), "Function") {
		return 2
	}
	if strings.HasSuffix(cursor.Type().Spelling(), "Func") {
		return 3
	}
	inenums := []string{
		"ComponentFormattingOptions",
		"FormattingOptions",
		"CategoryFilter",
		"KeyValues",
		"InterfaceFactory",
		"RootObjectHandler",
		"UpdateHandler",
		"QtMetaTypePrivate",
		"va_list",
	}
	for _, inenum := range inenums {
		if strings.Contains(cursor.Type().Spelling(), inenum) {
			return 4
		}
	}
	if cursor.Type().Spelling() == "Id" {
		return 5
	}
	// C_ZN18QThreadStorageDataC1EPFvPvE(void (*)(void *) func) {
	if cursor.Type().Spelling() == "void (*)(void *)" {
		return 6
	}

	if this.skipType(cursor.Type(), cursor) {
		return 7
	}

	return 0
}

func (this *GenFilterBase) skipType(ty clang.Type, cursor clang.Cursor) bool {
	skip := this.skipTypeImpl(ty, cursor)
	if skip > 0 {
		// log.Println(skip)
	}
	return skip > 0
}

func (this *GenFilterBase) skipTypeImpl(ty clang.Type, cursor clang.Cursor) int {

	switch ty.Kind() {
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_RValueReference:
		fallthrough
	case clang.Type_Pointer:
		// is template
		if ty.PointeeType().NumTemplateArguments() != -1 {
			return 1
		}
	case clang.Type_MemberPointer:
		return 2
	case clang.Type_Typedef:
		if false {
			log.Println(ty.Kind().Spelling(), ty.CanonicalType().Kind().Spelling())
		}
		return this.skipTypeImpl(ty.CanonicalType(), cursor)
	default:
		if ty.NumTemplateArguments() != -1 {
			// if strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags<") 这个会过滤掉太多方法
			if !strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags<") {
				return 3
			}
		}
	}

	return 0
}

type GenFilterInc struct {
	GenFilterBase
}

type GenFilterGo struct {
	GenFilterBase
}

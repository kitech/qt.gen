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
	cname := cursor.Spelling()
	prefixes := []string{
		"QMetaTypeId", "QTypeInfo", "QOpenGLFunctions",
		"QOpenGLExtraFunctions", "QOpenGLVersion", "QOpenGL",
		"QAbstract", "QPrivate",
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
			return true
		}
	}
	for _, equal := range equals {
		if equal == cname {
			return true
		}
	}

	// 这个也许是因为qt有bug，也许是因为arch上的qt包有问题。QT_OPENGL_ES_2相关。
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "CoreBackend") {
		return true
	}
	if strings.HasPrefix(cname, "QOpenGLFunctions_") &&
		strings.Contains(cname, "DeprecatedBackend") {
		return true
	}

	if !cursor.IsCursorDefinition() {
		return true
	}
	// pure virtual class check
	pure_virtual_class := false
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.CXXMethod_IsPureVirtual() {
			pure_virtual_class = true
		}
		return clang.ChildVisit_Continue
	})
	if pure_virtual_class {
		return true
	}

	// if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True
	if cursor.SpecializedCursorTemplate().IsNull() == false {
		return true
	}
	// inner class
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return true
	}
	if cursor.SemanticParent().Kind() == clang.Cursor_StructDecl {
		return true
	}
	// test
	fixclasses := []string{"QDebug", "QNoDebug", "QDebugStateSaver", "QFileDevice",
		"QLibraryInfo", "QInternal", "QAccessibleObject", "QGraphicsObject"}
	for _, c := range fixclasses {
		if cursor.Spelling() == c {
			return true
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
		return true
	}

	return false
}
func (this *GenFilterBase) skipMethod(cursor, parent clang.Cursor) bool {
	if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
		return true
	}

	cname := cursor.Spelling()
	metamths := []string{"qt_metacall", "qt_metacast", "qt_check_for_"}
	for _, mm := range metamths {
		if strings.HasPrefix(cname, mm) {
			return true
		}
	}

	for _, mm := range []string{"tr", "trUtf8", "data_ptr"} {
		if cname == mm {
			return true
		}
	}

	if strings.HasPrefix(cname, "operator") {
		return true
	}

	for _, mm := range []string{"rend", "append", "insert", "rbegin", "prepend", "crend", "crbegin"} {
		if cname == mm {
			return true
		}
	}

	if cursor.IsVariadic() {
		return true
	}
	// TODO move ctor and copy ctor?
	if cursor.CXXConstructor_IsCopyConstructor() {
		return true
	}
	if cursor.CXXConstructor_IsMoveConstructor() {
		return true
	}

	//
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		if this.skipArg(cursor.Argument(uint32(idx)), cursor) {
			return true
		}
	}

	return false
}
func (this *GenFilterBase) skipArg(cursor, parent clang.Cursor) bool {
	// C_ZN16QCoreApplication11aboutToQuitENS_14QPrivateSignalE(void *this_, QCoreApplication::QPrivateSignal a0)
	if strings.Contains(cursor.Type().Spelling(), "QPrivate") {
		return true
	}
	if strings.HasSuffix(cursor.Type().Spelling(), "Function") {
		return true
	}
	if strings.HasSuffix(cursor.Type().Spelling(), "Func") {
		return true
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
			return true
		}
	}
	if cursor.Type().Spelling() == "Id" {
		return true
	}
	// C_ZN18QThreadStorageDataC1EPFvPvE(void (*)(void *) func) {
	if cursor.Type().Spelling() == "void (*)(void *)" {
		return true
	}

	if this.skipType(cursor.Type(), cursor) {
		return true
	}

	return false
}

func (this *GenFilterBase) skipType(ty clang.Type, cursor clang.Cursor) bool {

	switch ty.Kind() {
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_RValueReference:
		fallthrough
	case clang.Type_Pointer:
		// is template
		if ty.PointeeType().NumTemplateArguments() != -1 {
			return true
		}
	case clang.Type_MemberPointer:
		return true
	case clang.Type_Typedef:
		if false {
			log.Println(ty.Kind().Spelling(), ty.CanonicalType().Kind().Spelling())
		}
		return this.skipType(ty.CanonicalType(), cursor)
	default:
		if ty.NumTemplateArguments() != -1 {
			return true
		}
	}

	return false
}

type GenFilterInc struct {
	GenFilterBase
}

type GenFilterGo struct {
	GenFilterBase
}

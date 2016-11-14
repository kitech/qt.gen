package main

import (
	"strings"

	"github.com/go-clang/v3.9/clang"
)

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
	}
	equals := []string{
		"QAbstractOpenGLFunctionsPrivate",
		"QOpenGLFunctionsPrivate",
		"QOpenGLExtraFunctionsPrivate",
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
	// if cursor.kind == clidx.CursorKind.CLASS_TEMPLATE: return True
	// test
	if cname != "QString" {
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

	return false
}
func (this *GenFilterBase) skipArg(cursor, parent clang.Cursor) bool {
	return false
}

type GenFilterInc struct {
	GenFilterBase
}

type GenFilterGo struct {
	GenFilterBase
}

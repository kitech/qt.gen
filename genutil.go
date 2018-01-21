package main

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-clang/v3.9/clang"
)

// # like core
func get_decl_mod(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, _, _, _ := loc.FileLocation()
	// log.Println(file.Name())
	if !strings.HasPrefix(file.Name(), "/usr/include/qt") {
		return "stdglobal"
	}
	segs := strings.Split(file.Name(), "/")
	// log.Println(segs[len(segs)-2], segs[len(segs)-2][2:])
	dmod := strings.ToLower(segs[len(segs)-2][2:])
	premods := []string{"core", "gui", "widgets", "network", "qml", "quick"}

	found := false
	for _, m := range premods {
		if m == dmod {
			found = true
			break
		}
	}
	if !found {
		if false {
			log.Printf("unknown module: %s, %s %s, %s\n",
				dmod, cursor.Spelling(), file.Name(), filepath.Base(file.Name()))
			time.Sleep(500 * time.Millisecond)
		}
	}

	return dmod
}

// 计算包名补全
func calc_package_suffix(curc, refc clang.Cursor) string {
	curmod := get_decl_mod(curc)
	refmod := get_decl_mod(refc)
	if refmod != curmod {
		return "qt" + refmod + "."
	}
	return ""
}

func is_qt_class(ty clang.Type) bool {
	nty := get_bare_type(ty)
	name := nty.Spelling()
	if len(name) < 2 {
		return false
	}
	if name[0:1] == "Q" && strings.ToUpper(name[1:2]) == name[1:2] && !strings.Contains(name, "::") {
		return true
	}
	return false
}

func is_private_method(c clang.Cursor) bool {
	return c.Kind() == clang.Cursor_CXXMethod &&
		c.AccessSpecifier() == clang.AccessSpecifier_Private
}

func get_bare_type(ty clang.Type) clang.Type {
	switch ty.Kind() {
	case clang.Type_LValueReference, clang.Type_Pointer:
		return get_bare_type(ty.PointeeType())
	}

	return ty.Declaration().Type()
}

func is_go_keyword(s string) bool {
	keywords := map[string]int{"match": 1, "type": 1, "move": 1, "select": 1, "map": 1,
		"range": 1, "var": 1}
	_, ok := keywords[s]
	return ok
}

// 包含1个以上的纯虚方法
func is_pure_virtual_class(cursor clang.Cursor) bool {
	// pure virtual class check
	pure_virtual_class := false
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.CXXMethod_IsPureVirtual() {
			pure_virtual_class = true
		}
		return clang.ChildVisit_Continue
	})
	return pure_virtual_class
}

func find_base_classes(cursor clang.Cursor) []clang.Cursor {
	bcs := make([]clang.Cursor, 0)
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		if c.Kind() == clang.Cursor_CXXBaseSpecifier {
			bcs = append(bcs, c.Definition())
		}
		if c.Kind() == clang.Cursor_CXXMethod {
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return bcs
}

func TypeIsCharPtrPtr(ty clang.Type) bool {
	return isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S
}

func TypeIsCharPtr(ty clang.Type) bool {
	return ty.Kind() == clang.Type_Pointer && ty.PointeeType().Kind() == clang.Type_Char_S
}

func TypeIsQFlags(ty clang.Type) bool {
	if ty.Kind() == clang.Type_Typedef &&
		strings.HasPrefix(ty.CanonicalType().Spelling(), "QFlags") &&
		strings.ContainsAny(ty.CanonicalType().Spelling(), "<>") {
		if strings.Contains(ty.Spelling(), "::") { // for QFlags<xxx> Qxxx::xxx
			return true
		}
	}
	return false
}

func TypeIsFuncPointer(ty clang.Type) bool {
	return strings.Contains(ty.Spelling(), " (*)(")
}

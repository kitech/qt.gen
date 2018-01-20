package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

// # like core
func get_decl_mod(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, _, _, _ := loc.FileLocation()
	// log.Println(file)
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
		log.Fatalln(fmt.Sprintf("unknown module: %s, %s ", dmod, cursor.Spelling()))
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
	name := ty.Spelling()
	if len(name) < 2 {
		return false
	}
	if name[0:1] == "Q" && strings.ToUpper(name[1:2]) == name[1:2] && !strings.Contains(name, "::") {
		return true
	}
	return false
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

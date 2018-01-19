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

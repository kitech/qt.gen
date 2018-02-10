package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
)

// # like core
func get_decl_mod(cursor clang.Cursor) string {
	loc := cursor.Location()
	file, _, _, _ := loc.FileLocation()
	// log.Println(file.Name())
	if !strings.HasPrefix(file.Name(), "/usr/include/qt") {
		if strings.Contains(file.Name(), "headers/QtCore/") { // fix qRegisterResourceData
			return "core"
		}
		return "stdglobal"
	}
	segs := strings.Split(file.Name(), "/")
	// log.Println(segs[len(segs)-2], segs[len(segs)-2][2:])
	dmod := strings.ToLower(segs[len(segs)-2][2:])
	if !strings.HasPrefix(filepath.Base(filepath.Dir(file.Name())), "Qt") {
		dmod = filepath.Base(filepath.Dir(file.Name()))
	}
	log.Println(cursor.Spelling(), dmod, file.Name())
	if dmod == "android" || dmod == "jni" {
		dmod = "androidextras"
	}
	log.Println(cursor.Spelling(), dmod, file.Name())
	if _, found := modDeps[dmod]; !found {
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

// 去掉reference和pointer,并查找其定义类型名，不带const
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
// TODO 父类也有纯虚方法，并且当前类没有实现该方法
func is_pure_virtual_class(cursor clang.Cursor) bool {
	// pure virtual class check
	pure_virtual_class := false
	extraPureClses := map[string]int{"QAnimationGroup": 1}
	if _, ok := extraPureClses[cursor.Spelling()]; ok {
		return true
	}
	if strings.HasPrefix(cursor.Spelling(), "QAbstract") {
		return true
	}
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.CXXMethod_IsPureVirtual() {
			pure_virtual_class = true
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return pure_virtual_class
}

func find_base_classes(cursor clang.Cursor) []clang.Cursor {
	bcs := make([]clang.Cursor, 0)
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		if c.Kind() == clang.Cursor_CXXBaseSpecifier {
			if _, ok := privClasses[c.Definition().Spelling()]; !ok {
				bcs = append(bcs, c.Definition())
			}
		}
		if c.Kind() == clang.Cursor_CXXMethod {
			return clang.ChildVisit_Break
		}
		return clang.ChildVisit_Continue
	})
	return bcs
}

func has_qobject_base_class(cursor clang.Cursor) bool {
	bcs := find_base_classes(cursor)
	for _, bc := range bcs {
		if bc.Spelling() == "QObject" {
			return true
		}
		if has_qobject_base_class(bc) {
			return true
		}
	}
	return false
}

func has_copy_ctor(cursor clang.Cursor) bool {
	has := false
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Constructor:
			if c.AccessSpecifier() == clang.AccessSpecifier_Public &&
				(c.CXXConstructor_IsCopyConstructor() ||
					c.CXXConstructor_IsMoveConstructor()) {
				has = true
				return clang.ChildVisit_Break
			}
		case clang.Cursor_CXXMethod:

		}
		return clang.ChildVisit_Continue
	})
	return has
}

// has default ctor, no virtual method, no virtual base class
func is_trivial_class(cursor clang.Cursor) bool {
	hasDefaultCtor := false
	hasVirtMethod := false
	hasVirtBaseCls := false
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Constructor:
			if c.AccessSpecifier() == clang.AccessSpecifier_Public &&
				c.CXXConstructor_IsDefaultConstructor() {
				hasDefaultCtor = true
			}
		case clang.Cursor_CXXMethod:
			hasVirtMethod = hasVirtMethod || c.CXXMethod_IsVirtual()
		case clang.Cursor_CXXBaseSpecifier:
			hasVirtBaseCls = hasVirtBaseCls || c.IsVirtualBase()
		}
		return clang.ChildVisit_Recurse
	})
	return hasDefaultCtor && !hasVirtMethod && !hasVirtBaseCls
}

func is_deleted_class(cursor clang.Cursor) bool {
	deleted := false
	arr := map[string]int{"QClipboard": 1, "QInputMethod": 1, "QSessionManager": 1,
		"QPaintDevice": 1, "QPagedPaintDevice": 1, "QScroller": 1, "QStandardPaths": 1}
	if _, ok := arr[cursor.Spelling()]; ok {
		return true
	}
	cursor.Visit(func(c, p clang.Cursor) clang.ChildVisitResult {
		switch c.Kind() {
		case clang.Cursor_Destructor:
		}
		return clang.ChildVisit_Recurse
	})
	return deleted
}

// TODO
func is_qt_private_class(cursor clang.Cursor) bool {
	loc := cursor.Type().Declaration().Definition().Location()
	file, _, _, _ := loc.FileLocation()
	// log.Println(file.Name(), cursor.Spelling(), cursor.IsCursorDefinition(), cursor.Definition().IsCursorDefinition())
	if strings.Contains(file.Name(), "/private/") && strings.HasSuffix(file.Name(), "_p.h") {
		return true
	}
	return false
}

// TODO
func is_qt_inner_class(cursor clang.Cursor) bool {
	return false
}

func is_projected_dtor_class(cursor clang.Cursor) bool {
	protectedDtors := map[string]int{
		"QTextCodec": 1, "QAccessibleInterface": 1, "QTextBlockGroup": 1,
		"QTextObject": 1, "QAccessibleWidget": 1,
	}
	_, ok := protectedDtors[cursor.Spelling()]
	return ok
}

func is_qt_global_func(cursor clang.Cursor) bool {
	// qputenv,qsrand,qCompress
	reg := regexp.MustCompile(`q[A-Z].+`) // 需要生成的全局函数名正则规范
	reg = regexp.MustCompile(`q.+`)       // 需要生成的全局函数名正则规范
	return reg.MatchString(cursor.Spelling())
}

func is_qstring_cls(retPlace string) bool {
	if funk.ContainsString(strings.FieldsFunc(retPlace, func(c rune) bool {
		return strings.Contains(" .*/", string(c))
	}), "QString") {
		return true
	}
	return false
}

func TypeIsCharPtrPtr(ty clang.Type) bool {
	return (isPrimitivePPType(ty) && ty.PointeeType().PointeeType().Kind() == clang.Type_Char_S) ||
		(ty.Kind() == clang.Type_IncompleteArray && TypeIsCharPtr(ty.ElementType()))
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

func TypeIsConsted(ty clang.Type) bool {
	return ty.IsConstQualifiedType() || strings.HasPrefix(ty.Spelling(), "const ")
}

func rewriteOperatorMethodName(name string) string {
	replaces := []string{
		"&=", "_and_equal",
		"^=", "_jian_equal",
		"|=", "_or_equal",
		"+=", "_add_equal",
		"-=", "_minus_equal",
		"==", "_equal_equal",
		"!=", "_not_equal",
		"!", "_not", "=", "_equal",
		"<<", "_left_shift",
		">>", "_right_shift",
		"[]", "_get_index",
		"()", "_fncall",
		"->", "_pointer_selector",
		"<", "_less_than", ">", "_greater_than",
		"&", "_and", "^", "_jian", "|", "_or", "~", "_pozhehao",
		"/", "_div", "*", "_mul", "-", "_minus", "+", "_add",
		" ", "_"}
	valiname := name
	for i := 0; i < len(replaces)/2; i += 1 {
		valiname = strings.Replace(valiname, replaces[i*2], replaces[i*2+1], -1)
	}
	return valiname
}

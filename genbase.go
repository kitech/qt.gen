package main

import (
	"log"

	"github.com/go-clang/v3.9/clang"
)

type Generator interface {
	// init(cursor, parent clang.Cursor)
	// genPassHeader(cursor, parent clang.Cursor)
	genClass(cursor, parent clang.Cursor)

	putMethod(c clang.Cursor)
	putFunc(c clang.Cursor)
	putTmplCls(c clang.Cursor)
	putTmplSpec(c clang.Cursor)
	putEnum(c clang.Cursor)

	genFunctions(cursor, parent clang.Cursor)
	genEnumsGlobal(cursor, parent clang.Cursor)
	genTemplateSpecializedClasses()
}

func init() {
	if false {
		log.Println("hehre")
	}
}

type GenBase struct {
	tu *clang.TranslationUnit

	methods      []clang.Cursor
	funcs        []clang.Cursor
	tmplclses    []clang.Cursor
	tmplclsspecs []clang.Cursor
	enums        []clang.Cursor

	isPureVirtualClass  bool
	hasVirtualProtected bool
	isQObjectClass      bool

	funcMangles map[string]int
}

// TODO is what?
func (this *GenBase) isSignal() bool {
	return false
}

func (this *GenBase) isSlot() bool {
	return false
}

func (this *GenBase) putMethod(c clang.Cursor) {
	this.methods = append(this.methods, c)
}

func (this *GenBase) putFunc(c clang.Cursor) {
	this.funcs = append(this.funcs, c)
}

func (this *GenBase) putTmplCls(c clang.Cursor) {
	this.tmplclses = append(this.tmplclses, c)
}

func (this *GenBase) putTmplSpec(c clang.Cursor) {
	this.tmplclsspecs = append(this.tmplclsspecs, c)
}

func (this *GenBase) putEnum(c clang.Cursor) {
	this.enums = append(this.enums, c)
}

//////
func (this *GenBase) groupFunctionsByModule() map[string][]clang.Cursor {
	rets := map[string][]clang.Cursor{}

	for _, fc := range this.funcs {
		qtmod := get_decl_mod(fc)
		if _, ok := modDeps[qtmod]; !ok {
			log.Println("wtf mod:", qtmod, fc.Spelling())
		} else {
			if _, ok := rets[qtmod]; !ok {
				rets[qtmod] = []clang.Cursor{}
			}
			rets[qtmod] = append(rets[qtmod], fc)
		}
	}

	return rets
}

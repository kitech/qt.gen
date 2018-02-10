package main

import (
	"fmt"
	"gopp"
	"log"

	"github.com/go-clang/v3.9/clang"
)

type Generator interface {
	// init(cursor, parent clang.Cursor)
	// genPassHeader(cursor, parent clang.Cursor)
	genClass(cursor, parent clang.Cursor)

	putMethod(c clang.Cursor)
	putFunc(c clang.Cursor)
	putTmplCls(c clang.Cursor)          // c 类型为clang.Cursor_ClassTemplate
	putPlainTmplClsInst(c clang.Cursor) // c类型为clang.Cursor_ClassDecl
	putTydefTmplClsInst(c clang.Cursor) // c类型为clang.Cursor_TypedefDecl
	putEnum(c clang.Cursor)

	genFunctions(cursor, parent clang.Cursor)
	genEnumsGlobal(cursor, parent clang.Cursor)
	genPlainTmplInstClses()
	genTydefTmplInstClses()
}

func init() {
	if false {
		log.Println("hehre")
	}
}

type GenBase struct {
	tu *clang.TranslationUnit

	methods            []clang.Cursor
	funcs              []clang.Cursor
	tmplclses          []clang.Cursor
	plaintmplinstclses []clang.Cursor
	tydeftmplinstclses []clang.Cursor
	enums              []clang.Cursor

	isPureVirtualClass  bool
	hasVirtualProtected bool
	isQObjectClass      bool
	isDeletedClass      bool
	hasProjectedDtor    bool
	hasExplictDtor      bool

	funcMangles map[string]int

	_argDesc1   []string
	_paramDesc1 []string
	_argtyDesc1 []string
	_argDesc2   []string
	_paramDesc2 []string
	_argtyDesc2 []string
	_argDesc3   []string
	_paramDesc3 []string
	_argtyDesc3 []string
	_argDesc4   []string
	_paramDesc4 []string
	_argtyDesc4 []string
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

func (this *GenBase) putPlainTmplClsInst(c clang.Cursor) {
	this.plaintmplinstclses = append(this.plaintmplinstclses, c)
}

func (this *GenBase) putTydefTmplClsInst(c clang.Cursor) {
	this.tydeftmplinstclses = append(this.tydeftmplinstclses, c)
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

func (this *GenBase) genParamRefName(cursor, parent clang.Cursor, aidx int) string {
	argName := cursor.Spelling()
	argName = gopp.IfElseStr(is_go_keyword(argName), argName+"_", argName)

	return gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName)
}

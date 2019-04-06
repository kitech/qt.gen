package main

import (
	"fmt"
	"gopp"
	"log"
	"strings"

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
	putConstant(c clang.Cursor)

	genFunctions(cursor, parent clang.Cursor)
	genEnumsGlobal(cursor, parent clang.Cursor)
	genPlainTmplInstClses()
	genTydefTmplInstClses()
	genConstantsGlobal(cursor, parent clang.Cursor)
}

func init() {
	if false {
		log.Println("hehre")
	}
}

type GenBase struct {
	tu *clang.TranslationUnit

	qtdir string
	qtver string

	methods            []clang.Cursor
	funcs              []clang.Cursor
	tmplclses          []clang.Cursor
	plaintmplinstclses []clang.Cursor
	tydeftmplinstclses []clang.Cursor
	enums              []clang.Cursor
	constants          []clang.Cursor

	isPureVirtualClass  bool
	hasVirtualProtected bool
	isQObjectClass      bool
	isDeletedClass      bool
	hasProjectedDtor    bool
	hasExplictDtor      bool
	hasNominMethod      bool

	funcMangles map[string]int

	// method indexes, reset perclass
	mthidxs map[string]int

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

	clsidx int
}

// 这个是全局的，不能放在类内吧
var tmplclsifgened = map[string]int{}

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
func (this *GenBase) putConstant(c clang.Cursor) {
	this.constants = append(this.constants, c)
}

func (this *GenBase) getFuncQulities(cursor clang.Cursor) []string {
	qualities := make([]string, 0)
	qualities = append(qualities, strings.Split(cursor.AccessSpecifier().Spelling(), "=")[1])
	if cursor.CXXMethod_IsStatic() {
		qualities = append(qualities, "static")
	}
	if cursor.IsFunctionInlined() {
		qualities = append(qualities, "inline")
	}
	if cursor.CXXMethod_IsPureVirtual() {
		qualities = append(qualities, "purevirtual")
	}
	if cursor.CXXMethod_IsVirtual() {
		qualities = append(qualities, "virtual")
	}
	qualities = append(qualities, cursor.Visibility().String())
	qualities = append(qualities, cursor.Availability().String())
	return qualities
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

// mod lower case, include need camel case
func (this *GenBase) getIncNameByMod(mod string) string {
	for name, _ := range modDepsAll {
		if strings.ToLower(name) == mod && name != mod {
			return "Qt" + name
		}
	}
	return ""
}

func (this *GenBase) nextclsidx() int {
	this.clsidx = gopp.IfElseInt(this.clsidx == 0, 10000, this.clsidx)
	this.clsidx += 1
	return this.clsidx
}

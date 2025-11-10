package main

import (
	"fmt"
	gopp "github.com/kitech/goplusplus"
	"log"
	"strings"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
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
	mangler  GenMangler

	qtdir string
	qtver string

	defcstpfxs []string // = "QT_"

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
	isTemplateClass     bool

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

	keywords map[string]int
	idfmtprop IdentRefmtProp
}

type IdentRefmtProp struct {
	// func name, arg name
	name_titled bool
	name_snake bool
	name_camel bool
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
	hastpl := hasTmplArgRet(cursor)
	if !hastpl {
		if cursor.Kind() == clang.Cursor_CXXMethod &&
	 		! cursor.CXXMethod_IsStatic() {
			xptr := cursor.GetFunctionProtoType()
			log.Println(xptr)
			// fni := clcg.ArrangeCXXMethodType(cursor, cursor)
			// retkd := cursor.ABIArgInfoKind(fni, -1)
			log.Println(gopp.Retn(clcg.GetCXXMethodRetinfo(cursor, cursor))...)
			retkind, _, _, _, _ := clcg.GetCXXMethodRetinfo(cursor, cursor)
			retkd := clang.CGABIArgInfoKind( retkind)
			qualities = append(qualities, retkd.String())
		} else if cursor.Kind() == clang.Cursor_CXXMethod &&
					 		cursor.CXXMethod_IsStatic()  {

		} else {
			fni := clcg.ArrangeFreeFunctionType(cursor)
			retkd := cursor.ABIArgInfoKind(fni, -1)
			qualities = append(qualities, retkd.String())
		}
	}
	qualities = append(qualities, cursor.Visibility().String())
	qualities = append(qualities, cursor.Availability().String())
	return qualities
}

// ////
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

// check empty
// check keyword
// check name rule
func (this *GenBase) genParamRefName(cursor, parent clang.Cursor, aidx int) string {
	argName := cursor.Spelling()
	argName = gopp.IfElseStr(is_go_keyword(argName), argName+"_", argName)

	return gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName)
}

func (this *GenBase) is_keyword(argName string) bool {
	_, ok := this.keywords [argName];
	return ok
}

type GenArgItem struct {
	idx int
	hasdft bool
	dftval string
	convtype GenFFIConvty // 0, 1

	argcs clang.Cursor
	prtcs clang.Cursor
	argty clang.Type

	oriname string
	orival any // nil when not need
	tmpname string
	tmpval any
	convname string
	convval any
	dvnme string
}

func NewGenArgItem(cursor, parent clang.Cursor, idx int) *GenArgItem {
	aitm := &GenArgItem{}
	aitm.idx = idx
	aitm.argcs = cursor
	aitm.prtcs = parent
	aitm.argty = cursor.Type()

	dv, has := has_default_value(cursor)
	aitm.hasdft = has
	aitm.dftval = dv

	aitm.tmpname = fmt.Sprintf("tmpArg%d", idx)
	aitm.convname = fmt.Sprintf("convArg%d", idx)

	return aitm
}

// Arg, Ret
type GenFFIConvty = int
const (
	none = iota
	get_cthis
	qt_record_class
	charptr
	charptrptr
	int_variant
	float_variant
)

func (this *GenBase) NewGenArgItem(cursor, parent clang.Cursor, idx int) * GenArgItem {
	aitm := NewGenArgItem(cursor, parent, idx)
	aitm.oriname = this.genParamRefName(aitm.argcs, aitm.prtcs, aitm.idx)

	aitm.convtype = this.typeToConvty(aitm.argty)
	return aitm
}

func (this *GenBase) typeToConvty(argty clang.Type) GenFFIConvty {
	if TypeIsCharPtrPtr(argty) {
		return charptrptr
	}else if TypeIsCharPtr(argty) {
		return charptr
	}else if   is_qt_class(argty) &&
		funk.ContainsString([]string{"QString", "QByteArray", "QVariant", "QModelIndex", "QUrl",
			"QSize", "QAbstractState" /*"QScreen", "QAction"*/}, get_bare_type(argty).Spelling()) {
		return qt_record_class
	} else if is_qt_class(argty) && !isPrimitiveType(argty.ClassType().CanonicalType()) {
		return get_cthis
	} else {
		// should be direct assign/forword
	}
	return none
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

var gclsidx = 10000

func (this *GenBase) nextclsidx() int {
	// this.clsidx = gopp.IfElseInt(this.clsidx == 0, 10000, this.clsidx)
	// this.clsidx += 1
	// return this.clsidx
	gclsidx += 1
	return gclsidx
}

// FunctionDecl/CXXMethodDecl
func (this *GenBase) protoMatch(c1, cx clang.Cursor) bool {
	c2 := cx

	mgname1 := this.mangler.origin(c1)
	mgname2 := this.mangler.origin(c2)
	log.Println(c1.Spelling(), mgname1, mgname2)

	rety1 := c1.ResultType()
	rety2 := c2.ResultType()
	argc1 := c1.NumArguments()
	argc2 := c2.NumArguments()
	if (c1.Spelling() == c2.Spelling() || "x"+c1.Spelling() == c2.Spelling()) &&
		rety1.Equal(rety2) && argc1 == argc2 {
		matched := true
		for i := 0; i < int(argc1); i++ {
			arg1 := c1.Argument(uint32(i))
			arg2 := c2.Argument(uint32(i))
			aty1 := arg1.Type()
			aty2 := arg2.Type()
			if !aty1.Equal(aty2) {
				matched = false
				break
			}
		}
		if matched {
			isconst1 := c1.CXXMethod_IsConst()
			isconst2 := c2.CXXMethod_IsConst()
			if isconst1 == isconst2 {
				return true
			}
		}
	}

	return false
}

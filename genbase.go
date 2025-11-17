package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gopp "github.com/kitech/goplusplus"

	"github.com/go-clang/v3.9/clang"
	// funk "github.com/thoas/go-funk"
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
	genClassSizes(cursor, parent clang.Cursor)
}

func init() {
	if false {
		log.Println("hehre")
	}
}

type GenBase struct {
	tu *clang.TranslationUnit
	mangler  GenMangler
	cp *CodePager

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

	outdir string
	file_ext string // lang src file ext
	fmt_exe string
	fmt_args []string
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

func (this *GenBase) genClassSizes(cursor, parent clang.Cursor) {
	log.Println("clslens", len(clslens))
	cp := this.cp
	// cp.initblocks()
	cp.APf("main", "#include <string.h>\n")
	cp.APf("main", "#include <stdlib.h>\n")
	cp.APf("main", "extern \"C\" \nint qtinline_get_class_size(char* name) {\n")
	cp.APf("main", "    if (0) {}\n")
	for k,v := range clslens {
		  cp.APf("main", "   else if (strcmp(name, \"%s\") == 0) { return %d; }\n", k, v)
	}
	cp.APf("main", "  return 386;\n")
	cp.APf("main", "}\n")
	// oldcp := this.cp
	qtmod := gopp.IfElseStr(isgenqt3(), "3", "core")
	this.saveCodeToFile(qtmod, "qtclass_sizes")
	// this.cp = oldcp
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

	argcs clang.Cursor
	prtcs clang.Cursor
	argty clang.Type

	oriname string
	orival any // nil when not need
	tmpname string
	tmpval any
	convname string
	convval any
	dftname string

	tycv_item *TypeConvItem
	// dest_tyname string
	// ffi_tyname string

	// swap/tmp value
	ffiprm string // only name part
	// sigtprm string // lang func signature
}

func NewGenArgItem(cursor, parent clang.Cursor, idx int) *GenArgItem {
	aitm := &GenArgItem{}
	aitm.idx = idx
	aitm.argcs = cursor
	aitm.prtcs = parent
	aitm.argty = cursor.Type()
	// for const &, const *, must first trim &*, or no effect
	aitm.argty = aitm.argty.RemoveLocalConst() // const morest useless

	dv, has := has_default_value(cursor)
	aitm.hasdft = has
	aitm.dftval = dv

	aitm.tmpname = fmt.Sprintf("tmpArg%d", idx)
	aitm.convname = fmt.Sprintf("convArg%d", idx)
	aitm.dftname = fmt.Sprintf("dftArg%d", idx)

	return aitm
}

func (this *GenBase) NewGenArgItem(cursor, parent clang.Cursor, idx int) * GenArgItem {
	aitm := NewGenArgItem(cursor, parent, idx)
	aitm.oriname = this.genParamRefName(aitm.argcs, aitm.prtcs, aitm.idx)

	return aitm
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


// save code
func (this *GenBase) saveCode(cursor, _ clang.Cursor) {
	// qtx{yyy}, only yyy
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]
	modname = get_decl_mod(cursor)
	log.Println(file.Name(), modname, filepath.Dir(file.Name()), filepath.Base(filepath.Dir(file.Name())))

	clsname := strings.ToLower(cursor.Spelling())
	this.saveCodeToFile(modname, clsname)

	hasnominmth := false
	for _, mth := range this.methods {
		if ismthnomin(mth) {
			hasnominmth = true
			break
		}
	}
	if hasnominmth { // TODO what
		// this.saveCodeToFileWithCode(modname, clsname+".nomin", this.cpnomin.ExportAll())
	}
}

func (this *GenBase) saveCodeToFile(modname, file string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.%s", modname, file, this.file_ext)
	log.Println(savefile, gopp.FileExist("src/"+modname))
	if !gopp.FileExist("src/" + modname) {
		os.Mkdir("src/"+modname+".miss", 0644)
	}

	// log.Println(this.cp.AllPoints())
	bcc := this.cp.ExportAll()
	if strings.HasPrefix(bcc, "//") {
		bcc = bcc[strings.Index(bcc, "\n"):]
	}
	err := ioutil.WriteFile(savefile, []byte(bcc), 0644)
	gopp.ErrPrint(err, savefile)
	if err != nil {
		// log.Panicln(savefile)
	}

	if genctx.refmt_gened_code {
	// gofmt the code
	cmd := exec.Command(this.fmt_exe, append(this.fmt_args, savefile)...)
	err = cmd.Run()
	gopp.ErrPrint(err, cmd)
	}

}

func (this *GenBase) saveCodeToFileWithCode(modname, file string, bcc string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.%s", modname, file, this.file_ext)
	log.Println(savefile)

	// log.Println(this.cp.AllPoints())
	err := ioutil.WriteFile(savefile, []byte(bcc), 0644)
	gopp.ErrPrint(err, savefile)

	if genctx.refmt_gened_code {
	// gofmt the code
	cmd := exec.Command(this.fmt_exe, append(this.fmt_args, savefile)...)
	err = cmd.Run()
	gopp.ErrPrint(err, cmd)
	}
}

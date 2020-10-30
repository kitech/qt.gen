package main

import (
	"fmt"
	"gopp"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
)

type GenerateNim struct {
	// TODO move to base
	filter   GenFilter
	mangler  GenMangler
	tyconver TypeConvertor

	maxClassSize int64 // 暂存一下类的大小的最大值

	cp          *CodePager
	cpnomin     *CodePager
	cpcs        map[string]*CodePager // mod =>
	cptys       map[string]*CodePager // mod =>, all types declare file/code
	cpfwds      map[string]*CodePager // mod =>,  forward declare file/code
	argDesc     []string              // origin c/c++ language syntax
	paramDesc   []string
	destArgDesc []string // dest language syntax

	genedFiles []string

	GenBase
}

func NewGenerateNim(qtdir, qtver string) *GenerateNim {
	this := &GenerateNim{}
	this.qtdir, this.qtver = qtdir, qtver
	this.filter = &GenFilterGo{}
	this.mangler = NewGoMangler()
	this.tyconver = NewTypeConvertGo()

	this.GenBase.funcMangles = map[string]int{}

	this.initBlocks()

	// long lifetime vars
	this.cptys = make(map[string]*CodePager)
	this.cpfwds = make(map[string]*CodePager)
	for _, mod := range []string{"stdglobal", "core", "gui", "widgets"} {
		this.cptys[mod] = NewCodePager()
		this.cpfwds[mod] = NewCodePager()
	}

	return this
}

func (this *GenerateNim) initBlocks() {
	this.cp = NewCodePager()
	this.cpnomin = NewCodePager()

	this.cpcs = make(map[string]*CodePager)
	blocks := []string{"header", "main", "use", "ext", "body", "keep"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
		this.cp.APf(block, "") // for keep block order
		this.cpnomin.AddPointer(block)
		this.cpnomin.APf(block, "")
		// this.cp.APf(block, "// block begin--- %s", block)
	}
}
func (this *GenerateNim) genClass(cursor, parent clang.Cursor) {
	if false {
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	}
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	this.genFileHeader(cursor, parent)
	this.walkClass(cursor, parent)
	// this.genExterns(cursor, parent)
	this.genImports(cursor, parent)
	this.genProtectedCallbacks(cursor, parent)
	this.genClassDef(cursor, parent)
	this.genMethods(cursor, parent)
	this.genClassEnums(cursor, parent)
	this.final(cursor, parent)
	if cursor.Spelling() == "QMimeType" {

	}

}

func (this *GenerateNim) final(cursor, parent clang.Cursor) {
	// log.Println(this.cp.ExportAll())
	this.saveCode(cursor, parent)

	this.initBlocks()
}
func (this *GenerateNim) saveCode(cursor, parent clang.Cursor) {
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
	this.genedFiles = append(this.genedFiles, modname+"/"+clsname)

	hasnominmth := false
	for _, mth := range this.methods {
		if ismthnomin(mth) {
			hasnominmth = true
			break
		}
	}
	if hasnominmth {
		this.saveCodeToFileWithCode(modname, clsname+".nomin", this.cpnomin.ExportAll("#"))
		this.genedFiles = append(this.genedFiles, modname+"/"+clsname+".nomin")
	}
}

func (this *GenerateNim) saveCodeToFile(modname, file string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.nim", modname, file)
	log.Println(savefile, gopp.FileExist("src/"+modname))
	if !gopp.FileExist("src/" + modname) {
		os.Mkdir("src/"+modname+".miss", 0644)
	}

	// log.Println(this.cp.AllPoints())
	bcc := this.cp.ExportAll("#")
	if strings.HasPrefix(bcc, "#") {
		bcc = bcc[strings.Index(bcc, "\n"):]
	}
	err := ioutil.WriteFile(savefile, []byte(bcc), 0644)
	gopp.ErrPrint(err, savefile)
	if err != nil {
		// log.Panicln(savefile)
	}

	// gofmt the code
	// cmd := exec.Command("/usr/bin/gofmt", []string{"-w", savefile}...)
	// err = cmd.Run()
	// gopp.ErrPrint(err, cmd)

}

func (this *GenerateNim) saveCodeToFileWithCode(modname, file string, bcc string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.nim", modname, file)
	log.Println(savefile)

	// log.Println(this.cp.AllPoints())
	ioutil.WriteFile(savefile, []byte(bcc), 0644)

	// gofmt the code
	// cmd := exec.Command("/usr/bin/gofmt", []string{"-w", savefile}...)
	// err := cmd.Run()
	// gopp.ErrPrint(err, cmd)
}
func (this *GenerateNim) genFileList() map[string]string {
	var ret = make(map[string]string) // modname => code of file list
	for _, name := range this.genedFiles {
		var segs = strings.Split(name, "/")
		if name == "stdglobal" {
			continue
		}
		if _, ok := ret[segs[0]]; !ok {
			ret[segs[0]] = ""
		}
		ret[segs[0]] += fmt.Sprintf("include \"./%s.nim\"\n", segs[1])
	}
	return ret
}
func (this *GenerateNim) final2() {
	for modname, cp := range this.cptys {
		this.saveCodeToFileWithCode(modname, "modtypes", cp.ExportAll("#"))
	}
	for modname, cp := range this.cpfwds {
		this.saveCodeToFileWithCode(modname, "modfwdecls", cp.ExportAll("#"))
	}
}

func (this *GenerateNim) genFileHeader(cursor, parent clang.Cursor) {
	this.genFileHeaderWithcp(this.cp, cursor, parent)
	this.genFileHeaderWithcp(this.cpnomin, cursor, parent)
}
func (this *GenerateNim) genFileHeaderWithcp(cp *CodePager, cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	if !strings.HasPrefix(fullModname, "Qt") { // fix cross platform win/mac
		fullModname = "Qt" + fullModname
	}
	modName := "qt" + get_decl_mod(cursor)
	if fullModname == "Qtandroid" {
		fullModname = "QtAndroidExtras"
	}

	ftpath := strings.ToLower(fmt.Sprintf("%s/%s", modName, filepath.Base(file.Name())))
	if _, ok := clts.qtreqcfgs[ftpath]; ok || cp == this.cpnomin {
		cp.APf("header", "#// +build !minimal")
		cp.APf("header", "")
	}

	cp.APf("header", "# package %s", modName)
	cp.APf("header", "# // %s", fix_inc_name(file.Name()))
	cp.APf("header", "# // #include <%s>", filepath.Base(file.Name()))
	cp.APf("header", "# // #include <%s>", fullModname)
	cp.APf("header", "")
	cp.APf("ext", "")
	cp.APf("ext", "#[")
	cp.APf("ext", "#include <stdlib.h>")
	cp.APf("ext", "// extern C begin: %d", len(this.methods))
}
func (this *GenerateNim) isinnomin(cursor clang.Cursor) bool {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	if !strings.HasPrefix(fullModname, "Qt") { // fix cross platform win/mac
		fullModname = "Qt" + fullModname
	}
	modName := "qt" + get_decl_mod(cursor)
	if fullModname == "Qtandroid" {
		fullModname = "QtAndroidExtras"
	}

	ftpath := strings.ToLower(fmt.Sprintf("%s/%s", modName, filepath.Base(file.Name())))
	if _, ok := clts.qtreqcfgs[ftpath]; ok {
		return ok
	}
	return false
}
func (this *GenerateNim) walkClass(cursor, parent clang.Cursor) {

	methods := make([]clang.Cursor, 0)
	enums := make([]clang.Cursor, 0)

	// pcursor := cursor
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			fallthrough
		case clang.Cursor_Destructor:
			fallthrough
		case clang.Cursor_CXXMethod:
			if !this.filter.skipMethod(cursor, parent) {
				methods = append(methods, cursor)
			} else {
				log.Println("filtered:", cursor.DisplayName(), parent.Spelling())
			}
		case clang.Cursor_UnexposedDecl:
			// log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			file, line, col, _ := cursor.Location().FileLocation()
			if false {
				log.Println(file.Name(), line, col, file.Time())
			}
		case clang.Cursor_EnumDecl:
			enums = append(enums, cursor)
		default:
			if false {
				log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			}
		}
		return clang.ChildVisit_Continue
	})

	this.methods = methods
	this.enums = enums
}

func (this *GenerateNim) genExterns(cursor, parent clang.Cursor) {
	for idx, cursor := range this.methods {
		parent := cursor.SemanticParent()
		if false {
			log.Println(idx, parent)
		}
		// log.Println(cursor.Kind().String(), cursor.DisplayName())
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			fallthrough
		case clang.Cursor_Destructor:
			fallthrough
		default:
			if cursor.CXXMethod_IsStatic() {
			} else {
			}
			this.cp.APf("ext", "extern void %s();", this.mangler.origin(cursor))
		}
	}

	this.cp.APf("ext", "// extern C end: %d", len(this.methods))
}

func (this *GenerateNim) genImports(cursor, parent clang.Cursor) {
	this.genImportsWithcp(this.cp, cursor, parent)
	this.genImportsWithcp(this.cpnomin, cursor, parent)
}
func (this *GenerateNim) genImportsWithcp(cp *CodePager, cursor, parent clang.Cursor) {

	cp.APf("ext", "]#")
	cp.APf("ext", "#// import \"C\"") // 直接import "C"导致编译速度下降n倍

	file, _, _, _ := cursor.Location().FileLocation()
	log.Println(file.Name(), cursor.Spelling(), parent.Spelling())
	modname := get_decl_mod(cursor)

	cp.APf("ext", "import nimlog")
	cp.APf("ext", "import nimplus")
	cp.APf("ext", "import typeinfo")
	cp.APf("ext", "import strformat")
	cp.APf("ext", "import strutils")
	cp.APf("ext", "import qt5/qtrt")
	cp.APf("ext", "import qt5/unsafe")
	for _, dep := range modDeps[modname] {
		cp.APf("ext", "import qt5/qt%s", dep)
	}

	cp.APf("keep", "")
	cp.APf("keep", "#func init_unused_%d() {", this.nextclsidx())
	cp.APf("keep", "#  if false {reflect.TypeOf(123)}")
	cp.APf("keep", "#  if false {reflect.TypeOf(unsafe.Sizeof(0))}")
	cp.APf("keep", "#  if false {fmt.Println(123)}")
	cp.APf("keep", "#  if false {log.Println(123)}")
	cp.APf("keep", "#  if false {qtrt.KeepMe()}")
	for _, dep := range modDeps[modname] {
		cp.APf("keep", "#if false {qt%s.KeepMe()}", dep)
	}
	cp.APf("keep", "#}")
	cp.APf("keep", "{.push hint[XDeclaredButNotUsed]:off.}")
}

func (this *GenerateNim) genClassDef(cursor, parent clang.Cursor) {
	bcs := find_base_classes(cursor)
	bcs = this.filter_base_classes(bcs)

	cp := this.cptys[get_decl_mod(cursor)]

	cp.APf("body", "\n#[")
	cp.APf("body", "%s", queryComment(cursor, this.qtdir, this.qtver))
	cp.APf("body", "]#")
	// genTypeStruct
	cp.APf("body", "type %s = ref object of CObject", cursor.Spelling())
	if len(bcs) == 0 {
		cp.APf("body", "#    *qtrt.CObject")
	} else {
		for _, bc := range bcs {
			cp.APf("body", "#    *%s%s", calc_package_prefix(cursor, bc), bc.Type().Spelling())
			// break
		}
	}
	cp.APf("body", "#}")

	// genTypeInterface, genTypeITF
	cp.APf("body", "#type %s_ITF interface {", cursor.Spelling())
	for _, bc := range bcs {
		cp.APf("body", "#    %s%s_ITF", calc_package_prefix(cursor, bc), bc.Type().Spelling())
		// break
	}
	cp.APf("body", "#    %s_PTR() *%s", cursor.Spelling(), cursor.Spelling())
	cp.APf("body", "#}")
	cp.APf("body", "#func (ptr *%s) %s_PTR() *%s { return ptr }",
		cursor.Spelling(), cursor.Spelling(), cursor.Spelling())
	cp.APf("body", "")

	// this.genGetCthis(cursor, cursor, 0) // 只要定义了结构体，就有GetCthis方法
	// this.genSetCthis(cursor, cursor, 0) // 只要定义了结构体，就有GetCthis方法
	this.genCtorFromPointer(cursor, cursor, 0)
	this.genYaCtorFromPointer(cursor, cursor, 0)
}

func (this *GenerateNim) filter_base_classes(bcs []clang.Cursor) []clang.Cursor {
	newbcs := make([]clang.Cursor, 0)
	for _, bc := range bcs {
		if !this.filter.skipClass(bc, bc.SemanticParent()) {
			newbcs = append(newbcs, bc)
		}
	}
	return newbcs
}

func (this *GenerateNim) genMethods(cursor, parent clang.Cursor) {
	// log.Println("process class:", len(this.methods), cursor.Spelling())
	grpMethods := this.groupMethods()
	// log.Println(len(grpMethods))

	seeDtor := false
	for _, cursors := range grpMethods {
		// this.genMethodHeader(cursors[0], cursors[0].SemanticParent())
		// this.genMethodInit(cursors[0], cursors[0].SemanticParent())

		/*
			for idx, cursor := range cursors {
				this.genVTableTypes(cursor, cursor.SemanticParent(), idx)
			}
		*/

		// this.genNameLookup(cursors[0], cursors[0].SemanticParent())

		// TODO is this range orderer?
		// case x
		for idx, cursor := range cursors {
			parent := cursor.SemanticParent()
			funco, found := qdi.findCoMethodObj(cursor)
			_, _ = funco, found
			// log.Println(idx, cursor.Kind().String(), cursor.DisplayName())
			switch cursor.Kind() {
			case clang.Cursor_Constructor:
				this.genCtor(cursor, parent, idx)
				this.genCtorDvs(cursor, parent, idx)
			case clang.Cursor_Destructor:
				seeDtor = true
				this.genDtor(cursor, parent, idx)
			default:
				if cursor.CXXMethod_IsStatic() {
					this.genStaticMethod(cursor, parent, idx)
					this.genStaticMethodNoThis(cursor, parent, idx)
					this.genStaticMethodDvs(cursor, parent, idx)
				} else {
					this.genNonStaticMethod(cursor, parent, idx)
					this.genNonStaticMethodDvs(cursor, parent, idx)
					log.Println(cursor.DisplayName(), parent.Spelling(), cursor.NumArguments(), num_default_value(cursor))
				}
			}
		}

		// this.genMethodFooter(cursors[0], cursors[0].SemanticParent())
	}
	if !seeDtor {
		this.genDtorNoCode(cursor, parent, 0)
	}
}

// 按名字/重载overload分组
func (this *GenerateNim) groupMethods() [][]clang.Cursor {
	methods2 := make(map[string]int, 0)
	idx := 0
	for _, cursor := range this.methods {
		name := cursor.Spelling()
		if _, ok := methods2[name]; ok {
		} else {
			methods2[name] = idx
			idx += 1
		}
	}
	methods := make([][]clang.Cursor, idx)
	for i := 0; i < idx; i++ {
		methods[i] = make([]clang.Cursor, 0)
	}
	for _, cursor := range this.methods {
		name := cursor.Spelling()
		if eidx, ok := methods2[name]; ok {
			methods[eidx] = append(methods[eidx], cursor)
		} else {
			log.Fatalln(idx, name)
		}
	}
	return methods
}

func ismthnomin_nim(cursor clang.Cursor) bool {
	file, _, _, _ := cursor.Location().FileLocation()
	fileName := strings.Replace(file.Name(), os.Getenv("HOME"), "/home/me", -1)
	pathkey := filepath.Base(filepath.Dir(fileName)) + "/" + filepath.Base(fileName)
	pathkey = strings.ToLower(pathkey)
	return iscursorinLineRanges(cursor, clts.qtcfgexps[pathkey])
}
func (this *GenerateNim) getpropercp(cursor clang.Cursor) *CodePager {
	if ismthnomin(cursor) {
		return this.cpnomin
	}
	if cursor.Kind() == clang.Cursor_FunctionDecl {
		if this.isinnomin(cursor) {
			return this.cpnomin
		}
	}
	return this.cp
}
func (this *GenerateNim) genMethodHeader(cursor, parent clang.Cursor, midx int) {
	file, lineno, _, _ := cursor.Location().FileLocation()
	fileName := strings.Replace(file.Name(), os.Getenv("HOME"), "/home/me", -1)
	var cp = this.getpropercp(cursor)

	cp.APf("body", "#// %s:%d", fix_inc_name(fileName), lineno)
	cp.APf("body", "#// index:%d", midx)

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
	qualities = this.getFuncQulities(cursor)
	if len(qualities) > 0 {
		cp.APf("body", "#// %s", strings.Join(qualities, " "))
	}

	cp.APf("body", "#// [%d] %s %s%s", cursor.ResultType().SizeOf(),
		cursor.ResultType().Spelling(), strings.Replace(cursor.DisplayName(), "class ", "", -1),
		gopp.IfElseStr(cursor.CXXMethod_IsConst(), " const", ""))

	cp.APf("body", "\n#[")
	// cp.APf("body", "%s", queryComment(cursor, this.qtdir, this.qtver))
	cp.APf("body", "]#")
}

func (this *GenerateNim) genMethodInit(cursor, parent clang.Cursor) {
	if cursor.Kind() == clang.Cursor_Constructor {
	}
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		this.cp.APf("body", "func (this *%s) %s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()))
	case clang.Cursor_Destructor:
		this.cp.APf("body", "func (this *%s) Delete%s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()[1:]))
	default:
		this.cp.APf("body", "func (this *%s) %s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()))
	}
	this.cp.AP("body", "  var vtys = make(map[uint8]map[uint8]reflect.Type)")
	this.cp.AP("body", "  if false {fmt.Println(vtys)}")
	this.cp.AP("body", "  var dargExists = make(map[uint8]map[uint8]bool)")
	this.cp.AP("body", "  if false {fmt.Println(dargExists)}")
	this.cp.AP("body", "  var dargValues = make(map[uint8]map[uint8]interface{})")
	this.cp.AP("body", "  if false {fmt.Println(dargValues)}")

	// TODO fill types, default args
}

func (this *GenerateNim) genMethodSignature(cursor, parent clang.Cursor, midx int) {
	if cursor.Kind() == clang.Cursor_Constructor {
	}
	var cp = this.getpropercp(cursor)
	var cpfwd = this.cpfwds[get_decl_mod(parent)]

	this.genArgsDest(cursor, parent, true)
	argStr := strings.Join(this.destArgDesc, ", ")

	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("%d", midx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		prms := funk.Map(this.destArgDesc, func(s string) string {
			return strings.TrimRight(strings.Split(s, " ")[0], " :")
		})
		prmStr := strings.Join(prms.([]string), ", ")
		cp.APf("body", "proc NewForInherit%s(empty: %s, %s) : %s =",
			overloadSuffix, strings.Title(parent.Spelling()), argStr, parent.Spelling())
		cp.APf("body", "    return New%s%s(%s)", cursor.Spelling(), overloadSuffix, prmStr)
		cp.APf("body", "#}")

		cp.APf("body", "proc New%s%s*(%s) : %s =",
			strings.Title(cursor.Spelling()), overloadSuffix, argStr, parent.Spelling())

		cpfwd.APf("body", "proc NewForInherit%s(empty: %s, %s) : %s",
			overloadSuffix, strings.Title(parent.Spelling()), argStr, parent.Spelling())
		cpfwd.APf("body", "proc New%s%s*(%s) : %s",
			strings.Title(cursor.Spelling()), overloadSuffix, argStr, parent.Spelling())
	case clang.Cursor_Destructor:
		cp.APf("body", "proc Delete%s%s*(this : %s) =",
			strings.Title(cursor.Spelling()[1:]), overloadSuffix, parent.Spelling())
		cpfwd.APf("body", "proc Delete%s%s*(this : %s)",
			strings.Title(cursor.Spelling()[1:]), overloadSuffix, parent.Spelling())
	default:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if is_qstring_cls(retPlace) {
			retPlace = "string" /*444*/
		}
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = "" /*333*/
		}
		retPlace = strings.TrimLeft(retPlace, "*")
		cmtpos := strings.Index(retPlace, "/*")
		if cmtpos > 0 {
			retPlace = retPlace[:cmtpos]
		}
		retPlace = gopp.IfElseStr(retPlace == "", retPlace, ":"+retPlace)

		mthname := gopp.IfElseStr(strings.HasPrefix(cursor.Spelling(), "operator"),
			rewriteOperatorMethodName(cursor.Spelling()), cursor.Spelling())
		argStr = gopp.IfElseStr(argStr == "", argStr, ", "+argStr)
		cp.APf("body", "proc %s%s*(this: %s %s) %s =",
			strings.Title(mthname), overloadSuffix, parent.Spelling(), argStr, retPlace)
		cpfwd.APf("body", "proc %s%s*(this: %s %s)%s",
			strings.Title(mthname), overloadSuffix, parent.Spelling(), argStr, retPlace)
	}

	// TODO fill types, default args
}

func (this *GenerateNim) genMethodSignatureDv(cursor, parent clang.Cursor, midx int, dvidx int) {
	if cursor.Kind() == clang.Cursor_Constructor {
	}

	dvn := num_default_value(cursor)
	this.genArgsDest(cursor, parent, true)
	this.destArgDesc = this.dvTrimArg(this.destArgDesc, dvn, dvidx)
	argStr := strings.Join(this.destArgDesc, ", ")
	var cp = this.getpropercp(cursor)
	var cpfwd = this.cpfwds[get_decl_mod(parent)]

	// 后缀有两条下划线的都是处理默认参数的
	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("%d", midx))
	overloadSuffix += gopp.IfElseStr(dvidx == 0, "p", fmt.Sprintf("p%d", dvidx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		prms := funk.Map(this.destArgDesc, func(s string) string {
			return strings.TrimRight(strings.Split(s, " ")[0], " :")
		})
		prmStr := strings.Join(prms.([]string), ", ")
		cp.APf("body", "proc NewForInherit%s(empty: %s, %s) : %s =",
			overloadSuffix, strings.Title(parent.Spelling()), argStr, parent.Spelling())
		cp.APf("body", "  return New%s%s(%s)", cursor.Spelling(), overloadSuffix, prmStr)
		cp.APf("body", "#}")

		cp.APf("body", "proc New%s%s*(%s) : %s =",
			strings.Title(cursor.Spelling()), overloadSuffix, argStr, parent.Spelling())

		cpfwd.APf("body", "proc NewForInherit%s(empty: %s, %s) : %s",
			overloadSuffix, strings.Title(parent.Spelling()), argStr, parent.Spelling())
		cpfwd.APf("body", "proc New%s%s*(%s) : %s",
			strings.Title(cursor.Spelling()), overloadSuffix, argStr, parent.Spelling())

	case clang.Cursor_Destructor:
	default:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if is_qstring_cls(retPlace) {
			retPlace = "string"
		}
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = ""
		}
		retPlace = strings.TrimLeft(retPlace, "*")
		retPlace = gopp.IfElseStr(retPlace == "", retPlace, ":"+retPlace)
		cmtpos := strings.Index(retPlace, "/*")
		if cmtpos > 0 {
			retPlace = retPlace[:cmtpos]
		}
		mthname := gopp.IfElseStr(strings.HasPrefix(cursor.Spelling(), "operator"),
			rewriteOperatorMethodName(cursor.Spelling()), cursor.Spelling())
		argStr = gopp.IfElseStr(argStr == "", argStr, ","+argStr)
		cp.APf("body", "proc %s%s*(this:%s %s) %s =",
			strings.Title(mthname), overloadSuffix, parent.Spelling(), argStr, retPlace)

		cpfwd.APf("body", "proc %s%s*(this:%s %s)%s",
			strings.Title(mthname), overloadSuffix, parent.Spelling(), argStr, retPlace)
	}

	// TODO fill types, default args
}

// only for static member
func (this *GenerateNim) genMethodSignatureNoThis(cursor, parent clang.Cursor, midx int) {
	this.genArgsDest(cursor, parent, true)
	argStr := strings.Join(this.destArgDesc, ", ")
	var cp = this.getpropercp(cursor)

	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("%d", midx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
	case clang.Cursor_Destructor:
	case clang.Cursor_CXXMethod:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if is_qstring_cls(retPlace) {
			retPlace = "string"
		}
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = ""
		}
		retPlace = strings.TrimLeft(retPlace, "*")
		cmtpos := strings.Index(retPlace, "/*")
		if cmtpos > 0 {
			retPlace = retPlace[:cmtpos]
		}
		retPlace = gopp.IfElseStr(retPlace == "", "", ": "+retPlace)
		mthname := gopp.IfElseStr(strings.HasPrefix(cursor.Spelling(), "operator"),
			rewriteOperatorMethodName(cursor.Spelling()), cursor.Spelling())
		cp.APf("body", "proc %s_%s%s*(%s) %s =",
			parent.Spelling(), strings.Title(mthname), overloadSuffix, argStr, retPlace)
	default:
		// wtf
	}

	// TODO fill types, default args
}

func (this *GenerateNim) genMethodFooter(cursor, parent clang.Cursor) {
	var cp = this.getpropercp(cursor)

	cp.APf("body", "  default:")
	cp.APf("body", "    qtrt.ErrorResolve(\"%s\", \"%s\", args)",
		parent.Spelling(), cursor.Spelling())
	cp.APf("body", "  } // end switch")
	cp.APf("body", "}")
}

func (this *GenerateNim) genMethodFooterFFI(cursor, parent clang.Cursor, midx int) {
	var cp = this.getpropercp(cursor)

	cp.APf("body", "#}")
	cp.APf("body", "")
}

func (this *GenerateNim) genVTableTypes(cursor, parent clang.Cursor, midx int) {

	this.cp.APf("body", "  // vtypes %d // dargExists %d // dargValues %d", midx, midx, midx)
	this.cp.APf("body", "  vtys[%d] = make(map[uint8]reflect.Type)", midx)
	this.cp.APf("body", "  dargExists[%d] = make(map[uint8]bool)", midx)
	this.cp.APf("body", "  dargValues[%d] = make(map[uint8]interface{})", midx)

	tyconv := this.tyconver.(*TypeConvertGo)
	for aidx := 0; aidx < int(cursor.NumArguments()); aidx++ {

		arg := cursor.Argument(uint32(aidx))
		aty := arg.Type()
		this.cp.APf("body", "  vtys[%d][%d] = %s", midx, aidx,
			tyconv.toDestMetaType(aty, arg))

		// has default value?
		dvres := arg.Evaluate()
		if cursor.Spelling() == "QCoreApplication" {
			log.Println(dvres.Kind().Spelling(), arg, cursor.DisplayName())
		}
		switch dvres.Kind() {
		case clang.Eval_Int:
			this.cp.APf("body", "  dargExists[%d][%d] = true", midx, aidx)
			this.cp.APf("body", "  dargValues[%d][%d] = %d", midx, aidx, dvres.AsInt())
		case clang.Eval_UnExposed:
			fallthrough
		default:
			this.cp.APf("body", "  dargExists[%d][%d] = false", midx, aidx)
			this.cp.APf("body", "  dargValues[%d][%d] = nil", midx, aidx)
		}

	}
}

func (this *GenerateNim) genNameLookup(cursor, parent clang.Cursor) {
	this.cp.AP("body", "")
	this.cp.AP("body", "  var matchedIndex = qtrt.SymbolResolve(args, vtys)")
	this.cp.AP("body", "  if false {fmt.Println(matchedIndex)}")
	this.cp.AP("body", "  switch matchedIndex {")
}

func (this *GenerateNim) genCtor(cursor, parent clang.Cursor, midx int) {
	// log.Println(this.mangler.origin(cursor))
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)

	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr
	var cp = this.getpropercp(cursor)

	if parent.Type().SizeOf() > this.maxClassSize {
		this.maxClassSize = parent.Type().SizeOf()
	}
	// this.cp.APf("body", "    cthis := qtrt.Calloc(1, 256) // %d", parent.Type().SizeOf())
	this.genArgsConvFFI(cursor, parent, midx)
	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_POINTER, %s)",
		this.mangler.origin(cursor), paramStr)
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	cp.APf("body", "    var gothis = New%sFromPointer(cast[pointer](rv))", parent.Spelling())
	if !has_qobject_base_class(parent) {
		cp.APf("body", "    qtrt.SetFinalizer(gothis, Delete%s)", parent.Spelling())
	} else {
		cp.APf("body", "    qtrt.ConnectDestroyed(gothis, \"%s\")", parent.DisplayName())
	}
	cp.APf("body", "    return gothis")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genCtorDvs(cursor, parent clang.Cursor, midx int) {
	dvn := num_default_value(cursor)
	if dvn == 0 {
		return
	}

	for dvidx := 0; dvidx < dvn; dvidx++ {
		this.genCtorDv(cursor, parent, midx, dvidx)
	}
}

func (this *GenerateNim) genCtorDv(cursor, parent clang.Cursor, midx int, dvidx int) {
	// log.Println(this.mangler.origin(cursor))
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignatureDv(cursor, parent, midx, dvidx)

	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr
	var cp = this.getpropercp(cursor)

	if parent.Type().SizeOf() > this.maxClassSize {
		this.maxClassSize = parent.Type().SizeOf()
	}
	// this.cp.APf("body", "    cthis := qtrt.Calloc(1, 256) // %d", parent.Type().SizeOf())
	this.genArgsConvFFIDv(cursor, parent, midx, dvidx)
	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_POINTER, %s)",
		this.mangler.origin(cursor), paramStr)
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	cp.APf("body", "    var gothis = New%sFromPointer(cast[pointer](rv))", parent.Spelling())
	if !has_qobject_base_class(parent) {
		cp.APf("body", "    qtrt.SetFinalizer(gothis, Delete%s)", parent.Spelling())
	} else {
		cp.APf("body", "    qtrt.ConnectDestroyed(gothis, \"%s\")", parent.DisplayName())
	}
	cp.APf("body", "    return gothis")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genCtorFromPointer(cursor, parent clang.Cursor, midx int) {
	if midx > 0 { // 忽略更多重载
		return
	}
	bcs := find_base_classes(parent)
	bcs = this.filter_base_classes(bcs)

	cp := this.cptys[get_decl_mod(cursor)]

	cp.APf("body", "proc New%sFromPointer*(cthis: pointer) : %s =",
		cursor.Spelling(), cursor.Spelling())
	if len(bcs) == 0 {
		cp.APf("body", "    return %s(Cthis:cthis)", cursor.Spelling())
	} else {
		bcobjs := []string{}
		for i, bc := range bcs {
			pkgSuff := calc_package_prefix(cursor, bc)
			cp.APf("body", "    var bcthis%d = %sNew%sFromPointer(cthis)", i, pkgSuff, bc.Spelling())
			bcobjs = append(bcobjs, fmt.Sprintf("bcthis%d", i))
			// break // TODO multiple base classes
		}
		// bcobjArgs := strings.Join(bcobjs, ", ")
		cp.APf("body", "    return %s(Cthis:cthis)", parent.Spelling())
	}
	cp.APf("body", "#}")
}

func (this *GenerateNim) genYaCtorFromPointer(cursor, parent clang.Cursor, midx int) {
	if midx > 0 { // 忽略更多重载
		return
	}
	cp := this.cptys[get_decl_mod(cursor)]
	// can use ((*Qxxx)nil).NewFromPointer
	cp.APf("body", "proc NewFromPointer(empty: %s, cthis:pointer) : %s =",
		cursor.Spelling(), cursor.Spelling())
	cp.APf("body", "    return New%sFromPointer(cthis)", cursor.Spelling())
	cp.APf("body", "#}")
	cp.APf("body", "")
}

func (this *GenerateNim) genGetCthis(cursor, parent clang.Cursor, midx int) {
	if midx > 0 { // 忽略更多重载
		return
	}
	bcs := find_base_classes(parent)
	bcs = this.filter_base_classes(bcs)

	this.cp.APf("body", "func (this *%s) GetCthis() unsafe.Pointer {", parent.Spelling())
	if len(bcs) == 0 {
		this.cp.APf("body", "    if this == nil{ return nil } else { return this.Cthis }")
	} else {
		for _, bc := range bcs {
			this.cp.APf("body", "    if this == nil {return nil} else {return this.%s.GetCthis() }", bc.Spelling())
			break
		}
	}
	this.cp.APf("body", "}")
}

// 用于动态生成实例，new(Qxxx).SetCthis(cthis)
// 像NewQxxxFromPointer，但是可以先创建空实例，再初始化
func (this *GenerateNim) genSetCthis(cursor, parent clang.Cursor, midx int) {
	if midx > 0 { // 忽略更多重载
		return
	}
	bcs := find_base_classes(parent)
	bcs = this.filter_base_classes(bcs)

	this.cp.APf("body", "func (this *%s) SetCthis(cthis unsafe.Pointer) {", parent.Spelling())
	if len(bcs) == 0 {
		this.cp.APf("body", "    if this.CObject == nil {")
		this.cp.APf("body", "        this.CObject = &qtrt.CObject{cthis}")
		this.cp.APf("body", "    }else{")
		this.cp.APf("body", "        this.CObject.Cthis = cthis")
		this.cp.APf("body", "    }")
	} else {
		for _, bc := range bcs {
			pkgSuff := calc_package_prefix(cursor, bc)
			this.cp.APf("body", "    this.%s = %sNew%sFromPointer(cthis)", bc.Spelling(), pkgSuff, bc.Spelling())
			// break
		}
	}
	this.cp.APf("body", "}")
}

func (this *GenerateNim) genDtor(cursor, parent clang.Cursor, midx int) {
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)
	var cp = this.getpropercp(cursor)

	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_VOID, this.GetCthis())",
		this.mangler.origin(cursor))
	cp.APf("body", "    qtrt.Cmemset(this.GetCthis(), 9, %d)", parent.Type().SizeOf())
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	cp.APf("body", "    this.SetCthis(nil)")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genDtorNoCode(cursor, parent clang.Cursor, midx int) {
	// this.genMethodHeader(cursor, parent, midx)
	// this.genMethodSignature(cursor, parent, midx)
	var cp = this.getpropercp(cursor)

	cp.APf("body", "")
	cp.APf("body", "func Delete%s*(this: %s) =", cursor.Spelling(), cursor.Spelling())
	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"_ZN%d%sD2Ev\", qtrt.FFI_TYPE_VOID, this.GetCthis())",
		len(cursor.Spelling()), cursor.Spelling())
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	cp.APf("body", "    this.SetCthis(nil)")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genNonStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr

	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
	}
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)
	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
	}
	var cp = this.getpropercp(cursor)

	this.genArgsConvFFI(cursor, parent, midx)

	retype := cursor.ResultType() // move like sementic, compiler auto behaiver
	mvexpr := ""                  // move expr
	if retype.Kind() == clang.Type_Record {
		// this.cp.APf("body", "    mv := qtrt.Calloc(1, 256)")
		// mvexpr = ", mv"
	}

	ffirety := "qtrt.FFI_TYPE_POINTER"
	if retype.CanonicalType().Kind() == clang.Type_Float ||
		retype.CanonicalType().Kind() == clang.Type_Double {
		ffirety = "qtrt.FFI_TYPE_DOUBLE"
	}
	if retype.Kind() == clang.Type_Record &&
		(retype.Spelling() == "QSize" || retype.Spelling() == "QSizeF") {
		cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", %s %s, this.GetCthis(), %s)",
			this.mangler.origin(cursor), ffirety, mvexpr, paramStr)
	} else {
		cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", %s %s, this.GetCthis(), %s)",
			this.mangler.origin(cursor), ffirety, mvexpr, paramStr)
	}
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	if retype.Kind() == clang.Type_Record {
		// this.cp.APf("body", "   rv = uint64(uintptr(mv))")
	}
	this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

// default argument value
func (this *GenerateNim) genNonStaticMethodDvs(cursor, parent clang.Cursor, midx int) {
	dvn := num_default_value(cursor)
	if dvn == 0 {
		return
	}
	for dvidx := 0; dvidx < dvn; dvidx++ {
		this.genNonStaticMethodDv(cursor, parent, midx, dvidx)
	}
}

// dvidx keep default argument num
func (this *GenerateNim) genNonStaticMethodDv(cursor, parent clang.Cursor, midx int, dvidx int) {
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr

	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignatureDv(cursor, parent, midx, dvidx)

	this.genArgsConvFFIDv(cursor, parent, midx, dvidx)
	var cp = this.getpropercp(cursor)

	retype := cursor.ResultType() // move like sementic, compiler auto behaiver
	mvexpr := ""                  // move expr
	if retype.Kind() == clang.Type_Record {
		// this.cp.APf("body", "    mv := qtrt.Calloc(1, 256)")
		// mvexpr = ", mv"
	}

	ffirety := "qtrt.FFI_TYPE_POINTER"
	if retype.CanonicalType().Kind() == clang.Type_Float ||
		retype.CanonicalType().Kind() == clang.Type_Double {
		ffirety = "qtrt.FFI_TYPE_DOUBLE"
	}
	if retype.Kind() == clang.Type_Record &&
		(retype.Spelling() == "QSize" || retype.Spelling() == "QSizeF") {
		cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", %s %s, this.GetCthis(), %s)",
			this.mangler.origin(cursor), ffirety, mvexpr, paramStr)
	} else {
		cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", %s %s, this.GetCthis(), %s)",
			this.mangler.origin(cursor), ffirety, mvexpr, paramStr)
	}
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")
	if retype.Kind() == clang.Type_Record {
		// this.cp.APf("body", "   rv = uint64(uintptr(mv))")
	}
	this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
	}
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)
	if cursor.IsVariadic() && cursor.NumArguments() > 0 {
	}

	this.genArgsConvFFI(cursor, parent, midx)
	var cp = this.getpropercp(cursor)

	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_POINTER, %s)",
		this.mangler.origin(cursor), paramStr)
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")

	this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genStaticMethodDvs(cursor, parent clang.Cursor, midx int) {
	dvn := num_default_value(cursor)
	if dvn == 0 {
		return
	}
	for dvidx := 0; dvidx < dvn; dvidx++ {
		this.genStaticMethodDv(cursor, parent, midx, dvidx)
	}
}

func (this *GenerateNim) genStaticMethodDv(cursor, parent clang.Cursor, midx int, dvidx int) {
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignatureDv(cursor, parent, midx, dvidx)
	this.genArgsConvFFIDv(cursor, parent, midx, dvidx)
	var cp = this.getpropercp(cursor)

	cp.APf("body", "    var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_POINTER, %s)",
		this.mangler.origin(cursor), paramStr)
	cp.APf("body", "    qtrt.ErrPrint(err, rv)")

	this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genStaticMethodNoThis(cursor, parent clang.Cursor, midx int) {
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr

	// this.genMethodHeaderLongName(cursor, parent, midx)
	this.genMethodSignatureNoThis(cursor, parent, midx)
	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("%d", midx))
	mthname := gopp.IfElseStr(strings.HasPrefix(cursor.Spelling(), "operator"),
		rewriteOperatorMethodName(cursor.Spelling()), cursor.Spelling())
	var cp = this.getpropercp(cursor)

	// this.cp.APf("body", "    // %d: (%s), (%s)", midx, argStr, paramStr)
	cp.APf("body", "    var nilthis : %s", parent.Spelling())
	if cursor.ResultType().Kind() == clang.Type_Void {
		cp.APf("body", "    nilthis.%s%s(%s)",
			strings.Title(mthname), overloadSuffix, paramStr)
	} else {
		cp.APf("body", "    var rv = nilthis.%s%s(%s)",
			strings.Title(mthname), overloadSuffix, paramStr)
		cp.APf("body", "    return rv")
	}
	// this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateNim) genNonVirtualMethod(cursor, parent clang.Cursor, midx int) {

}

func (this *GenerateNim) genProtectedCallbacks(cursor, parent clang.Cursor) {
	log.Println("process class:", len(this.methods), cursor.Spelling())
	mod := get_decl_mod(cursor)
	if _, ok := this.cpcs[mod]; !ok {
		cp := NewCodePager()
		cp.AddPointer("package")
		cp.AddPointer("extern")
		cp.AddPointer("header")
		cp.AddPointer("body")
		cp.APf("package", "package qt%s", mod)
		cp.APf("package", "/*")
		cp.APf("package", "#include <stdint.h>")
		cp.APf("package", "#include <stdbool.h>")
		cp.APf("header", "*/")
		cp.APf("header", "import \"C\"")
		cp.APf("header", "import \"unsafe\"")
		cp.APf("header", "import \"gopp\"")
		cp.APf("header", "// import \"log\"")
		this.cpcs[mod] = cp
	}
	for midx, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())

		if cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
			this.genProtectedCallback(cursor, parent, midx)
		}
	}

	this.cp.APf("body", "")
}

var inheritMethods_nim = map[string]int{}

func (this *GenerateNim) genProtectedCallback(cursor, parent clang.Cursor, midx int) {
	// this.genMethodHeader(cursor, parent, 0)
	mod := get_decl_mod(cursor)
	cp, _ := this.cpcs[mod]

	this.genArgsCGO(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	argStr = gopp.IfElseStr(len(argStr) > 0, ", "+argStr, argStr)

	this.genArgsCGOSign(cursor, parent)
	argStrSign := strings.Join(this.argDesc, ", ")
	argStrSign = gopp.IfElseStr(len(argStrSign) > 0, ", "+argStrSign, argStrSign)

	this.genParams(cursor, parent)
	prmStr := strings.Join(this.paramDesc, ", ")
	prmStr = gopp.IfElseStr(len(prmStr) > 0, ", "+prmStr, prmStr)

	// cp.APf("extern", "extern void set_callback%s(void* fnptr);", cursor.Mangling())
	cp.APf("extern", "extern void callback%s(void* fnptr %s);", cursor.Mangling(), argStrSign)
	cp.APf("body", "// %s %s", getTyDesc(cursor.ResultType(), ArgTyDesc_CPP_SIGNAUTE, cursor), cursor.DisplayName())
	cp.APf("body", "//export callback%s", cursor.Mangling())
	cp.APf("body", "func callback%s(cthis unsafe.Pointer %s) {", cursor.Mangling(), argStr)
	cp.APf("body", "  // log.Println(cthis, \"%s.%s\")", parent.Spelling(), cursor.Spelling())
	cp.APf("body", "  rvx := qtrt.CallbackAllInherits(cthis, \"%s\" %s)", cursor.Spelling(), prmStr)
	cp.APf("body", "  qtrt.ErrPrint(nil, rvx)")
	cp.APf("body", "}")
	cp.APf("body", "func init(){ qtrt.SetInheritCallback2c(\"%s\", C.callback%s /*nil*/) }", cursor.Mangling(), cursor.Mangling())
	cp.APf("body", "")

	// inherit impl
	if cursor.Kind() != clang.Cursor_Constructor && cursor.Kind() != clang.Cursor_Destructor {
		key := fmt.Sprintf("%s::%s", parent.Spelling(), cursor.Spelling())
		if _, ok := inheritMethods[key]; !ok {
			inheritMethods[key] = 1

			this.genArgsDest(cursor, parent, false)
			argStr := strings.Join(this.destArgDesc, ", ")
			retStr := getTyDesc(cursor.ResultType(), AsGoReturn, parent)
			log.Println(parent.Spelling(), cursor.DisplayName(), argStr, retStr, get_decl_loc(cursor), get_decl_mod(cursor))
			retStr = strings.TrimLeft(retStr, "*")
			cmtpos := strings.Index(retStr, "/*")
			if cmtpos >= 0 {
				retStr = retStr[:cmtpos]
			}
			retStr = gopp.IfElseStr(retStr == "", retStr, ":"+retStr)

			mthname := gopp.IfElseStr(strings.HasPrefix(cursor.Spelling(), "operator"),
				rewriteOperatorMethodName(cursor.Spelling()), cursor.Spelling())

			this.cp.APf("body", "#// %s %s", getTyDesc(cursor.ResultType(), ArgTyDesc_CPP_SIGNAUTE, cursor), cursor.DisplayName())
			this.cp.APf("body", "proc Inherit%s(this: %s, f: proc(%s) %s) =",
				strings.Title(mthname), parent.Spelling(), argStr, retStr)
			this.cp.APf("body", "  qtrt.SetAllInheritCallback(this, \"%s\", f)", cursor.Spelling())
			this.cp.APf("body", "#}")
			this.cp.APf("body", "")
		}
	}
}

func (this *GenerateNim) genArgs(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArg(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.origin(cursor))
}

func (this *GenerateNim) genArg(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.origin(parent))

	if len(cursor.Spelling()) == 0 {
		this.argDesc = append(this.argDesc, fmt.Sprintf("%s arg%d", cursor.Type().Spelling(), idx))
	} else {
		if cursor.Type().Kind() == clang.Type_LValueReference {
			// 转成指针
		}
		if strings.Contains(cursor.Type().CanonicalType().Spelling(), "QFlags<") {
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
				cursor.Type().CanonicalType().Spelling(), cursor.Spelling()))
		} else {
			if cursor.Type().Kind() == clang.Type_IncompleteArray ||
				cursor.Type().Kind() == clang.Type_ConstantArray {
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s unsafe.Pointer",
					cursor.Spelling()))
				// log.Println(cursor.Type().Spelling(), cursor.Type().ArrayElementType().Spelling())
				// idx := strings.Index(cursor.Type().Spelling(), " [")
				// this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s %s",
				//	cursor.Type().Spelling()[0:idx], cursor.Spelling(), cursor.Type().Spelling()[idx+1:]))
			} else {
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
					cursor.Spelling(), cursor.Type().Spelling()))
			}
		}
	}
}

func (this *GenerateNim) genArgsDest(cursor, parent clang.Cursor, asitf bool) {
	this.destArgDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgDest(argc, cursor, idx, asitf)
	}
	// log.Println(strings.Join(this.destArgDesc, ", "), this.mangler.origin(cursor))
}

func (this *GenerateNim) genArgDest(cursor, parent clang.Cursor, idx int, asitf bool) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.origin(parent), get_bare_type(cursor.Type()).Spelling(), is_qt_class(cursor.Type()))

	argName := this.genParamRefName(cursor, parent, idx)

	destTy := this.tyconver.toDest(cursor.Type(), cursor)
	if pos := strings.Index(destTy, "/*"); pos > 0 {
		destTy = destTy[:pos]
	}
	destTy = strings.TrimLeft(destTy, "*")
	if destTy == "[]string" {
		destTy = "openArray[string]"
	}

	if cursor.Type().Kind() == clang.Type_LValueReference {
		// 转成指针
	}
	if strings.HasPrefix(cursor.Type().CanonicalType().Spelling(), "QFlags<") {
		this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: int", argName))
	} else if is_qt_class(cursor.Type()) && get_bare_type(cursor.Type()).Spelling() == "QString" {
		this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: string", argName))
	} else if is_qt_class(cursor.Type()) {
		destTyITF := destTy
		if asitf && (strings.HasPrefix(destTy, "*Q") || strings.Contains(destTy, ".Q")) {
			if pos := strings.Index(destTy, "/*"); pos > 0 {
				// destTyITF = destTy[1:pos] + "_ITF" + destTy[pos:]
				destTyITF = destTyITF[1:pos]
			} else {
				// destTyITF = strings.TrimLeft(destTy, "*") + "_ITF"
			}
		}
		this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: %s", argName, destTyITF))
	} else {
		if cursor.Type().Kind() == clang.Type_IncompleteArray {
			this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: %s", argName, destTy))
		} else if cursor.Type().Kind() == clang.Type_ConstantArray {
			this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: %s", argName, destTy))
			// idx := strings.Index(cursor.Type().Spelling(), " [")
			// this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s %s %s",
			// 	cursor.Type().Spelling()[0:idx], argName, cursor.Type().Spelling()[idx+1:]))
		} else {
			this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s: %s", argName, destTy))
		}
	}
}

func (this *GenerateNim) dvTrimArg(argsDesc []string, dvn int, dvidx int) []string {
	return argsDesc[:len(argsDesc)-dvn+dvidx]
}

// midx method index
func (this *GenerateNim) genArgsConv(cursor, parent clang.Cursor, midx int) {
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgConv(argc, cursor, midx, idx)
	}
}

// midx method index
// aidx method index
func (this *GenerateNim) genArgConv(cursor, parent clang.Cursor, midx, aidx int) {
	var cp = this.getpropercp(parent)

	cp.APf("body", "	   var arg%d %s", aidx, this.tyconver.toCall(cursor.Type(), parent))
	cp.APf("body", "	   // if %d >= len(args) {", aidx)
	cp.APf("body", "	   //	  arg%d = defaultargx", aidx)
	cp.APf("body", "	   // } else {")
	cp.APf("body", "	   //	  arg%d = argx.toBind", aidx)
	cp.APf("body", "	   // }")
}

// midx method index
func (this *GenerateNim) genArgsConvFFI(cursor, parent clang.Cursor, midx int) {
	log.Println("gggggggggg", cursor.Spelling(), cursor.ResultType().Kind(), cursor.ResultType().Spelling(), parent.Spelling())
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgConvFFI(argc, cursor, midx, idx)
	}
}

// midx method index
// aidx method index
func (this *GenerateNim) genArgConvFFI(cursor, parent clang.Cursor, midx, aidx int) {
	var cp = this.getpropercp(parent)

	argty := cursor.Type()
	// barety := get_bare_type(argty)
	if TypeIsCharPtrPtr(argty) {
		cp.APf("body", "    var convArg%d = qtrt.StringSliceToCCharPP(%s)", aidx,
			this.genParamRefName(cursor, parent, aidx))
	} else if TypeIsCharPtr(argty) {
		cp.APf("body", "    var convArg%d = qtrt.CString(%s)", aidx,
			this.genParamRefName(cursor, parent, aidx))
		cp.APf("body", "    defer: qtrt.FreeMem(convArg%d)", aidx)
	} else if is_qt_class(argty) && get_bare_type(argty).Spelling() == "QString" {
		usemod := get_decl_mod(cursor)
		pkgPref := gopp.IfElseStr(usemod == "core", "", "qtcore.")
		cp.APf("body", "    var tmpArg%d = %sNewQString5(%s)", aidx, pkgPref,
			this.genParamRefName(cursor, parent, aidx))
		// this.cp.APf("body", "    defer: %sDeleteQString(tmpArg%d)", pkgPref, aidx) // not needed
		cp.APf("body", "    var convArg%d = tmpArg%d.GetCthis()", aidx, aidx)
	} else if is_qt_class(argty) && !isPrimitiveType(argty.CanonicalType()) {
		if argty.Spelling() == "QRgb" {
			log.Fatalln(argty.Spelling(), argty.CanonicalType().Kind().String())
		}
		refmod := get_decl_mod(argty.PointeeType().Declaration())
		usemod := get_decl_mod(cursor)
		log.Println("kkkkk", refmod, usemod, parent.Spelling())
		if _, ok := privClasses[argty.PointeeType().Spelling()]; ok {
		} else if usemod == "core" && refmod == "widgets" {
		} else if usemod == "gui" && refmod == "widgets" {
		} else {
			cp.APf("body", "    var convArg%d : pointer", aidx)
			cp.APf("body", "    if %s != nil:",
				this.genParamRefName(cursor, parent, aidx))
			cp.APf("body", "        convArg%d = %s.GetCthis()", aidx,
				this.genParamRefName(cursor, parent, aidx))
			cp.APf("body", "    #}")
		}
	} else { // no convert needed
		// log.Fatalln("wtf", argty.Kind(), argty.Spelling(), parent.Spelling())
	}
}

// midx method index
func (this *GenerateNim) genArgsConvFFIDv(cursor, parent clang.Cursor, midx int, dvidx int) {
	log.Println("gggggggggg", cursor.Spelling(), cursor.ResultType().Kind(), cursor.ResultType().Spelling(), parent.Spelling())
	dvn := num_default_value(cursor)
	argn := int(cursor.NumArguments())
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		if idx < (argn - dvn + dvidx) {
			this.genArgConvFFI(argc, cursor, midx, idx)
		} else {
			this.genArgConvFFIDv(argc, cursor, midx, idx, dvidx)
		}
	}
}

// midx method index
// aidx method index
func (this *GenerateNim) genArgConvFFIDv(cursor, parent clang.Cursor, midx, aidx int, dvidx int) {
	argdv, _ := has_default_value(cursor)
	argty := cursor.Type()
	barety := get_bare_type(argty)
	undty := barety.Declaration().TypedefDeclUnderlyingType()
	var cp = this.getpropercp(parent)

	argdvs := map[string]string{
		"SH_Default":       "QStyleHintReturn__SH_Default",
		"SO_Default":       "QStyleOption__SO_Default",
		"SO_Complex":       "QStyleOption__SO_Complex",
		"ApplicationFlags": "0",
		"Q_NULLPTR":        "unsafe.Pointer(nil)",
		"nullptr":          "unsafe.Pointer(nil)",
		"Type":             "0",
		"USHRT_MAX":        "-1",
		"ULONG_MAX":        "-1",
	}
	_ = argdvs

	cp.APf("body", "    #// arg: %d, %s=%s, %s=%s, %s, %s", aidx,
		argty.Spelling(), argty.Kind().String(), barety.Spelling(), barety.Kind().String(),
		undty.Spelling(), undty.Kind().String())

	if TypeIsCharPtrPtr(argty) {
		cp.APf("body", "    var convArg%d = qtrt.StringSliceToCCharPP(%s)", aidx,
			this.genParamRefName(cursor, parent, aidx))
	} else if TypeIsCharPtr(argty) {
		cp.APf("body", "    var convArg%d : pointer", aidx)
	} else if funk.Contains([]clang.TypeKind{clang.Type_Enum, clang.Type_Elaborated}, argty.Kind()) {
		cp.APf("body", "    var %s = 0", this.genParamRefName(cursor, parent, aidx))
	} else if argty.Kind() == clang.Type_LValueReference &&
		funk.Contains([]clang.TypeKind{clang.Type_Enum, clang.Type_Elaborated}, argty.PointeeType().Kind()) {
		cp.APf("body", "    var %s = 0", this.genParamRefName(cursor, parent, aidx))
	} else if funk.Contains([]clang.TypeKind{clang.Type_Int, clang.Type_Long, clang.Type_ULong, clang.Type_LongLong, clang.Type_Double, clang.Type_UShort, clang.Type_Float}, argty.Kind()) {
		if strings.HasPrefix(argdv, "Qt::") || argdv == "Type" ||
			(strings.HasPrefix(argdv, "Q") && strings.Contains(argdv, "::")) {
			cp.APf("body", "    var %s = 0 # %s", this.genParamRefName(cursor, parent, aidx), argdv)
		} else if tmpdv, ok := argdvs[argdv]; ok {
			cp.APf("body", "    var %s = %s", this.genParamRefName(cursor, parent, aidx), tmpdv)
		} else {
			cp.APf("body", "    var %s = %s(%s)", this.genParamRefName(cursor, parent, aidx), this.tyconver.toDest(argty, cursor), strings.TrimRight(argdv, "f"))
		}
	} else if barety.Kind() == clang.Type_Typedef &&
		funk.Contains([]clang.TypeKind{clang.Type_Int, clang.Type_UInt, clang.Type_Long, clang.Type_LongLong, clang.Type_Double, clang.Type_UShort, clang.Type_UChar}, barety.Declaration().TypedefDeclUnderlyingType().Kind()) {
		if tmpdv, ok := argdvs[argdv]; ok {
			cp.APf("body", "    var %s = %s", this.genParamRefName(cursor, parent, aidx), tmpdv)
		} else {
			cp.APf("body", "    var %s = %s(%s)", this.genParamRefName(cursor, parent, aidx), this.tyconver.toDest(barety.Declaration().TypedefDeclUnderlyingType(), cursor), argdv)
		}
	} else if funk.Contains([]clang.TypeKind{clang.Type_Bool}, argty.Kind()) {
		cp.APf("body", "    var %s = %s", this.genParamRefName(cursor, parent, aidx), argdv)
	} else if funk.Contains([]clang.TypeKind{clang.Type_Char_S}, argty.Kind()) {
		cp.APf("body", "    var %s = %s", this.genParamRefName(cursor, parent, aidx), argdv)
	} else if TypeIsBoolPtr(argty) || TypeIsVoidPtr(argty) || TypeIsIntPtr(argty) || TypeIsUCharPtr(argty) {
		cp.APf("body", "    var %s : unsafe.Pointer", this.genParamRefName(cursor, parent, aidx))
	} else if TypeIsQFlags(argty) {
		cp.APf("body", "    var %s = 0", this.genParamRefName(cursor, parent, aidx))
	} else if is_qt_class(argty) &&
		funk.ContainsString([]string{"QString", "QByteArray", "QVariant", "QModelIndex", "QUrl",
			"QSize", "QAbstractState" /*"QScreen", "QAction"*/}, get_bare_type(argty).Spelling()) {
		usemod := get_decl_mod(cursor)
		pkgPref := gopp.IfElseStr(usemod == "core", "", "qtcore.")
		cp.APf("body", "    var convArg%d = %sNew%s()", aidx, pkgPref, get_bare_type(argty).Spelling())
	} else if is_qt_class(argty) && get_bare_type(argty).Spelling() == "QChar" {
		usemod := get_decl_mod(cursor)
		pkgPref := gopp.IfElseStr(usemod == "core", "", "qtcore.")
		cp.APf("body", "    var convArg%d  = %sNewQChar8('%s'.byte)", aidx,
			pkgPref, strings.Split(argdv, "'")[1])
	} else if is_qt_class(argty) && !isPrimitiveType(argty.CanonicalType()) {
		if argty.Spelling() == "QRgb" {
			log.Fatalln(argty.Spelling(), argty.CanonicalType().Kind().String())
		}
		refmod := get_decl_mod(argty.PointeeType().Declaration())
		usemod := get_decl_mod(cursor)
		log.Println("kkkkk", refmod, usemod, parent.Spelling())
		if _, ok := privClasses[argty.PointeeType().Spelling()]; ok {
		} else if usemod == "core" && refmod == "widgets" {
			cp.APf("body", "    var %s : pointer", this.genParamRefName(cursor, parent, aidx))
		} else if usemod == "gui" && refmod == "widgets" {
			cp.APf("body", "    var %s : pointer", this.genParamRefName(cursor, parent, aidx))
		} else {
			cp.APf("body", "    var convArg%d : pointer", aidx)
		}
	} else if argty.Spelling() == "WId" {
		cp.APf("body", "    var %s : pointer", this.genParamRefName(cursor, parent, aidx))
	} else if barety.Kind() == clang.Type_Typedef && TypeIsFuncPointer(undty) {
		cp.APf("body", "    var %s : pointer ", this.genParamRefName(cursor, parent, aidx))
	} else { // no convert needed
		// log.Fatalln("wtf", argty.Kind(), argty.Spelling(), parent.Spelling())
		cp.APf("body", "    #// var %s : pointer #// 111", this.genParamRefName(cursor, parent, aidx))
	}
}

func (this *GenerateNim) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateNim) genParam(cursor, parent clang.Cursor, aidx int) {
	argName := cursor.Spelling()
	argName = gopp.IfElseStr(is_nim_keyword(argName), argName+"_e", argName)
	this.paramDesc = append(this.paramDesc,
		gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName))
}

func (this *GenerateNim) genParamRefName(cursor, _ clang.Cursor, aidx int) string {
	argName := cursor.Spelling()
	if strings.HasPrefix(argName, "_") {
		argName = "e" + argName
	}
	if strings.HasSuffix(argName, "_") {
		argName += "e"
	}
	argName = gopp.IfElseStr(is_nim_keyword(argName), argName+"_e", argName)

	return gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName)
}

func (this *GenerateNim) genParamsFFI(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamFFI(argc, cursor, idx)
	}
}

func (this *GenerateNim) genParamFFI(cursor, parent clang.Cursor, idx int) {
	argty := cursor.Type()
	if TypeIsCharPtrPtr(argty) {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
	} else if TypeIsCharPtr(argty) {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
	} else if is_qt_class(argty) && get_bare_type(argty).Spelling() == "QString" {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
	} else if is_qt_class(argty) && !isPrimitiveType(argty.CanonicalType()) {
		usemod := get_decl_mod(cursor)
		refmod := get_decl_mod(argty.PointeeType().Declaration())
		if _, ok := privClasses[argty.PointeeType().Spelling()]; ok {
		} else if usemod == "core" && refmod == "widgets" {
			this.paramDesc = append(this.paramDesc, cursor.Spelling())
		} else if usemod == "gui" && refmod == "widgets" {
			this.paramDesc = append(this.paramDesc, cursor.Spelling())
		} else {
			this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
		}
	} else {
		argName := cursor.Spelling()
		argName = gopp.IfElseStr(is_nim_keyword(argName), argName+"_e", argName)

		useand := argty.Kind() == clang.Type_LValueReference &&
			isPrimitiveType(argty.PointeeType())
		if argty.Kind() == clang.Type_Pointer && isPrimitiveType(argty.PointeeType()) &&
			argty.PointeeType().Kind() == clang.Type_UChar { // UChar, SChar是字符串或者字节串
			useand = false
		} else if argty.Kind() == clang.Type_Pointer && isPrimitiveType(argty.PointeeType()) &&
			argty.PointeeType().Kind() == clang.Type_Bool {
			useand = false
		}
		andop := gopp.IfElseStr(useand, "&", "")
		this.paramDesc = append(this.paramDesc,
			andop+gopp.IfElseStr(cursor.Spelling() == "",
				fmt.Sprintf("arg%d", idx), fmt.Sprintf("%s", argName)))
	}
}

func (this *GenerateNim) genRetFFI(cursor, parent clang.Cursor, midx int) {
	var cp = this.getpropercp(cursor)

	rety := cursor.ResultType()
	retybare := get_bare_type(rety.CanonicalType()).Declaration()
	defmod := get_decl_mod(retybare)
	if retybare.Spelling() == "QList" {
		defmod = get_decl_mod(rety.Declaration())
		if defmod == "stdglobal" {
			if strings.Contains(rety.Spelling(), "QObjectList") {
				defmod = "core"
			}
		}
		if strings.Contains(rety.Spelling(), "QCameraInfo") {
			defmod = "multimedia"
		} else if strings.Contains(rety.Spelling(), "QGraphicsItem") {
			defmod = "widgets"
		} else if strings.Contains(rety.Spelling(), "QQuickItem") {
			defmod = "quick"
		}
	}
	usemod := get_decl_mod(cursor)
	log.Println("hhhhh use ==? ref", retybare.Spelling(), defmod, usemod, rety.Spelling(), cursor.DisplayName(), parent.Spelling())
	pkgPrefix := gopp.IfElseStr(defmod == usemod, "", fmt.Sprintf("qt%s.", defmod))

	switch rety.Kind() {
	case clang.Type_Void:
	case clang.Type_Int, clang.Type_UInt, clang.Type_Long, clang.Type_ULong,
		clang.Type_Short, clang.Type_UShort,
		clang.Type_Char_S, clang.Type_Char_U, clang.Type_UChar,
		clang.Type_Float, clang.Type_Double, clang.Type_LongDouble:
		cp.APf("body", "    return qtrt.Cretval2go[%s](\"%s\", rv) #// 1111",
			this.tyconver.toDest(rety, cursor), this.tyconver.toDest(rety, cursor))
		// cp.APf("body", "    return %s(rv) // 111", this.tyconver.toDest(rety, cursor))
	case clang.Type_Typedef:
		if TypeIsQFlags(rety) {
			cp.APf("body", "    return int(rv)")
		} else if is_qt_class(rety.CanonicalType()) &&
			(rety.Spelling() == "QObjectList" || rety.Spelling() == "QModelIndexList" ||
				rety.Spelling() == "QFileInfoList" || rety.Spelling() == "QVariantList" ||
				(TypeIsConsted(rety) && (strings.HasSuffix(rety.Spelling(), "QVariantList"))) ||
				rety.Spelling() == "QWindowList" || rety.Spelling() == "QWidgetList" ||
				rety.Spelling() == "QCameraFocusZoneList" || rety.Spelling() == "QMediaResourceList") {
			if strings.HasPrefix(rety.Spelling(), "QWidget") || strings.HasPrefix(rety.Spelling(), "QGraphicsItem") {
				pkgPrefix = "/*222*/"
			}
			cp.APf("body", "    var rv2 = %sNew%sFromPointer(cast[pointer](rv)) #//5551",
				pkgPrefix, gopp.IfElseStr(TypeIsConsted(rety), rety.Spelling()[6:], rety.Spelling()))
			cp.APf("body", "    return rv2")
		} else if is_qt_class(rety.CanonicalType()) {
			cp.APf("body", "    var rv2 = %sNew%sFromPointer(cast[pointer](rv)) #//555",
				// pkgPrefix, rety.Spelling())
				pkgPrefix, get_bare_type(rety.CanonicalType()).Spelling())
			cp.APf("body", "    return rv2")
		} else if TypeIsFuncPointer(rety.CanonicalType()) {
			cp.APf("body", "    return cast[pointer](rv)")
		} else if rety.Spelling() == "qreal" {
			cp.APf("body", "    return qtrt.Cretval2go[float64](\"%s\", rv) #// 1111",
				this.tyconver.toDest(rety, cursor), this.tyconver.toDest(rety, cursor))
		} else if TypeIsCharPtr(rety.CanonicalType()) {
			cp.APf("body", "    return qtrt.GoStringI(rv)")
			// TODO iterator is pointer, don't convert to string
		} else if TypeIsPtr(rety.CanonicalType()) {
			cp.APf("body", "    return cast[pointer](rv)")
		} else if TypeIsIter(rety.CanonicalType()) {
			cp.APf("body", "    return cast[pointer](rv)")
		} else if strings.HasPrefix(this.tyconver.toDest(rety, cursor), "unsafe.Pointer") {
			cp.APf("body", "    return cast[pointer](rv)")
		} else {
			cp.APf("body", "    return %s(rv) #// 222", this.tyconver.toDest(rety, cursor))
		}
	case clang.Type_Record:
		if is_qt_class(rety) && get_bare_type(rety).Spelling() == "QString" {
			cp.APf("body", "    var rv2 = %sNewQStringFromPointer(cast[pointer](rv))", pkgPrefix)
			cp.APf("body", "    var rv3 = rv2.ToUtf8().Data()")
			cp.APf("body", "    %sDeleteQString(rv2)", pkgPrefix)
			cp.APf("body", "    return rv3")
		} else if is_qt_class(rety) {
			barety := get_bare_type(rety)
			cp.APf("body", "    var rv2 = %sNew%sFromPointer(cast[pointer](rv)) #// 333",
				pkgPrefix, barety.Spelling())
			cp.APf("body", "    qtrt.SetFinalizer(rv2, %sDelete%s)", pkgPrefix, barety.Spelling())
			cp.APf("body", "    return rv2")
		} else {
			cp.APf("body", "    return cast[pointer](rv)")
		}

	case clang.Type_LValueReference:
		if is_qt_class(rety) && get_bare_type(rety).Spelling() == "QString" {
			cp.APf("body", "    var rv2 = %sNewQStringFromPointer(cast[pointer](rv))", pkgPrefix)
			cp.APf("body", "    var rv3 = rv2.ToUtf8().Data()")
			cp.APf("body", "    %sDeleteQString(rv2)", pkgPrefix)
			cp.APf("body", "    return rv3")
		} else if is_qt_class(rety) {
			barety := get_bare_type(rety)
			cp.APf("body", "    var rv2 = %sNew%sFromPointer(cast[pointer](rv)) #// 4441",
				pkgPrefix, barety.Spelling())
			cp.APf("body", "    qtrt.SetFinalizer(rv2, %sDelete%s)", pkgPrefix, barety.Spelling())
			cp.APf("body", "    return rv2")
		} else if TypeIsCharPtr(rety) {
			cp.APf("body", "    return qtrt.GoStringI(rv)")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UChar {
			cp.APf("body", "    return byte(rv) /*2221*/")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UShort {
			cp.APf("body", "    return uint16(rv)")
		} else if isPrimitiveType(rety.PointeeType()) {
			// int(*(*C.int)(unsafe.Pointer(uintptr(rv))))
			cp.APf("body", "    return qtrt.Cpretval2go[%s](\"%s\", rv) #// 3331",
				this.tyconver.toDest(rety.PointeeType(), cursor),
				this.tyconver.toDest(rety.PointeeType(), cursor))
			// this.cp.APf("body", "    return %s(rv) // 3331", this.tyconver.toDest(rety.PointeeType(), cursor))
		} else {
			cp.APf("body", "    return cast[pointer](rv)")
		}
	case clang.Type_Pointer:
		if is_qt_class(rety) && get_bare_type(rety).Spelling() == "QString" {
			cp.APf("body", "    var rv2 = %sNewQStringFromPointer(cast[pointer](rv))", pkgPrefix)
			cp.APf("body", "    var rv3 = rv2.ToUtf8().Data()")
			cp.APf("body", "    %sDeleteQString(rv2)", pkgPrefix)
			cp.APf("body", "    return rv3")
		} else if is_qt_class(rety) {
			if _, ok := privClasses[rety.PointeeType().Spelling()]; ok {
				cp.APf("body", "    return cast[pointer](rv)")
			} else if usemod == "core" && defmod == "widgets" {
				cp.APf("body", "    return cast[pointer](rv)")
			} else if usemod == "gui" && defmod == "widgets" {
				cp.APf("body", "    return cast[pointer](rv)")
			} else {
				barety := get_bare_type(rety)
				cp.APf("body", "    return %sNew%sFromPointer(cast[pointer](rv)) #// 444",
					pkgPrefix, barety.Spelling())
			}
		} else if TypeIsCharPtrPtr(rety) {
			cp.APf("body", "    return qtrt.CCharPPToStringSlice(cast[pointer](rv))")
		} else if TypeIsCharPtr(rety) {
			cp.APf("body", "    return qtrt.GoStringI(rv)")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UChar {
			cp.APf("body", "    return cast[pointer](rv)")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UShort {
			cp.APf("body", "    return cast[pointer](rv)")
		} else if isPrimitiveType(rety.PointeeType()) {
			cp.APf("body", "    return cast[pointer](rv)")
			// this.cp.APf("body", "    return %s(rv) // 333", this.tyconver.toDest(rety.PointeeType(), cursor))
		} else {
			cp.APf("body", "    return cast[pointer](rv)")
		}
	case clang.Type_RValueReference:
		cp.APf("body", "    return cast[pointer](rv) #//777")
	case clang.Type_Bool:
		cp.APf("body", "    return rv!=0")
	case clang.Type_Enum:
		cp.APf("body", "    return int(rv)")
	case clang.Type_Elaborated:
		cp.APf("body", "    return int(rv)")
	case clang.Type_Unexposed:
		if strings.HasPrefix(rety.Spelling(), "QList<") {
			cp.APf("body", "    var rv2 = %sNew%sListFromPointer(cast[pointer](rv)) #//5552",
				pkgPrefix, strings.TrimRight(rety.Spelling()[6:], " *>"))
			cp.APf("body", "    return rv2")
		} else {
			cp.APf("body", "    return rv #/*-222*/")
		}
	default:
		cp.APf("body", "    return rv #/*-111*/")
	}
}

func (this *GenerateNim) genArgsCGO(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgCGO(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.origin(cursor))
}

func (this *GenerateNim) genArgCGO(cursor, parent clang.Cursor, idx int) {
	argty := cursor.Type()
	argName := gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", idx), cursor.Spelling())
	argName = gopp.IfElseStr(is_nim_keyword(argName), argName+"_e", argName)

	dstr := getTyDesc(argty, ArgTyDesc_CGO_SIGNATURE, cursor)
	this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s", argName, dstr))
}

func (this *GenerateNim) genArgsCGOSign(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgCGOSign(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.origin(cursor))
}

func (this *GenerateNim) genArgCGOSign(cursor, parent clang.Cursor, idx int) {
	argty := cursor.Type()
	argName := gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", idx), cursor.Spelling())
	argName = gopp.IfElseStr(is_nim_keyword(argName), argName+"_e", argName)

	tystr := getTyDesc(argty, ArgTyDesc_C_SIGNATURE_USED_IN_CGO_EXTERN, cursor)
	this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s", tystr, argName))
}

func (this *GenerateNim) genClassEnums(cursor, parent clang.Cursor) {
	// log.Println("yyyyyyyy", cursor.DisplayName(), parent.DisplayName())
	isobjty := has_qobject_base_class(cursor)
	for _, enum := range this.enums {
		comment := queryComment(enum, this.qtdir, this.qtver)
		pcomment, elems := extractEnumElem(comment)
		_ = pcomment
		this.cp.APf("body", "")
		this.cp.APf("body", "#[")
		//this.cp.APf("body", "%s", pcomment)
		this.cp.APf("body", "]#")
		// must use uint, because on android
		this.cp.APf("body", "type %s_Enum_%s = int", cursor.DisplayName(), enum.DisplayName())
		enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
			switch c1.Kind() {
			case clang.Cursor_EnumConstantDecl:
				log.Println("yyyyyyyyy", c1.EnumConstantDeclValue(), c1.DisplayName(), p1.DisplayName(), cursor.DisplayName())
				this.cp.APf("body", "#// %s", elems[c1.DisplayName()])
				this.cp.APf("body", "const %s_Enum_%s : %s_Enum_%s = %d",
					cursor.DisplayName(), c1.DisplayName(),
					cursor.DisplayName(), p1.DisplayName(),
					c1.EnumConstantDeclValue())
			}

			return clang.ChildVisit_Continue
		})

		// generate get enum item name by enum value
		revalmap := map[int64][]string{} // reverse enum val => enum names
		enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
			switch c1.Kind() {
			case clang.Cursor_EnumConstantDecl:
				eival := c1.EnumConstantDeclValue()
				if _, ok := revalmap[eival]; ok {
					revalmap[eival] = append(revalmap[eival], c1.DisplayName())
				} else {
					revalmap[eival] = []string{c1.DisplayName()}
				}
			}

			return clang.ChildVisit_Continue
		})

		this.cp.APf("body", "proc %sItemName*(this:%s, val: int) : string =",
			enum.DisplayName(), cursor.DisplayName())
		if isobjty {
			this.cp.APf("body", "  return qtrt.GetClassEnumItemName(this, val)")
		} else {
			this.cp.APf("body", "  case val:")
			enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
				switch c1.Kind() {
				case clang.Cursor_EnumConstantDecl:
					eival := c1.EnumConstantDeclValue()
					_, keyok := revalmap[eival]
					commentit := gopp.IfElseStr(keyok, "", "#//")
					this.cp.APf("body", "    %sof %s_Enum_%s: #// %d",
						commentit, cursor.DisplayName(), c1.DisplayName(), eival)
					this.cp.APf("body", "      %sreturn \"%s\"", commentit, strings.Join(revalmap[eival], ","))
					if keyok {
						delete(revalmap, eival)
					}
				}
				return clang.ChildVisit_Continue
			})
			this.cp.APf("body", "    else: return \"$#\" %% [$val]")
			this.cp.APf("body", "#}")
		}
		this.cp.APf("body", "#}")
		this.cp.APf("body", "func %s_%sItemName*(val : int) : string =",
			cursor.DisplayName(), enum.DisplayName())
		this.cp.APf("body", "  var nilthis : %s", cursor.DisplayName())
		this.cp.APf("body", "  return nilthis.%sItemName(val)", enum.DisplayName())
		this.cp.APf("body", "#}")
		this.cp.APf("body", "")
	}
}

// enum一定要使用int类型，而不能用uint。注意-1值的处理
func (this *GenerateNim) genEnumsGlobal(cursor, parent clang.Cursor) {
	// log.Println("yyyyyyyy", cursor.DisplayName(), parent.DisplayName())
	dedups := map[string]int{}
	for _, enum := range this.enums {
		if enum.DisplayName() == "" || enum.DisplayName() == "Uninitialized" ||
			enum.DisplayName() == "timeout" || enum.DisplayName() == "deferred" ||
			enum.DisplayName() == "GuardValues" || enum.DisplayName() == "cv_status" ||
			enum.DisplayName() == "future_statu" || enum.DisplayName() == "launch" {
			continue
		}
		if _, ok := dedups[enum.DisplayName()]; ok {
			continue
		}
		dedups[enum.DisplayName()] = 1

		comment := queryComment(enum, this.qtdir, this.qtver)
		pcomment, elems := extractEnumElem(comment)
		_ = pcomment
		qtmod := get_decl_mod(enum)
		this.cp.APf("body", "")
		this.cp.APf("body", "#[")
		//this.cp.APf("body", "%s", pcomment)
		this.cp.APf("body", "]#")
		this.cp.APUf("body", "type %s_Enum_%s = int #// %s", "Qt", enum.DisplayName(), qtmod)
		enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
			switch c1.Kind() {
			case clang.Cursor_EnumConstantDecl:
				log.Println("yyyyyyyyy", c1.EnumConstantDeclValue(), c1.DisplayName(), p1.DisplayName(), cursor.DisplayName())
				if _, ok := dedups[c1.DisplayName()]; ok {
					break
				}
				dedups[c1.DisplayName()] = 1

				this.cp.APUf("body", "#// %s", elems[c1.DisplayName()])
				this.cp.APUf("body", "const %s_Enum_%s : %s_Enum_%s = %d",
					"Qt", c1.DisplayName(), "Qt", p1.DisplayName(),
					c1.EnumConstantDeclValue())
			}

			return clang.ChildVisit_Continue
		})

		// generate get enum item name by enum value
		revalmap := map[int64][]string{} // reverse enum val => enum names
		enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
			switch c1.Kind() {
			case clang.Cursor_EnumConstantDecl:
				eival := c1.EnumConstantDeclValue()
				if _, ok := revalmap[eival]; ok {
					revalmap[eival] = append(revalmap[eival], c1.DisplayName())
				} else {
					revalmap[eival] = []string{c1.DisplayName()}
				}
			}

			return clang.ChildVisit_Continue
		})

		this.cp.APf("body", "proc %sItemName*(val : int) : string =", enum.DisplayName())
		this.cp.APf("body", "  case val:")
		enum.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
			switch c1.Kind() {
			case clang.Cursor_EnumConstantDecl:
				eival := c1.EnumConstantDeclValue()
				_, keyok := revalmap[eival]
				commentit := gopp.IfElseStr(keyok, "", "#//")
				this.cp.APf("body", "    %sof %s_Enum_%s: #// %d", commentit, "Qt", c1.DisplayName(), eival)
				this.cp.APf("body", "      %sreturn \"%s\"", commentit, strings.Join(revalmap[eival], ","))
				if keyok {
					delete(revalmap, eival)
				}
			}
			return clang.ChildVisit_Continue
		})
		this.cp.APf("body", "    else: return \"$#\" %% [$val]")
		this.cp.APf("body", "#}")
		this.cp.APf("body", "#}")
		this.cp.APf("body", "")

	}
}

func (this *GenerateNim) genEnum() {

}

func (this *GenerateNim) genFunctions(cursor clang.Cursor, parent clang.Cursor) {
	// this.genHeader(cursor, parent)
	skipKeys := []string{"QKeySequence", "QVector2D", "QPointingDeviceUniqueId", "QFont", "QMatrix",
		"QTransform", "QPixelFormat", "QRawFont", "QVector3D", "QVector4D",
		"QOpenGLVersionStatus", "QOpenGLVersionProfile"}
	hasSkipKey := func(c clang.Cursor) bool {
		for _, k := range skipKeys {
			if strings.Contains(c.DisplayName(), k) {
				return true
			}
		}
		return false
	}

	grfuncs := this.groupFunctionsByModule()
	qtmods := []string{}
	for qtmod, _ := range grfuncs {
		qtmods = append(qtmods, qtmod)
	}
	sort.Strings(qtmods)

	for _, qtmod := range qtmods {
		funcs := grfuncs[qtmod]
		log.Println(qtmod, len(funcs))
		this.cp = NewCodePager()
		this.cpnomin = NewCodePager()

		// write code
		writehead := func(cp *CodePager) {
			cp.APf("header", "package qt%s", qtmod)
			cp.APf("header", "import \"unsafe\"")
			cp.APf("header", "import \"github.com/kitech/qt.go/qtrt\"")
			for _, mod := range modDeps[qtmod] {
				cp.APf("header", "import \"github.com/kitech/qt.go/qt%s\"", mod)
			}
			cp.APf("header", "func init_unused_%d(){", this.nextclsidx())
			cp.APf("header", "  if false{_=unsafe.Pointer(uintptr(0))}")
			cp.APf("header", "  if false{qtrt.KeepMe()}")
			cp.APf("header", "  if false{qtrt.KeepMe()}")
			for _, dep := range modDeps[qtmod] {
				cp.APf("header", "if false {qt%s.KeepMe()}", dep)
			}
			cp.APf("header", "}")
		}
		writehead(this.cp)
		writehead(this.cpnomin)

		// 这个是一个包范围内的排序还是所有包范围内的排序呢？
		sort.Slice(funcs, func(i int, j int) bool {
			return funcs[i].Mangling() > funcs[j].Mangling()

		})
		for _, fc := range funcs {
			log.Println(fc.Spelling(), fc.Mangling(), fc.DisplayName(), fc.IsCursorDefinition(), is_qt_global_func(fc))
			if !is_qt_global_func(fc) {
				log.Println("skip global function ", fc.Spelling())
				continue
			}

			if strings.ContainsAny(fc.DisplayName(), "<>") {
				log.Println("skip global function ", fc.Spelling())
				continue
			}
			if strings.Contains(fc.DisplayName(), "Rgba64") {
				log.Println("skip global function ", fc.Spelling())
				continue
			}
			if strings.Contains(fc.ResultType().Spelling(), "Rgba64") {
				log.Println("skip global function ", fc.Spelling())
				continue
			}
			if hasSkipKey(fc) {
				log.Println("skip global function ", fc.Spelling())
				continue
			}

			if this.filter.skipFunc(fc) {
				log.Println("skip global function ", fc.Spelling())
				continue
			}

			if _, ok := this.funcMangles[fc.Spelling()]; ok {
				this.funcMangles[fc.Spelling()] += 1
			} else {
				this.funcMangles[fc.Spelling()] = 0
			}
			olidx := this.funcMangles[fc.Spelling()]
			log.Println("wtf ", qtmod, fc.Spelling())
			this.genFunction(fc, olidx)
		}

		this.saveCodeToFile(qtmod, "qfunctions")
	}
}

func (this *GenerateNim) genFunction(cursor clang.Cursor, olidx int) {
	this.genParamsFFI(cursor, cursor.SemanticParent())
	paramStr := strings.Join(this.paramDesc, ", ")
	_ = paramStr
	var cp = this.getpropercp(cursor)

	this.genMethodHeader(cursor, cursor.SemanticParent(), olidx)
	this.genBareFunctionSignature(cursor, cursor.SemanticParent(), olidx)

	this.genArgsConvFFI(cursor, cursor.SemanticParent(), olidx)
	cp.APf("body", "  var (rv, err) = qtrt.InvokeQtFunc6(\"%s\", qtrt.FFI_TYPE_POINTER, %s)",
		cursor.Mangling(), paramStr)
	cp.APf("body", "  qtrt.ErrPrint(err, rv)")

	this.genRetFFI(cursor, cursor.SemanticParent(), olidx)
	this.genMethodFooterFFI(cursor, cursor.SemanticParent(), olidx)
	cp.APf("body", "")
}

// only for static member
func (this *GenerateNim) genBareFunctionSignature(cursor, parent clang.Cursor, midx int) {
	this.genArgsDest(cursor, parent, true)
	argStr := strings.Join(this.destArgDesc, ", ")
	if strings.Contains(argStr, "DropActions::enum_type") {
		log.Fatalln(parent.Spelling(), cursor.DisplayName(), cursor.Spelling(), argStr)
	}
	var cp = this.getpropercp(cursor)

	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("%d", midx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
	case clang.Cursor_Destructor:
	default:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if is_qstring_cls(retPlace) {
			retPlace = "string"
		}
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = ""
		}
		retPlace = strings.TrimLeft(retPlace, "*")
		retPlace = gopp.IfElseStr(retPlace == "", "", ":"+retPlace)
		rewname := rewriteOperatorMethodName(cursor.Spelling())
		cp.APf("body", "proc %s%s*(%s) %s {",
			strings.Title(rewname), overloadSuffix, argStr, retPlace)
	}
}

// TODO sperate to modules
func (this *GenerateNim) genConstantsGlobal(cursor, parent clang.Cursor) {
	dedups := map[string]int{}
	for _, macro := range this.constants {
		if strings.HasPrefix(macro.Spelling(), "_") {
			continue
		}
		qtmod := get_decl_mod(macro)
		if qtmod == "stdglobal" {
			continue
		}
		macroval, macroty := readDefineRange(macro.Extent())
		if macroty == "" {
			continue
		}
		if strings.ContainsAny(macroval, "()\\") {
			continue
		}
		if _, ok := dedups[macro.Spelling()]; ok {
			continue
		}
		dedups[macro.Spelling()] = 1

		macroval = gopp.IfElseStr(strings.HasPrefix(macroty, "num"), strings.TrimRight(macroval, "ACDL"), macroval)

		log.Println(qtmod, macro.Spelling(), macroval, macroty)
		this.cp.APf("body", "const %s = %s #// %s @ %s", macro.Spelling(), macroval, macroty, qtmod)
	}
}

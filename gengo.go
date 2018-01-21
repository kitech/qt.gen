package main

import (
	"fmt"
	"gopp"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

type GenerateGo struct {
	// TODO move to base
	filter   GenFilter
	mangler  GenMangler
	tyconver TypeConvertor

	maxClassSize int64 // 暂存一下类的大小的最大值

	methods     []clang.Cursor
	cp          *CodePager
	argDesc     []string // origin c/c++ language syntax
	paramDesc   []string
	destArgDesc []string // dest language syntax
}

func NewGenerateGo() *GenerateGo {
	this := &GenerateGo{}
	this.filter = &GenFilterGo{}
	this.mangler = NewGoMangler()
	this.tyconver = NewTypeConvertGo()

	this.cp = NewCodePager()
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
		// this.cp.APf(block, "// block begin--- %s", block)
	}

	return this
}

func (this *GenerateGo) genClass(cursor, parent clang.Cursor) {
	if false {
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	}
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	this.genHeader(cursor, parent)
	this.walkClass(cursor, parent)
	// this.genExterns(cursor, parent)
	this.genImports(cursor, parent)
	this.genClassDef(cursor, parent)
	this.genMethods(cursor, parent)
	this.final(cursor, parent)
	if cursor.Spelling() == "QMimeType" {

	}

}

func (this *GenerateGo) final(cursor, parent clang.Cursor) {
	// log.Println(this.cp.ExportAll())
	this.saveCode(cursor, parent)

	this.cp = NewCodePager()
}
func (this *GenerateGo) saveCode(cursor, parent clang.Cursor) {
	// qtx{yyy}, only yyy
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]
	savefile := fmt.Sprintf("src/%s/%s.go", modname, strings.ToLower(cursor.Spelling()))

	// TODO gofmt the code
	// log.Println(this.cp.AllPoints())
	bcc := this.cp.ExportAll()
	if strings.HasPrefix(bcc, "//") {
		bcc = bcc[strings.Index(bcc, "\n"):]
	}
	ioutil.WriteFile(savefile, []byte(bcc), 0644)
	cmd := exec.Command("/usr/bin/gofmt", []string{"-w", savefile}...)
	err := cmd.Run()
	if err != nil {
		log.Println(err, cmd)
	}
}

func (this *GenerateGo) genHeader(cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	this.cp.APf("header", "package %s", strings.ToLower(fullModname[0:]))
	this.cp.APf("header", "// %s", file.Name())
	this.cp.APf("header", "// #include <%s>", filepath.Base(file.Name()))
	this.cp.APf("header", "// #include <%s>", fullModname)
	this.cp.APf("header", "")
	this.cp.APf("ext", "")
	this.cp.APf("ext", "/*")
	this.cp.APf("ext", "#include <stdlib.h>")
	this.cp.APf("ext", "// extern C begin: %d", len(this.methods))
}

func (this *GenerateGo) walkClass(cursor, parent clang.Cursor) {

	methods := make([]clang.Cursor, 0)

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
				// log.Println("filtered:", cursor.Spelling())
			}
		case clang.Cursor_UnexposedDecl:
			// log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			file, line, col, _ := cursor.Location().FileLocation()
			if false {
				log.Println(file.Name(), line, col, file.Time())
			}
		default:
			if false {
				log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			}
		}
		return clang.ChildVisit_Continue
	})

	this.methods = methods
}

func (this *GenerateGo) genExterns(cursor, parent clang.Cursor) {
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
			this.cp.APf("ext", "extern void %s();", this.mangler.convTo(cursor))
		}
	}

	this.cp.APf("ext", "// extern C end: %d", len(this.methods))
}

func (this *GenerateGo) genImports(cursor, parent clang.Cursor) {

	this.cp.APf("ext", "*/")
	this.cp.APf("ext", "// import \"C\"") // 直接import "C"导致编译速度下降n倍

	file, _, _, _ := cursor.Location().FileLocation()
	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]

	this.cp.APf("ext", "import \"unsafe\"")
	this.cp.APf("ext", "import \"reflect\"")
	this.cp.APf("ext", "import \"fmt\"")
	this.cp.APf("ext", "import \"qtrt\"")
	this.cp.APf("ext", "import \"mkuse/cffiqt\"")
	this.cp.APf("ext", "import \"gopp\"")
	for _, dep := range modDeps[modname] {
		this.cp.APf("ext", "import \"qt%s\"", dep)
	}

	this.cp.APf("ext", "")
	this.cp.APf("ext", "func init() {")
	this.cp.APf("ext", "  if false {reflect.TypeOf(123)}")
	this.cp.APf("ext", "  if false {reflect.TypeOf(unsafe.Sizeof(0))}")
	this.cp.APf("ext", "  if false {fmt.Println(123)}")
	this.cp.APf("ext", "  if false {qtrt.KeepMe()}")
	this.cp.APf("ext", "  if false {ffiqt.KeepMe()}")
	this.cp.APf("ext", "  if false {gopp.KeepMe()}")
	for _, dep := range modDeps[modname] {
		this.cp.APf("ext", "if false {qt%s.KeepMe()}", dep)
	}
	this.cp.APf("ext", "}")
}

func (this *GenerateGo) genClassDef(cursor, parent clang.Cursor) {
	bcs := find_base_classes(cursor)
	bcs = this.filter_base_classes(bcs)

	this.cp.APf("body", "type %s struct {", cursor.Spelling())
	if len(bcs) == 0 {
		this.cp.APf("body", "    *qtrt.CObject")
	} else {
		for _, bc := range bcs {
			this.cp.APf("body", "    *%s%s", calc_package_suffix(cursor, bc), bc.Type().Spelling())
			// break // TODO multiple base class
		}
	}
	// this.cp.APf("body", "    cthis unsafe.Pointer")
	this.cp.APf("body", "}")

	this.genGetCthis(cursor, cursor, 0) // 只要定义了结构体，就有GetCthis方法
	this.genCtorFromPointer(cursor, cursor, 0)
}

func (this *GenerateGo) filter_base_classes(bcs []clang.Cursor) []clang.Cursor {
	newbcs := make([]clang.Cursor, 0)
	for _, bc := range bcs {
		if !this.filter.skipClass(bc, bc.SemanticParent()) {
			newbcs = append(newbcs, bc)
		}
	}
	return newbcs
}

func (this *GenerateGo) genMethods(cursor, parent clang.Cursor) {
	// log.Println("process class:", len(this.methods), cursor.Spelling())
	grpMethods := this.groupMethods()
	// log.Println(len(grpMethods))

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
			// log.Println(idx, cursor.Kind().String(), cursor.DisplayName())
			switch cursor.Kind() {
			case clang.Cursor_Constructor:
				this.genCtor(cursor, parent, idx)
				// this.genCtorFromPointer(cursor, parent, idx)
				// this.genGetCthis(cursor, parent, idx)
			case clang.Cursor_Destructor:
				this.genDtor(cursor, parent, idx)
			default:
				if cursor.CXXMethod_IsStatic() {
					this.genStaticMethod(cursor, parent, idx)
					this.genStaticMethodNoThis(cursor, parent, idx)
				} else {
					this.genNonStaticMethod(cursor, parent, idx)
				}
			}
		}

		// this.genMethodFooter(cursors[0], cursors[0].SemanticParent())
	}
}

// 按名字分组
func (this *GenerateGo) groupMethods() [][]clang.Cursor {
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

func (this *GenerateGo) genMethodHeader(cursor, parent clang.Cursor, midx int) {
	file, lineno, _, _ := cursor.Location().FileLocation()
	this.cp.APf("body", "// %s:%d", file.Name(), lineno)
	this.cp.APf("body", "// index:%d", midx)

	qualities := make([]string, 0)
	qualities = append(qualities, strings.Split(cursor.AccessSpecifier().Spelling(), "=")[1])
	if cursor.CXXMethod_IsStatic() {
		qualities = append(qualities, "static")
	}
	if cursor.IsFunctionInlined() {
		qualities = append(qualities, "inline")
	}
	if cursor.CXXMethod_IsPureVirtual() {
		qualities = append(qualities, "pure")
	}
	if cursor.CXXMethod_IsVirtual() {
		qualities = append(qualities, "virtual")
	}

	if len(qualities) > 0 {
		this.cp.APf("body", "// %s", strings.Join(qualities, " "))
	}

	this.cp.APf("body", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())
}

func (this *GenerateGo) genMethodInit(cursor, parent clang.Cursor) {
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

func (this *GenerateGo) genMethodSignature(cursor, parent clang.Cursor, midx int) {
	if cursor.Kind() == clang.Cursor_Constructor {
	}

	this.genArgsDest(cursor, parent)
	argStr := strings.Join(this.destArgDesc, ", ")

	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("_%d", midx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		this.cp.APf("body", "func New%s%s(%s) *%s {",
			strings.Title(cursor.Spelling()),
			overloadSuffix, argStr, parent.Spelling())
	case clang.Cursor_Destructor:
		this.cp.APf("body", "func Delete%s%s(*%s) {",
			strings.Title(cursor.Spelling()[1:]),
			overloadSuffix, parent.Spelling())
	default:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = ""
		}
		this.cp.APf("body", "func (this *%s) %s%s(%s) %s {",
			parent.Spelling(), strings.Title(cursor.Spelling()),
			overloadSuffix, argStr, retPlace)
	}

	// TODO fill types, default args
}

// only for static member
func (this *GenerateGo) genMethodSignatureNoThis(cursor, parent clang.Cursor, midx int) {
	this.genArgsDest(cursor, parent)
	argStr := strings.Join(this.destArgDesc, ", ")

	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("_%d", midx))
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
	case clang.Cursor_Destructor:
	default:
		retPlace := "interface{}"
		retPlace = this.tyconver.toDest(cursor.ResultType(), cursor)
		if cursor.ResultType().Kind() == clang.Type_Void {
			retPlace = ""
		}
		this.cp.APf("body", "func %s_%s%s(%s) %s {",
			parent.Spelling(), strings.Title(cursor.Spelling()),
			overloadSuffix, argStr, retPlace)
	}

	// TODO fill types, default args
}

func (this *GenerateGo) genMethodFooter(cursor, parent clang.Cursor) {
	this.cp.APf("body", "  default:")
	this.cp.APf("body", "    qtrt.ErrorResolve(\"%s\", \"%s\", args)",
		parent.Spelling(), cursor.Spelling())
	this.cp.APf("body", "  } // end switch")
	this.cp.APf("body", "}")
}

func (this *GenerateGo) genMethodFooterFFI(cursor, parent clang.Cursor, midx int) {
	this.cp.APf("body", "}")
}

func (this *GenerateGo) genVTableTypes(cursor, parent clang.Cursor, midx int) {

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

func (this *GenerateGo) genNameLookup(cursor, parent clang.Cursor) {
	this.cp.AP("body", "")
	this.cp.AP("body", "  var matchedIndex = qtrt.SymbolResolve(args, vtys)")
	this.cp.AP("body", "  if false {fmt.Println(matchedIndex)}")
	this.cp.AP("body", "  switch matchedIndex {")
}

func (this *GenerateGo) genCtor(cursor, parent clang.Cursor, midx int) {
	// log.Println(this.mangler.convTo(cursor))
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)

	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	_, _ = argStr, paramStr
	if parent.Type().SizeOf() > this.maxClassSize {
		this.maxClassSize = parent.Type().SizeOf()
	}
	this.cp.APf("body", "    cthis := qtrt.Calloc(1, 256) // %d", parent.Type().SizeOf())
	this.genArgsConvFFI(cursor, parent, midx)
	this.cp.APf("body", "    rv, err := ffiqt.InvokeQtFunc6(\"%s\", ffiqt.FFI_TYPE_VOID, cthis, %s)",
		this.mangler.origin(cursor), paramStr)
	this.cp.APf("body", "    gopp.ErrPrint(err, rv)")
	// this.cp.APf("body", "    return &%s{qtrt.CObject{cthis}}", parent.Spelling())
	this.cp.APf("body", "    gothis := New%sFromPointer(cthis)", parent.Spelling())
	this.cp.APf("body", "    return gothis")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateGo) genCtorFromPointer(cursor, parent clang.Cursor, midx int) {
	if midx > 0 { // 忽略更多重载
		return
	}
	bcs := find_base_classes(parent)
	bcs = this.filter_base_classes(bcs)

	this.cp.APf("body", "func New%sFromPointer(cthis unsafe.Pointer) *%s {",
		cursor.Spelling(), cursor.Spelling())
	if len(bcs) == 0 {
		this.cp.APf("body", "    return &%s{&qtrt.CObject{cthis}}", cursor.Spelling())
	} else {
		bcobjs := []string{}
		for i, bc := range bcs {
			pkgSuff := calc_package_suffix(cursor, bc)
			this.cp.APf("body", "    bcthis%d := %sNew%sFromPointer(cthis)", i, pkgSuff, bc.Spelling())
			bcobjs = append(bcobjs, fmt.Sprintf("bcthis%d", i))
			// break // TODO multiple base classes
		}
		bcobjArgs := strings.Join(bcobjs, ", ")
		this.cp.APf("body", "    return &%s{%s}", parent.Spelling(), bcobjArgs)
	}
	this.cp.APf("body", "}")
}

func (this *GenerateGo) genGetCthis(cursor, parent clang.Cursor, midx int) {
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

func (this *GenerateGo) genDtor(cursor, parent clang.Cursor, midx int) {
	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)

	/*
		this.cp.APf("body", "	 case %d:", midx)
		this.cp.APf("body", "	   var cthis unsafe.Pointer = this.cthis")
		this.cp.APf("body", "	   // ffi not need this C.%s(cthis)", this.mangler.convTo(cursor))
		this.cp.APf("body", "      // %s", cursor.DisplayName())
		this.cp.APf("body", "	   this.cthis = nil")
	*/
	this.cp.APf("body", "    rv, err := ffiqt.InvokeQtFunc6(\"%s\", ffiqt.FFI_TYPE_VOID)",
		this.mangler.origin(cursor))
	this.cp.APf("body", "    gopp.ErrPrint(err, rv)")

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateGo) genNonStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParamsFFI(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	if len(argStr) > 0 {
		argStr = ", " + argStr
	}

	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)

	_, _ = argStr, paramStr

	/*
		this.cp.APf("body", "    // %d: (%s), (%s)", midx, argStr, paramStr)
		this.cp.APf("body", "    case %d: // (%s), (%s)", midx, argStr, paramStr)
		this.genArgsConv(cursor, parent, midx)
		this.cp.APf("body", "      // ffi not need this C.%s(%s)", this.mangler.convTo(cursor), paramStr)
		this.cp.APf("body", "      // %s", cursor.DisplayName())
	*/

	this.genArgsConvFFI(cursor, parent, midx)
	this.cp.APf("body", "    rv, err := ffiqt.InvokeQtFunc6(\"%s\", ffiqt.FFI_TYPE_POINTER, this.GetCthis(), %s)",
		this.mangler.origin(cursor), paramStr)
	this.cp.APf("body", "    gopp.ErrPrint(err, rv)")
	if cursor.ResultType().Kind() != clang.Type_Void {
		this.cp.APf("body", "   //  return rv")
	}
	this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateGo) genStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.genMethodHeader(cursor, parent, midx)
	this.genMethodSignature(cursor, parent, midx)

	_, _ = argStr, paramStr
	/*
		this.cp.APf("body", "    // %d: (%s), (%s)", midx, argStr, paramStr)
		this.cp.APf("body", "    case %d: // (%s), (%s)", midx, argStr, paramStr)
		this.genArgsConv(cursor, parent, midx)
		this.cp.APf("body", "      // ffi not need this C.%s(%s)", this.mangler.convTo(cursor), paramStr)
		this.cp.APf("body", "      // %s", cursor.DisplayName())
	*/

	this.cp.APf("body", "    rv, err := ffiqt.InvokeQtFunc6(\"%s\", ffiqt.FFI_TYPE_POINTER, %s)",
		this.mangler.origin(cursor), paramStr)
	this.cp.APf("body", "    gopp.ErrPrint(err, rv)")
	if cursor.ResultType().Kind() != clang.Type_Void {
		this.cp.APf("body", "    // return rv")
	}
	this.genRetFFI(cursor, parent, midx)

	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateGo) genStaticMethodNoThis(cursor, parent clang.Cursor, midx int) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	_, _ = argStr, paramStr

	// this.genMethodHeaderLongName(cursor, parent, midx)
	this.genMethodSignatureNoThis(cursor, parent, midx)
	overloadSuffix := gopp.IfElseStr(midx == 0, "", fmt.Sprintf("_%d", midx))

	// this.cp.APf("body", "    // %d: (%s), (%s)", midx, argStr, paramStr)
	this.cp.APf("body", "    var nilthis *%s", parent.Spelling())
	if cursor.ResultType().Kind() == clang.Type_Void {
		this.cp.APf("body", "    nilthis.%s%s(%s)",
			strings.Title(cursor.Spelling()), overloadSuffix, paramStr)
	} else {
		this.cp.APf("body", "    rv := nilthis.%s%s(%s)",
			strings.Title(cursor.Spelling()), overloadSuffix, paramStr)
		this.cp.APf("body", "    return rv")
	}
	// this.genRetFFI(cursor, parent, midx)
	this.genMethodFooterFFI(cursor, parent, midx)
}

func (this *GenerateGo) genNonVirtualMethod(cursor, parent clang.Cursor, midx int) {

}

func (this *GenerateGo) genArgs(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArg(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateGo) genArg(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))

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

func (this *GenerateGo) genArgsDest(cursor, parent clang.Cursor) {
	this.destArgDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgDest(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.destArgDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateGo) genArgDest(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))

	argName := cursor.Spelling()
	if is_go_keyword(argName) {
		argName = argName + "_"
	}

	destTy := this.tyconver.toDest(cursor.Type(), cursor)
	if len(cursor.Spelling()) == 0 {
		this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("arg%d %s", idx, destTy))
	} else {
		if cursor.Type().Kind() == clang.Type_LValueReference {
			// 转成指针
		}
		if strings.HasPrefix(cursor.Type().CanonicalType().Spelling(), "QFlags<") {
			this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s int", argName))
		} else {
			if cursor.Type().Kind() == clang.Type_IncompleteArray {
				this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s %s", argName, destTy))
			} else if cursor.Type().Kind() == clang.Type_ConstantArray {
				idx := strings.Index(cursor.Type().Spelling(), " [")
				this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s %s %s",
					cursor.Type().Spelling()[0:idx], argName, cursor.Type().Spelling()[idx+1:]))
			} else {
				this.destArgDesc = append(this.destArgDesc, fmt.Sprintf("%s %s", argName, destTy))
			}
		}
	}
}

// midx method index
func (this *GenerateGo) genArgsConv(cursor, parent clang.Cursor, midx int) {
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgConv(argc, cursor, midx, idx)
	}
}

// midx method index
// aidx method index
func (this *GenerateGo) genArgConv(cursor, parent clang.Cursor, midx, aidx int) {
	this.cp.APf("body", "	   var arg%d %s", aidx, this.tyconver.toCall(cursor.Type(), parent))
	this.cp.APf("body", "	   // if %d >= len(args) {", aidx)
	this.cp.APf("body", "	   //	  arg%d = defaultargx", aidx)
	this.cp.APf("body", "	   // } else {")
	this.cp.APf("body", "	   //	  arg%d = argx.toBind", aidx)
	this.cp.APf("body", "	   // }")
}

// midx method index
func (this *GenerateGo) genArgsConvFFI(cursor, parent clang.Cursor, midx int) {
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgConvFFI(argc, cursor, midx, idx)
	}
}

// midx method index
// aidx method index
func (this *GenerateGo) genArgConvFFI(cursor, parent clang.Cursor, midx, aidx int) {
	argty := cursor.Type()
	if TypeIsCharPtrPtr(argty) {
		this.cp.APf("body", "    var convArg%d = qtrt.StringSliceToCCharPP(%s)", aidx,
			this.genParamRet(cursor, parent, aidx))
	} else if TypeIsCharPtr(argty) {
		this.cp.APf("body", "    var convArg%d = qtrt.CString(%s)", aidx,
			this.genParamRet(cursor, parent, aidx))
		this.cp.APf("body", "    defer qtrt.FreeMem(convArg%d)", aidx)
	} else if is_qt_class(argty) && !isPrimitiveType(argty.CanonicalType()) {
		if argty.Spelling() == "QRgb" {
			log.Fatalln(argty.Spelling(), argty.CanonicalType().Kind().String())
		}
		refmod := get_decl_mod(argty.PointeeType().Declaration())
		usemod := get_decl_mod(cursor)
		log.Println("kkkkk", refmod, usemod, parent.Spelling())
		if usemod == "core" && refmod == "widgets" {
		} else if usemod == "gui" && refmod == "widgets" {
		} else {
			this.cp.APf("body", "    var convArg%d = %s.GetCthis()", aidx,
				this.genParamRet(cursor, parent, aidx))
		}
	}
}

func (this *GenerateGo) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateGo) genParam(cursor, parent clang.Cursor, aidx int) {
	argName := cursor.Spelling()
	if is_go_keyword(argName) {
		argName = argName + "_"
	}
	this.paramDesc = append(this.paramDesc,
		gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName))
}

func (this *GenerateGo) genParamRet(cursor, parent clang.Cursor, aidx int) string {
	argName := cursor.Spelling()
	if is_go_keyword(argName) {
		argName = argName + "_"
	}

	return gopp.IfElseStr(cursor.Spelling() == "", fmt.Sprintf("arg%d", aidx), argName)
}

func (this *GenerateGo) genParamsFFI(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamFFI(argc, cursor, idx)
	}
}

func (this *GenerateGo) genParamFFI(cursor, parent clang.Cursor, idx int) {
	argty := cursor.Type()
	if TypeIsCharPtrPtr(argty) {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
	} else if TypeIsCharPtr(argty) {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
	} else if is_qt_class(argty) && !isPrimitiveType(argty.CanonicalType()) {
		usemod := get_decl_mod(cursor)
		refmod := get_decl_mod(argty.PointeeType().Declaration())
		if usemod == "core" && refmod == "widgets" {
			this.paramDesc = append(this.paramDesc, cursor.Spelling())
		} else if usemod == "gui" && refmod == "widgets" {
			this.paramDesc = append(this.paramDesc, cursor.Spelling())
		} else {
			this.paramDesc = append(this.paramDesc, fmt.Sprintf("convArg%d", idx))
		}
	} else {
		argName := cursor.Spelling()
		if is_go_keyword(argName) {
			argName = argName + "_"
		}
		andop := gopp.IfElseStr(isPrimitiveType(argty) ||
			(argty.Kind() == clang.Type_LValueReference &&
				isPrimitiveType(argty.PointeeType())), "&", "")
		this.paramDesc = append(this.paramDesc,
			andop+gopp.IfElseStr(cursor.Spelling() == "",
				fmt.Sprintf("arg%d", idx), fmt.Sprintf("%s", argName)))
	}
}

func (this *GenerateGo) genRetFFI(cursor, parent clang.Cursor, midx int) {
	rety := cursor.ResultType()
	refmod := get_decl_mod(get_bare_type(rety.CanonicalType()).Declaration())
	usemod := get_decl_mod(cursor)
	log.Println("hhhhh use ==? ref", refmod, usemod, rety.Spelling(), cursor.DisplayName(), parent.Spelling())
	pkgSuff := gopp.IfElseStr(refmod == usemod, "/*==*/", fmt.Sprintf("qt%s.", refmod))

	switch rety.Kind() {
	case clang.Type_Void:
	case clang.Type_Int, clang.Type_UInt, clang.Type_Long, clang.Type_ULong,
		clang.Type_Short, clang.Type_UShort, clang.Type_Char_S, clang.Type_Char_U,
		clang.Type_Float, clang.Type_Double, clang.Type_UChar:
		this.cp.APf("body", "    return %s(rv) // 111", this.tyconver.toDest(rety, cursor))
	case clang.Type_Typedef:
		if TypeIsQFlags(rety) {
			this.cp.APf("body", "    return int(rv)")
		} else if is_qt_class(rety.CanonicalType()) {
			this.cp.APf("body", "    rv2 := %sNew%sFromPointer(unsafe.Pointer(uintptr(rv))) //555",
				pkgSuff, get_bare_type(rety.CanonicalType()).Spelling())
			this.cp.APf("body", "    return rv2")
		} else if TypeIsFuncPointer(rety.CanonicalType()) {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		} else {
			this.cp.APf("body", "    return %s(rv) // 222", this.tyconver.toDest(rety, cursor))
		}
	case clang.Type_Record:
		if is_qt_class(rety) {
			barety := get_bare_type(rety)
			this.cp.APf("body", "    rv2 := %sNew%sFromPointer(unsafe.Pointer(uintptr(rv))) // 333",
				pkgSuff, barety.Spelling())
			this.cp.APf("body", "    return rv2")
		} else {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		}

	case clang.Type_LValueReference:
		if is_qt_class(rety) {
			barety := get_bare_type(rety)
			this.cp.APf("body", "    rv2 := %sNew%sFromPointer(unsafe.Pointer(uintptr(rv))) // 4441",
				pkgSuff, barety.Spelling())
			this.cp.APf("body", "    return rv2")
		} else if TypeIsCharPtr(rety) {
			this.cp.APf("body", "    return qtrt.GoStringI(rv)")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UChar {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv)) /*222*/")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UShort {
			this.cp.APf("body", "    return uint16(rv)")
		} else if isPrimitiveType(rety.PointeeType()) {
			this.cp.APf("body", "    return %s(rv) // 3331", this.tyconver.toDest(rety.PointeeType(), cursor))
		} else {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		}
	case clang.Type_Pointer:
		if is_qt_class(rety) {
			if usemod == "core" && refmod == "widgets" {
				this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
			} else if usemod == "gui" && refmod == "widgets" {
				this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
			} else {
				barety := get_bare_type(rety)
				this.cp.APf("body", "    rv2 := %sNew%sFromPointer(unsafe.Pointer(uintptr(rv))) // 444",
					pkgSuff, barety.Spelling())
				this.cp.APf("body", "    return rv2")
			}
		} else if TypeIsCharPtr(rety) {
			this.cp.APf("body", "    return qtrt.GoStringI(rv)")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UChar {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		} else if rety.PointeeType().CanonicalType().Kind() == clang.Type_UShort {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		} else if isPrimitiveType(rety.PointeeType()) {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
			// this.cp.APf("body", "    return %s(rv) // 333", this.tyconver.toDest(rety.PointeeType(), cursor))
		} else {
			this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv))")
		}
	case clang.Type_RValueReference:
		this.cp.APf("body", "    return unsafe.Pointer(uintptr(rv)) //777")
	case clang.Type_Bool:
		this.cp.APf("body", "    return rv!=0")
	case clang.Type_Enum:
		this.cp.APf("body", "    return int(rv)")
	case clang.Type_Elaborated:
		this.cp.APf("body", "    return int(rv)")
	default:
		this.cp.APf("body", "    return rv")
	}
}

func (this *GenerateGo) genEnums() {

}
func (this *GenerateGo) genEnum() {

}

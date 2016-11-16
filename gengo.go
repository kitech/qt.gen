package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

type GenerateGo struct {
	// TODO move to base
	filter   GenFilter
	mangler  GenMangler
	tyconver TypeConvertor

	methods   []clang.Cursor
	cp        *CodePager
	argDesc   []string
	paramDesc []string
}

func NewGenerateGo() *GenerateGo {
	this := &GenerateGo{}
	this.filter = &GenFilterGo{}
	this.mangler = NewGoMangler()

	this.cp = NewCodePager()
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
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
	this.genExterns(cursor, parent)
	this.genMethods(cursor, parent)
	this.final(cursor, parent)
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
	savefile := fmt.Sprintf("src/%s/%s.go", modname,
		strings.Split(filepath.Base(file.Name()), ".")[0])

	// TODO gofmt the code
	ioutil.WriteFile(savefile, []byte(this.cp.ExportAll()), 0644)
}

func (this *GenerateGo) genHeader(cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	this.cp.APf("header", "// %s", file.Name())
	this.cp.APf("header", "// #include <%s>", filepath.Base(file.Name()))
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	this.cp.APf("header", "// #include <%s>", fullModname)
	this.cp.APf("header", "package %s", strings.ToLower(fullModname[0:]))
	this.cp.APf("header", "")
	this.cp.APf("ext", "")
	this.cp.APf("ext", "/* // extern C begin: %d", len(this.methods))
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

	file, _, _, _ := cursor.Location().FileLocation()
	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]

	this.cp.APf("ext", "// extern C end: %d", len(this.methods))
	this.cp.APf("ext", "*/")
	this.cp.APf("ext", "import \"C\"")
	this.cp.APf("ext", "import \"unsafe\"")
	this.cp.APf("ext", "import \"reflect\"")
	this.cp.APf("ext", "import \"fmt\"")
	this.cp.APf("ext", "import \"qtrt\"")
	for _, dep := range modDeps[modname] {
		this.cp.APf("ext", "import \"qt%s\"", dep)
	}

	this.cp.APf("ext", "")
	this.cp.APf("ext", "func init() {")
	this.cp.APf("ext", "  if false {reflect.TypeOf(123)}")
	this.cp.APf("ext", "  if false {reflect.TypeOf(unsafe.Sizeof(0))}")
	this.cp.APf("ext", "  if false {fmt.Println(123)}")
	this.cp.APf("ext", "  if false {qtrt.KeepMe()}")
	for _, dep := range modDeps[modname] {
		this.cp.APf("ext", "if false {qt%s.KeepMe()}", dep)
	}
	this.cp.APf("ext", "}")
}

func (this *GenerateGo) genMethods(cursor, parent clang.Cursor) {
	log.Println("process class:", len(this.methods), cursor.Spelling())
	grpMethods := this.groupMethods()

	for _, cursors := range grpMethods {
		this.genMethodHeader(cursors[0], cursors[0].SemanticParent())
		this.genMethodInit(cursors[0], cursors[0].SemanticParent())

		for idx, cursor := range cursors {
			this.genVTableTypes(cursor, cursor.SemanticParent(), idx)
		}

		this.genNameLookup(cursors[0], cursors[0].SemanticParent())

		// case x
		for idx, cursor := range cursors {
			parent := cursor.SemanticParent()
			// log.Println(cursor.Kind().String(), cursor.DisplayName())
			switch cursor.Kind() {
			case clang.Cursor_Constructor:
				this.genCtor(cursor, parent, idx)
			case clang.Cursor_Destructor:
				this.genDtor(cursor, parent, idx)
			default:
				if cursor.CXXMethod_IsStatic() {
					this.genStaticMethod(cursor, parent, idx)
				} else {
					this.genNonStaticMethod(cursor, parent, idx)
				}
			}
		}

		this.genMethodFooter(cursors[0], cursors[0].SemanticParent())
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

func (this *GenerateGo) genMethodHeader(cursor, parent clang.Cursor) {
	file, lineno, _, _ := cursor.Location().FileLocation()
	this.cp.APf("main", "// %s:%d", file.Name(), lineno)

	qualities := make([]string, 0)
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
		this.cp.APf("main", "// %s", strings.Join(qualities, " "))
	}

	this.cp.APf("main", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())

}

func (this *GenerateGo) genMethodInit(cursor, parent clang.Cursor) {
	if cursor.Kind() == clang.Cursor_Constructor {
		this.cp.APf("main", "type %s struct {", cursor.Spelling())
		this.cp.APf("main", "    cthis unsafe.Pointer")
		this.cp.APf("main", "}")
	}
	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		this.cp.APf("main", "func (this *%s) %s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()))
	case clang.Cursor_Destructor:
		this.cp.APf("main", "func (this *%s) Delete%s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()[1:]))
	default:
		this.cp.APf("main", "func (this *%s) %s(args...interface{}) {",
			parent.Spelling(), strings.Title(cursor.Spelling()))
	}
	this.cp.AP("main", "  var vtys = make(map[uint8]map[uint8]reflect.Type)")
	this.cp.AP("main", "  if false {fmt.Println(vtys)}")
	this.cp.AP("main", "  var dargExists = make(map[uint8]map[uint8]bool)")
	this.cp.AP("main", "  if false {fmt.Println(dargExists)}")
	this.cp.AP("main", "  var dargValues = make(map[uint8]map[uint8]interface{})")
	this.cp.AP("main", "  if false {fmt.Println(dargValues)}")

	// TODO fill types, default args
}

func (this *GenerateGo) genMethodFooter(cursor, parent clang.Cursor) {
	this.cp.APf("main", "  default:")
	this.cp.APf("main", "    qtrt.ErrorResolve(\"%s\", \"%s\", args)",
		parent.Spelling(), cursor.Spelling())
	this.cp.APf("main", "  } // end switch")
	this.cp.APf("main", "}")
}

func (this *GenerateGo) genVTableTypes(cursor, parent clang.Cursor, midx int) {
	this.cp.APf("main", "  // vtypes %d", midx)
	this.cp.APf("main", "  // dargExists %d", midx)
	this.cp.APf("main", "  // dargValues %d", midx)
	this.cp.APf("main", "  vtys[%d] = make(map[uint8]reflect.Type)", midx)
	this.cp.APf("main", "  dargExists[%d] = make(map[uint8]bool)", midx)
	this.cp.APf("main", "  dargValues[%d] = make(map[uint8]interface{})", midx)
}

func (this *GenerateGo) genNameLookup(cursor, parent clang.Cursor) {
	this.cp.AP("main", "")
	this.cp.AP("main", "  var matchedIndex = qtrt.SymbolResolve(args, vtys)")
	this.cp.AP("main", "  if false {fmt.Println(matchedIndex)}")
	this.cp.AP("main", "  switch matchedIndex {")
}

func (this *GenerateGo) genCtor(cursor, parent clang.Cursor, midx int) {
	// log.Println(this.mangler.convTo(cursor))

	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.cp.APf("main", "    case %d: // (%s), (%s)", midx, argStr, paramStr)
	this.genArgsConv(cursor, parent, midx)
	this.cp.APf("main", "      C.%s(%s)", this.mangler.convTo(cursor), paramStr)
}

func (this *GenerateGo) genDtor(cursor, parent clang.Cursor, midx int) {
	this.cp.APf("main", "    case %d:", midx)
	this.cp.APf("main", "      var cthis unsafe.Pointer = this.cthis")
	this.cp.APf("main", "      C.%s(cthis)", this.mangler.convTo(cursor))
	this.cp.APf("main", "      this.cthis = nil")
}

func (this *GenerateGo) genNonStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	if len(argStr) > 0 {
		argStr = ", " + argStr
	}

	this.cp.APf("main", "    case %d: // (%s), (%s)", midx, argStr, paramStr)
	this.genArgsConv(cursor, parent, midx)
	this.cp.APf("main", "      C.%s(%s)", this.mangler.convTo(cursor), paramStr)
}

func (this *GenerateGo) genStaticMethod(cursor, parent clang.Cursor, midx int) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.cp.APf("main", "    case %d: // (%s), (%s)", midx, argStr, paramStr)
	this.genArgsConv(cursor, parent, midx)
	this.cp.APf("main", "      C.%s(%s)", this.mangler.convTo(cursor), paramStr)
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
		this.argDesc = append(this.argDesc, fmt.Sprintf("%s a%d", cursor.Type().Spelling(), idx))
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
				idx := strings.Index(cursor.Type().Spelling(), " [")
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s %s",
					cursor.Type().Spelling()[0:idx], cursor.Spelling(), cursor.Type().Spelling()[idx+1:]))
			} else {
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
					cursor.Type().Spelling(), cursor.Spelling()))
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
	this.cp.APf("main", "      // var arg%d %s", aidx, "wtype")
	this.cp.APf("main", "      // if %d >= len(args) {", aidx)
	this.cp.APf("main", "      //     arg%d = defaultargx", aidx)
	this.cp.APf("main", "      // } else {")
	this.cp.APf("main", "      //     arg%d = argx.toBind", aidx)
	this.cp.APf("main", "      // }")
}

func (this *GenerateGo) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateGo) genParam(cursor, parent clang.Cursor, idx int) {
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("arg%d", idx))
}

func (this *GenerateGo) genEnums() {

}
func (this *GenerateGo) genEnum() {

}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

type GenerateInline struct {
	// TODO move to base
	filter  GenFilter
	mangler GenMangler

	methods   []clang.Cursor
	cp        *CodePager
	argDesc   []string
	paramDesc []string
}

func NewGenerateInline() *GenerateInline {
	this := &GenerateInline{}
	this.filter = &GenFilterInc{}
	this.mangler = NewIncMangler()

	this.cp = NewCodePager()
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
	}

	return this
}

func (this *GenerateInline) genClass(cursor, parent clang.Cursor) {
	if false {
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	}
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	this.genHeader(cursor, parent)
	this.walkClass(cursor, parent)
	this.genMethods(cursor, parent)
	this.final(cursor, parent)
}

func (this *GenerateInline) final(cursor, parent clang.Cursor) {
	// log.Println(this.cp.ExportAll())
	this.saveCode(cursor, parent)

	this.cp = NewCodePager()
}
func (this *GenerateInline) saveCode(cursor, parent clang.Cursor) {
	// qtx{yyy}, only yyy
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	modname := strings.ToLower(filepath.Base(filepath.Dir(file.Name())))[2:]
	savefile := fmt.Sprintf("src/%s/%s.cxx", modname,
		strings.Split(filepath.Base(file.Name()), ".")[0])

	ioutil.WriteFile(savefile, []byte(this.cp.ExportAll()), 0644)
}

func (this *GenerateInline) genHeader(cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	this.cp.APf("header", "// %s", file.Name())
	this.cp.APf("header", "#include <%s>", filepath.Base(file.Name()))
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	this.cp.APf("header", "#include <%s>", fullModname)
	this.cp.APf("header", "")
}

func (this *GenerateInline) walkClass(cursor, parent clang.Cursor) {

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

func (this *GenerateInline) genMethods(cursor, parent clang.Cursor) {
	log.Println("process class:", len(this.methods), cursor.Spelling())
	for _, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())
		this.genMethodHeader(cursor, parent)

		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			this.genCtor(cursor, parent)
		case clang.Cursor_Destructor:
			this.genDtor(cursor, parent)
		default:
			if cursor.CXXMethod_IsStatic() {
				this.genStaticMethod(cursor, parent)
			} else {
				this.genNonStaticMethod(cursor, parent)
			}
		}
	}
}

func (this *GenerateInline) genMethodHeader(cursor, parent clang.Cursor) {
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

	file, lineno, _, _ := cursor.Location().FileLocation()
	this.cp.APf("main", "// %s:%d", file.Name(), lineno)
	this.cp.APf("main", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())
	this.cp.APf("main", "extern \"C\"")
}

func (this *GenerateInline) genCtor(cursor, parent clang.Cursor) {
	// log.Println(this.mangler.convTo(cursor))
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.cp.APf("main", "void* %s(%s) {", this.mangler.convTo(cursor), argStr)
	this.cp.APf("main", "  return new %s(%s);", parent.Spelling(), paramStr)
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genDtor(cursor, parent clang.Cursor) {
	this.cp.APf("main", "void %s(void *this_) {", this.mangler.convTo(cursor))
	this.cp.APf("main", "  delete (%s*)(this_);", parent.Spelling())
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genNonStaticMethod(cursor, parent clang.Cursor) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")
	if len(argStr) > 0 {
		argStr = ", " + argStr
	}

	this.cp.APf("main", "void %s(void *this_%s) {", this.mangler.convTo(cursor), argStr)

	if cursor.ResultType().Kind() == clang.Type_Void {
		this.cp.APf("main", "  ((%s*)this_)->%s(%s);", parent.Spelling(), cursor.Spelling(), paramStr)
	} else {
		this.cp.APf("main", "  /*return*/ ((%s*)this_)->%s(%s);", parent.Spelling(), cursor.Spelling(), paramStr)
	}
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genStaticMethod(cursor, parent clang.Cursor) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	this.cp.APf("main", "void %s(%s) {", this.mangler.convTo(cursor), argStr)
	if cursor.ResultType().Kind() == clang.Type_Void {
		this.cp.APf("main", "  %s::%s(%s);", parent.Spelling(), cursor.Spelling(), paramStr)
	} else {
		this.cp.APf("main", "  /*return*/ %s::%s(%s);", parent.Spelling(), cursor.Spelling(), paramStr)
	}
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genArgs(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArg(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateInline) genArg(cursor, parent clang.Cursor, idx int) {
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

func (this *GenerateInline) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateInline) genParam(cursor, parent clang.Cursor, idx int) {
	if len(cursor.Spelling()) == 0 {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("arg%d", idx))
	} else {
		this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s", cursor.Spelling()))
	}
}

func (this *GenerateInline) genEnums() {

}
func (this *GenerateInline) genEnum() {

}

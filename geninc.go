package main

import (
	"fmt"
	"gopp"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

type GenerateInline struct {
	// TODO move to base
	filter   GenFilter
	tyconver *TypeConvertGo
	mangler  GenMangler

	methods   []clang.Cursor
	cp        *CodePager
	cpcs      *CodeFS // mod => file =>,
	argDesc   []string
	argtyDesc []string
	paramDesc []string

	GenBase
}

func NewGenerateInline() *GenerateInline {
	this := &GenerateInline{}
	this.filter = &GenFilterInc{}
	this.mangler = NewIncMangler()
	this.tyconver = NewTypeConvertGo()

	this.GenBase.funcMangles = map[string]int{}

	this.cp = NewCodePager()
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
	}

	return this
}

func (this *GenerateInline) initBlocks(cp *CodePager) {
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		cp.AddPointer(block)
		cp.APf(block, "")
	}
}

func (this *GenerateInline) genClass(cursor, parent clang.Cursor) {
	if false {
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	}
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}

	this.isPureVirtualClass = false
	this.genFileHeader(cursor, parent)
	this.walkClass(cursor, parent)
	this.genProtectedCallbacks(cursor, parent)
	this.genProxyClass(cursor, parent)
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
	savefile := fmt.Sprintf("src/%s/%s.cxx", modname, strings.ToLower(cursor.Spelling()))

	ioutil.WriteFile(savefile, []byte(this.cp.ExportAll()), 0644)
}

func (this *GenerateInline) saveCodeToFile(modname, file string) {
	// qtx{yyy}, only yyy
	savefile := fmt.Sprintf("src/%s/%s.cxx", modname, file)
	log.Println(savefile)

	// log.Println(this.cp.AllPoints())
	bcc := this.cp.ExportAll()
	if strings.HasPrefix(bcc, "//") {
		bcc = bcc[strings.Index(bcc, "\n"):]
	}
	ioutil.WriteFile(savefile, []byte(bcc), 0644)

}

func (this *GenerateInline) genFileHeader(cursor, parent clang.Cursor) {
	file, line, col, _ := cursor.Location().FileLocation()
	if false {
		log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
	}
	this.cp.APf("header", "// %s", file.Name())
	hbname := filepath.Base(file.Name())
	if strings.HasSuffix(hbname, "_impl.h") {
		this.cp.APf("header", "#include <%s.h>", hbname[:len(hbname)-7])
	} else {
		this.cp.APf("header", "#include <%s>", filepath.Base(file.Name()))
	}
	fullModname := filepath.Base(filepath.Dir(file.Name()))
	this.cp.APf("header", "#include <%s>", fullModname)
	this.cp.APf("header", "#include \"callback_inherit.h\"")
	this.cp.APf("header", "")
}

func (this *GenerateInline) walkClass(cursor, parent clang.Cursor) {
	pureVirt := false
	methods := make([]clang.Cursor, 0)

	// pcursor := cursor
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			pureVirt = pureVirt || cursor.CXXMethod_IsPureVirtual()
			fallthrough
		case clang.Cursor_Destructor:
			fallthrough
		case clang.Cursor_CXXMethod:
			pureVirt = pureVirt || cursor.CXXMethod_IsPureVirtual()
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

	this.isPureVirtualClass = pureVirt
	if !pureVirt {
		this.isPureVirtualClass = is_pure_virtual_class(cursor)
	}
	this.cp.APf("header", "// %s is pure virtual: %v", cursor.Spelling(), pureVirt)
	this.methods = methods
}

func (this *GenerateInline) genProtectedCallbacks(cursor, parent clang.Cursor) {
	log.Println("process class:", len(this.methods), cursor.Spelling())
	for _, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())

		if cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
			this.genProtectedCallback(cursor, parent)
		}
	}

	this.cp.APf("main", "")
}

func (this *GenerateInline) genProxyClass(cursor, parent clang.Cursor) {
	this.hasVirtualProtected = false
	if is_deleted_class(cursor) {
		return
	}

	this.cp.APf("main", "class My%s : public %s {", cursor.Spelling(), cursor.Type().Spelling())
	this.cp.APf("main", "public:")
	this.cp.APf("main", "  virtual ~My%s() {}", cursor.Spelling())

	for _, mcs := range this.methods {
		if mcs.Kind() == clang.Cursor_Constructor {
			this.cp.APf("main", "// %s %s", mcs.ResultType().Spelling(), mcs.DisplayName())
			this.genArgsPxy(mcs, cursor)
			argStr := strings.Join(this.argDesc, ", ")
			this.genParamsPxy(mcs, cursor)
			paramStr := strings.Join(this.paramDesc, ", ")
			if len(argStr) > 0 {
				// argStr = ", " + argStr
			}
			this.cp.APf("main", "My%s(%s) : %s(%s) {}", mcs.Spelling(),
				argStr, cursor.Type().Spelling(), paramStr)
			continue
		}
		if mcs.AccessSpecifier() != clang.AccessSpecifier_Protected {
			continue
		}
		if !mcs.CXXMethod_IsVirtual() { // 是否只override virtual方法呢？
			// continue
		}

		this.cp.APf("main", "// %s", strings.Join(this.getFuncQulities(mcs), " "))
		this.cp.APf("main", "// %s %s", mcs.ResultType().Spelling(), mcs.DisplayName())
		if mcs.Kind() == clang.Cursor_Destructor {
			continue
		}
		if mcs.Spelling() == "drawItems" || mcs.Spelling() == "getPaintContext" { // temporary skip this
			continue
		}

		this.hasVirtualProtected = true

		// gen projected methods
		this.genArgsPxy(mcs, cursor)
		argStr := strings.Join(this.argDesc, ", ")
		this.genParamsPxy(mcs, cursor)
		paramStr := strings.Join(this.paramDesc, ", ")
		argStr2 := argStr
		if len(argStr2) > 0 {
			argStr2 = ", " + argStr2
		}
		paramStr2 := paramStr
		if len(paramStr2) > 0 {
			paramStr2 = ", " + paramStr2
		}
		argtyStr := strings.Join(this.argtyDesc, ", ")
		if len(argtyStr) > 0 {
			argtyStr = ", " + argtyStr
		}

		this.genArgs(mcs, cursor)
		argtyStr3 := strings.Join(this.argtyDesc, ", ")
		argtyStr3 = gopp.IfElseStr(len(argtyStr3) > 0, ", "+argtyStr3, argtyStr3)
		this.genParamsPxyCall(mcs, cursor)
		paramStr3 := strings.Join(this.paramDesc, ", ")
		paramStr3 = gopp.IfElseStr(len(paramStr3) > 0, ", "+paramStr3, paramStr3)

		prms := []string{}
		prms = append(prms, fmt.Sprintf("%d", mcs.NumArguments()))
		for i := int32(0); i < mcs.NumArguments(); i++ {
			arg := mcs.Argument(uint32(i))
			argty := arg.Type()
			// prmname := gopp.IfElseStr(arg.Spelling() == "", fmt.Sprintf("arg%d", i), arg.Spelling())
			prmname := this.genParamRefName(arg, mcs, int(i))
			switch argty.Kind() {
			case clang.Type_Record:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_LValueReference:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_Double, clang.Type_Float:
				prms = append(prms, "(uint64_t)&"+prmname)
			case clang.Type_Typedef:
				switch argty.Spelling() {
				case "qreal":
					prms = append(prms, "(uint64_t)&"+prmname)
				default:
					prms = append(prms, "(uint64_t)"+prmname)
				}
			default:
				prms = append(prms, "(uint64_t)"+prmname)
			}
		}
		for i := mcs.NumArguments(); i < 10; i++ {
			prms = append(prms, "0")
		}
		prmStr4 := strings.Join(prms, ", ")

		rety := mcs.ResultType()
		this.cp.APf("main", "  virtual %s %s(%s) {", mcs.ResultType().Spelling(), mcs.Spelling(), argStr)
		this.cp.APf("main", "    int handled = 0;")
		this.cp.APf("main", "    auto irv = callbackAllInherits_fnptr(this, (char*)\"%s\", &handled, %s);",
			mcs.Spelling(), prmStr4)
		this.cp.APf("main", "    if (handled) {")
		switch rety.Kind() {
		case clang.Type_Void:
		case clang.Type_Record:
			this.cp.APf("main", "    return *(%s*)(irv);", rety.Spelling())
		case clang.Type_Elaborated, clang.Type_Enum:
			this.cp.APf("main", "    return (%s)(int)(irv);", rety.Spelling())
		case clang.Type_Typedef:
			log.Println(rety.Spelling(), rety.CanonicalType().Kind(), rety.CanonicalType().Spelling(), rety.ClassType().Kind(), rety.ClassType().Spelling())
			if TypeIsQFlags(rety) {
				this.cp.APf("main", "    return (%s)(int)(irv);", rety.Spelling())
			} else if rety.CanonicalType().Kind() == clang.Type_Record {
				this.cp.APf("main", "    return *(%s*)(irv);", rety.Spelling())
			} else {
				this.cp.APf("main", "    return (%s)(irv);", rety.Spelling())
			}
		default:
			this.cp.APf("main", "    return (%s)(irv);", rety.Spelling())
		}
		this.cp.APf("main", "      // %s", rety.Kind().String()+rety.CanonicalType().Kind().String()+rety.CanonicalType().Spelling())
		this.cp.APf("main", "    } else {")
		// TODO check return and convert return if needed
		if mcs.ResultType().Kind() == clang.Type_Void {
			this.cp.APf("main", "    %s::%s(%s);", cursor.Spelling(), mcs.Spelling(), paramStr)
		} else {
			this.cp.APf("main", "    return %s::%s(%s);", cursor.Spelling(), mcs.Spelling(), paramStr)
		}
		this.cp.APf("main", "  }")
		this.cp.APf("main", "  }")
		this.cp.APf("main", "")
	}

	this.cp.APf("main", "};")
	this.cp.APf("main", "")

	// a hotfix
	if this.hasVirtualProtected && cursor.Spelling() == "QVariant" {
		this.hasVirtualProtected = false
	}
}

func (this *GenerateInline) genMethods(cursor, parent clang.Cursor) {
	this.cp.APf("header", "// %s has virtual projected: %v", cursor.Spelling(), this.hasVirtualProtected)
	log.Println("process class:", len(this.methods), cursor.Spelling())

	seeDtor := false
	for _, cursor := range this.methods {
		parent := cursor.SemanticParent()
		// log.Println(cursor.Kind().String(), cursor.DisplayName())
		if cursor.AccessSpecifier() == clang.AccessSpecifier_Protected {
			continue
		}

		this.genMethodHeader(cursor, parent)
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			this.genCtor(cursor, parent)
		case clang.Cursor_Destructor:
			seeDtor = true
			this.genDtor(cursor, parent)
		default:
			if cursor.CXXMethod_IsStatic() {
				this.genStaticMethod(cursor, parent)
			} else {
				this.genNonStaticMethod(cursor, parent)
			}
		}
	}
	if !seeDtor && !is_deleted_class(cursor) && !is_projected_dtor_class(cursor) {
		this.genDtorNotsee(cursor, parent)
	}
}

// TODO move to base
func (this *GenerateInline) genMethodHeader(cursor, parent clang.Cursor) {
	qualities := this.getFuncQulities(cursor)
	if len(qualities) > 0 {
		this.cp.APf("main", "// %s", strings.Join(qualities, " "))
	}

	file, lineno, _, _ := cursor.Location().FileLocation()
	this.cp.APf("main", "// %s:%d", file.Name(), lineno)
	this.cp.APf("main", "// [%d] %s %s",
		cursor.ResultType().SizeOf(), cursor.ResultType().Spelling(), cursor.DisplayName())
	this.cp.APf("main", "extern \"C\"")
}

func (this *GenerateInline) genCtor(cursor, parent clang.Cursor) {
	// log.Println(this.mangler.convTo(cursor))
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	pparent := parent.SemanticParent()
	log.Println(cursor.Spelling(), parent.DisplayName(),
		cursor.SemanticParent().DisplayName(), cursor.LexicalParent().DisplayName(),
		pparent.Spelling(), parent.CanonicalCursor().DisplayName())

	pureVirtRetstr := gopp.IfElseStr(this.isPureVirtualClass, "0; //", "")

	this.cp.APf("main", "void* %s(%s) {", this.mangler.convTo(cursor), argStr)
	pxyclsp := ""
	if !is_deleted_class(parent) && this.hasVirtualProtected {
		this.cp.APf("main", "  auto _nilp = (My%s*)(0);", parent.Spelling())
		pxyclsp = "My"
		// pxyclsp = "" // TODO
	}
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		if pxyclsp == "" {
			this.cp.APf("main", "  return %s new %s::%s(%s);", pureVirtRetstr, pparent.Spelling(), parent.Spelling(), paramStr)
		} else {
			this.cp.APf("main", "  return %s new %s%s(%s);", pureVirtRetstr, pxyclsp, parent.Spelling(), paramStr)
		}
	} else {
		this.cp.APf("main", "  return %s new %s%s(%s);", pureVirtRetstr, pxyclsp, parent.Spelling(), paramStr)
	}

	this.cp.APf("main", "}")
}

func (this *GenerateInline) genDtor(cursor, parent clang.Cursor) {
	pparent := parent.SemanticParent()

	this.cp.APf("main", "void %s(void *this_) {", this.mangler.convTo(cursor))
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		this.cp.APf("main", "  delete (%s::%s*)(this_);", pparent.Spelling(), parent.Spelling())
	} else {
		this.cp.APf("main", "  delete (%s*)(this_);", parent.Type().Spelling())
	}
	this.cp.APf("main", "}")
}

// 在该类没有显式的声明析构方法时，补充生成一个默认析构方法
func (this *GenerateInline) genDtorNotsee(cursor, parent clang.Cursor) {
	// pparent := parent.SemanticParent()

	this.cp.APf("main", "")
	this.cp.APf("main", "extern \"C\"")
	this.cp.APf("main", "void C_ZN%d%sD2Ev(void *this_) {", len(cursor.Spelling()), cursor.Spelling())
	if strings.HasPrefix(parent.Spelling(), "Qt") {
		this.cp.APf("main", "  delete (%s::%s*)(this_);", parent.Spelling(), cursor.Spelling())
	} else {
		this.cp.APf("main", "  delete (%s*)(this_);", cursor.Type().Spelling())
	}
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

	pparent := parent.SemanticParent()
	pparentstr := ""
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		pparentstr = fmt.Sprintf("%s::", pparent.Spelling())
	}

	retstr := "void"
	retset := false
	rety := cursor.ResultType()
	cancpobj := has_copy_ctor(rety.Declaration()) || is_trivial_class(rety.Declaration())
	if rety.Kind() == clang.Type_Void {
	} else if isPrimitiveType(rety) {
		retstr = rety.Spelling()
		retset = true
	} else if rety.Kind() == clang.Type_Pointer {
		retstr = "void*"
		retset = true
	} else {
		if cancpobj {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_Typedef &&
			rety.CanonicalType().Kind() == clang.Type_Record {
			retstr = fmt.Sprintf("%s*", rety.Spelling())
		}
	}

	this.cp.APf("main", "%s %s(void *this_%s) {", retstr, this.mangler.convTo(cursor), argStr)
	log.Println(rety.Spelling(), rety.Declaration().Spelling(), rety.IsPODType())
	if cursor.ResultType().Kind() == clang.Type_Void {
		this.cp.APf("main", "  ((%s%s*)this_)->%s(%s);", pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)
	} else {
		if retset {
			this.cp.APf("main", "  return (%s)((%s%s*)this_)->%s(%s);", retstr, pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)
		} else {
			autoand := gopp.IfElseStr(rety.Kind() == clang.Type_LValueReference, "auto&", "auto")
			this.cp.APf("main", "  %s rv = ((%s%s*)this_)->%s(%s);",
				autoand, pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)

			if cancpobj {
				unconstystr := strings.Replace(rety.Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
				unconstystr := strings.Replace(rety.PointeeType().Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
				this.cp.APf("main", "return &rv;")
			} else if rety.Kind() == clang.Type_Typedef &&
				rety.CanonicalType().Kind() == clang.Type_Record {
				// QModelIndexList
				this.cp.APf("main", "return new %s(rv);", rety.Spelling())
			} else {
				this.cp.APf("main", "/*return rv;*/")
			}
		}
	}
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genStaticMethod(cursor, parent clang.Cursor) {
	this.genArgs(cursor, parent)
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, parent)
	paramStr := strings.Join(this.paramDesc, ", ")

	pparent := parent.SemanticParent()
	pparentstr := ""
	if strings.HasPrefix(pparent.Spelling(), "Qt") {
		pparentstr = fmt.Sprintf("%s::", pparent.Spelling())
	}

	retstr := "void"
	retset := false
	rety := cursor.ResultType()
	cancpobj := has_copy_ctor(rety.Declaration()) || is_trivial_class(rety.Declaration())
	if isPrimitiveType(rety) {
		retstr = rety.Spelling()
		retset = true
	} else if rety.Kind() == clang.Type_Pointer {
		retstr = "void*"
		retset = true
	} else {
		if cancpobj {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
			retstr = "void*"
		}
	}

	this.cp.APf("main", "%s %s(%s) {", retstr, this.mangler.convTo(cursor), argStr)
	if cursor.ResultType().Kind() == clang.Type_Void {
		this.cp.APf("main", "  %s%s::%s(%s);", pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)
	} else {
		if retset {
			this.cp.APf("main", "  return (%s)%s%s::%s(%s);", retstr, pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)
		} else {
			// this.cp.APf("main", "  /*return*/ %s%s::%s(%s);", pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)
			autoand := gopp.IfElseStr(rety.Kind() == clang.Type_LValueReference, "auto&", "auto")
			this.cp.APf("main", "  %s rv = %s%s::%s(%s);",
				autoand, pparentstr, parent.Spelling(), cursor.Spelling(), paramStr)

			if cancpobj {
				unconstystr := strings.Replace(rety.Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
				unconstystr := strings.Replace(rety.PointeeType().Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
				this.cp.APf("main", "return &rv;")
			} else {
				this.cp.APf("main", "/*return rv;*/")
			}
		}
	}
	this.cp.APf("main", "}")
}

func (this *GenerateInline) genProtectedCallback(cursor, parent clang.Cursor) {
	// this.genMethodHeader(cursor, parent)
}

func (this *GenerateInline) genArgsPxy(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	this.argtyDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArgPxy(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateInline) genArgPxy(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))
	csty := cursor.Type()
	argName := this.genParamRefName(cursor, parent, idx)
	if csty.Kind() == clang.Type_LValueReference {
		// 转成指针
	}
	if strings.Contains(csty.CanonicalType().Spelling(), "QFlags<") {
		this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
			csty.CanonicalType().Spelling(), argName))
		this.argtyDesc = append(this.argtyDesc, csty.CanonicalType().Spelling())
	} else {
		log.Println(cursor.Type().Kind(), csty.Spelling(), parent.SemanticParent().Spelling(), parent.DisplayName())
		if TypeIsFuncPointer(csty) {
			this.argDesc = append(this.argDesc,
				strings.Replace(csty.Spelling(), "(*)", fmt.Sprintf("(*%s)", argName), 1))
			this.argtyDesc = append(this.argtyDesc, cursor.Type().Spelling())
		} else if TypeIsCharPtrPtr(cursor.Type()) {
			this.argDesc = append(this.argDesc, fmt.Sprintf("char** %s", argName))
			this.argtyDesc = append(this.argtyDesc, "char**")
		} else if (csty.Kind() == clang.Type_IncompleteArray ||
			csty.Kind() == clang.Type_ConstantArray) &&
			csty.ElementType().Kind() == clang.Type_Pointer {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void** %s", argName))
			this.argtyDesc = append(this.argtyDesc, "void**")
		} else if csty.Kind() == clang.Type_IncompleteArray ||
			csty.Kind() == clang.Type_ConstantArray {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void* %s", argName))
			this.argtyDesc = append(this.argtyDesc, "void*")
			// idx := strings.Index(cursor.Type().Spelling(), " [")
			// this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s %s",
			//	cursor.Type().Spelling()[0:idx], cursor.Spelling(), cursor.Type().Spelling()[idx+1:]))
		} else {
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s", csty.Spelling(), argName))
			this.argtyDesc = append(this.argtyDesc, csty.Spelling())
		}
	}
}

func (this *GenerateInline) genArgs(cursor, parent clang.Cursor) {
	this.argDesc = make([]string, 0)
	this.argtyDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genArg(argc, cursor, idx)
	}
	// log.Println(strings.Join(this.argDesc, ", "), this.mangler.convTo(cursor))
}

func (this *GenerateInline) genArg(cursor, parent clang.Cursor, idx int) {
	// log.Println(cursor.DisplayName(), cursor.Type().Spelling(), cursor.Type().Kind() == clang.Type_LValueReference, this.mangler.convTo(parent))
	csty := cursor.Type()
	argName := this.genParamRefName(cursor, parent, idx)
	if cursor.Type().Kind() == clang.Type_LValueReference {
		// 转成指针
	}
	if strings.Contains(cursor.Type().CanonicalType().Spelling(), "QFlags<") {
		this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
			cursor.Type().CanonicalType().Spelling(), argName))
		this.argtyDesc = append(this.argtyDesc, cursor.Type().CanonicalType().Spelling())
	} else {
		log.Println(csty.Kind(), csty.Spelling(), parent.SemanticParent().Spelling(), parent.DisplayName())
		if csty.Kind() == clang.Type_Record {
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s* %s", cursor.Type().Spelling(), argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s*", cursor.Type().Spelling()))
		} else if csty.Kind() == clang.Type_LValueReference &&
			csty.PointeeType().Kind() == clang.Type_Record {
			if csty.PointeeType().NumTemplateArguments() > 0 {
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s* %s", csty.PointeeType().Spelling(), argName))
				this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s*", csty.PointeeType().Spelling()))
			} else {
				this.argDesc = append(this.argDesc, fmt.Sprintf("%s* %s", csty.PointeeType().Declaration().Spelling(), argName))
				this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s*", csty.PointeeType().Declaration().Spelling()))
			}
		} else if TypeIsFuncPointer(cursor.Type()) {
			this.argDesc = append(this.argDesc,
				strings.Replace(cursor.Type().Spelling(), "(*)",
					fmt.Sprintf("(*%s)", argName), 1))
			this.argtyDesc = append(this.argtyDesc, cursor.Type().Spelling())
		} else if TypeIsCharPtrPtr(cursor.Type()) {
			this.argDesc = append(this.argDesc, fmt.Sprintf("char** %s", argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("char**"))
		} else if (cursor.Type().Kind() == clang.Type_IncompleteArray ||
			cursor.Type().Kind() == clang.Type_ConstantArray) &&
			cursor.Type().ElementType().Kind() == clang.Type_Pointer {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void** %s", argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("void**"))
		} else if cursor.Type().Kind() == clang.Type_IncompleteArray ||
			cursor.Type().Kind() == clang.Type_ConstantArray {
			this.argDesc = append(this.argDesc, fmt.Sprintf("void* %s", argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("void*"))
			// idx := strings.Index(cursor.Type().Spelling(), " [")
			// this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s %s",
			//	cursor.Type().Spelling()[0:idx], cursor.Spelling(), cursor.Type().Spelling()[idx+1:]))
		} else {
			this.argDesc = append(this.argDesc, fmt.Sprintf("%s %s",
				cursor.Type().Spelling(), argName))
			this.argtyDesc = append(this.argtyDesc, fmt.Sprintf("%s",
				cursor.Type().Spelling()))
		}
	}
}

func (this *GenerateInline) genParamsPxy(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamPxy(argc, cursor, idx)
	}
}

func (this *GenerateInline) genParamPxy(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())
	if csty.Kind() == clang.Type_Record { //} &&
		// (parent.Kind() != clang.Cursor_Constructor && !this.hasVirtualProjected) {
		// forceConvStr = "*"
	} else if TypeIsCharPtrPtr(csty) {
		// forceConvStr = "(char**)"
	}
	argName := this.genParamRefName(cursor, parent, idx)
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
}

func (this *GenerateInline) genParamsPxyCall(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParamPxyCall(argc, cursor, idx)
	}
}

func (this *GenerateInline) genParamPxyCall(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())
	if csty.Kind() == clang.Type_Record { //} &&
		// (parent.Kind() != clang.Cursor_Constructor && !this.hasVirtualProjected) {
		forceConvStr = fmt.Sprintf("(%s*)&", csty.Declaration().Spelling())
	} else if csty.Kind() == clang.Type_LValueReference &&
		csty.PointeeType().Kind() == clang.Type_Record {
		forceConvStr = fmt.Sprintf("(%s*)&", csty.PointeeType().Declaration().Spelling())
	} else if TypeIsCharPtrPtr(csty) {
		// forceConvStr = "(char**)"
	}

	argName := this.genParamRefName(cursor, parent, idx)
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
}

func (this *GenerateInline) genParams(cursor, parent clang.Cursor) {
	this.paramDesc = make([]string, 0)
	for idx := 0; idx < int(cursor.NumArguments()); idx++ {
		argc := cursor.Argument(uint32(idx))
		this.genParam(argc, cursor, idx)
	}
}

func (this *GenerateInline) genParam(cursor, parent clang.Cursor, idx int) {
	csty := cursor.Type()
	forceConvStr := ""
	log.Println(csty.Kind().String(), csty.Spelling(), parent.Spelling(), csty.PointeeType().Kind().String(), csty.ArrayElementType().Kind().String())
	if csty.Kind() == clang.Type_Record { //} &&
		// (parent.Kind() != clang.Cursor_Constructor && !this.hasVirtualProjected) {
		forceConvStr = "*"
	} else if csty.Kind() == clang.Type_LValueReference &&
		csty.PointeeType().Kind() == clang.Type_Record {
		forceConvStr = "*"
	} else if TypeIsCharPtrPtr(csty) {
		// forceConvStr = "(char**)"
	}

	argName := this.genParamRefName(cursor, parent, idx)
	this.paramDesc = append(this.paramDesc, fmt.Sprintf("%s%s", forceConvStr, argName))
}

func (this *GenerateInline) genRet(cursor, parent clang.Cursor, idx int) {

}

//
func (this *GenerateInline) genCSignature(cursor, parent clang.Cursor, idx int) {

}

func (this *GenerateInline) genEnumsGlobal(cursor, parent clang.Cursor) {

}
func (this *GenerateInline) genEnum() {

}

func (this *GenerateInline) genFunctions(cursor, parent clang.Cursor) {
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
		// write code
		for _, mod := range modDeps[qtmod] {
			this.cp.APf("header", "#include <Qt%s>", strings.Title(mod))
		}
		this.cp.APf("header", "#include <Qt%s>", strings.Title(qtmod))
		this.cp.APf("header", "#include \"hidden_symbols.h\"")

		sort.Slice(funcs, func(i int, j int) bool {
			return funcs[i].Mangling() > funcs[j].Mangling()
		})
		for _, fc := range funcs {
			log.Println(fc.Spelling(), fc.Mangling(), fc.DisplayName(), fc.IsCursorDefinition())
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
			this.genFunction(fc, olidx)
		}

		this.saveCodeToFile(qtmod, "qfunctions")
	}
}

func (this *GenerateInline) genFunction(cursor clang.Cursor, olidx int) {
	log.Println(cursor.DisplayName(), len(this.funcs))
	this.genArgs(cursor, cursor.SemanticParent())
	argStr := strings.Join(this.argDesc, ", ")
	this.genParams(cursor, cursor.SemanticParent())
	paramStr := strings.Join(this.paramDesc, ", ")

	retstr := "void"
	retset := false
	rety := cursor.ResultType()
	cancpobj := has_copy_ctor(rety.Declaration()) || is_trivial_class(rety.Declaration())
	if isPrimitiveType(rety) {
		retstr = rety.Spelling()
		retset = true
	} else if rety.Kind() == clang.Type_Pointer {
		retstr = "void*"
		retset = true
	} else {
		if cancpobj {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
			retstr = "void*"
		} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
			retstr = "void*"
		}
	}

	overloadSuffix := gopp.IfElseStr(olidx == 0, "", fmt.Sprintf("_%d", olidx))
	this.genMethodHeader(cursor, cursor.SemanticParent())
	this.cp.APf("main", "%s %s%s(%s) {", retstr,
		this.mangler.convTo(cursor), overloadSuffix, argStr)
	if rety.Kind() == clang.Type_Void {
		this.cp.APf("main", "  %s(%s);", cursor.Spelling(), paramStr)
	} else {
		if retset {
			this.cp.APf("main", "  return (%s)%s(%s);", retstr, cursor.Spelling(), paramStr)
		} else {
			autoand := gopp.IfElseStr(rety.Kind() == clang.Type_LValueReference, "auto&", "auto")
			this.cp.APf("main", "  %s rv = %s(%s);", autoand, cursor.Spelling(), paramStr)

			if cancpobj {
				unconstystr := strings.Replace(rety.Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && TypeIsConsted(rety) {
				unconstystr := strings.Replace(rety.PointeeType().Spelling(), "const ", "", 1)
				this.cp.APf("main", "return new %s(rv);", unconstystr)
			} else if rety.Kind() == clang.Type_LValueReference && !TypeIsConsted(rety) {
				this.cp.APf("main", "return &rv;")
			} else {
				this.cp.APf("main", "/*return rv;*/")
			}
		}
	}
	this.cp.APf("main", "}")
	this.cp.APf("main", "")
}

func (this *GenerateInline) genConstantsGlobal(cursor, parent clang.Cursor) {
}

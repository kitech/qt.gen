package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

func (this *GenerateV) genPlainTmplInstClses() {
}

func (this *GenerateV) genTydefTmplInstClses() {
	log.Println("ddddddddd")
	reg := regexp.MustCompile(`^(Q[A-Z].*)([LSHM][ListSetHashMap]+)$`)
	for _, clsinst := range this.tydeftmplinstclses {
		mats := reg.FindAllStringSubmatch(clsinst.Spelling(), -1)
		// log.Println(clsinst.Spelling(), mats)
		if len(mats) == 0 {
			continue
		}
		var undty = clsinst.TypedefDeclUnderlyingType()
		var undcs = undty.TemplateArgumentAsType(0).Declaration()
		if undty.TemplateArgumentAsType(0).Kind() == clang.Type_Pointer {
			undcs = undty.TemplateArgumentAsType(0).PointeeType().Declaration()
		}
		tmplArgClsName := mats[0][1]
		tmplClsName := "Q" + mats[0][2]
		for _, tmplcls := range this.tmplclses {
			if tmplcls.Spelling() == tmplClsName {
				log.Println(tmplClsName, tmplArgClsName)

				this.cp = NewCodePager()
				this.genFileHeader(tmplcls, tmplcls.SemanticParent())
				this.genImports(tmplcls, tmplcls.SemanticParent())
				this.genTemplateInterface(tmplcls, clsinst)
				mod := get_decl_mod(tmplcls)
				if false {
					this.saveCodeToFile(mod, strings.ToLower(tmplcls.Spelling()))
				}

				this.cp = NewCodePager()
				this.genFileHeader(undcs, undcs.SemanticParent())
				this.genImports(undcs, undcs.SemanticParent())
				this.genTemplateInstant(tmplcls, clsinst)
				mod = get_decl_mod(undcs)
				log.Println(mod)
				this.saveCodeToFile(mod, strings.ToLower(clsinst.Spelling()))
				// os.Exit(0)
			}
		}
	}
}

func (this *GenerateV) genTemplateInstant(tmplClsCursor, argClsCursor clang.Cursor) {
	// tmplArgClsName := argClsCursor.Spelling()
	// tmplClsName := tmplClsCursor.Spelling()

	this.cp.APf("body", "pub struct %s  {", argClsCursor.Spelling())
	this.cp.APf("body", "    // *qtrt.CObject")
	this.cp.APf("body", "  pub mut: cthis voidptr")
	this.cp.APf("body", "}")

	this.mthidxs = map[string]int{}
	tmplClsCursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.Mangling())
		case clang.Cursor_CXXMethod:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.NumTemplateArguments(), cursor.Mangling())
			this.genTemplateMethod(cursor, parent, argClsCursor)
		}
		return clang.ChildVisit_Continue
	})

}

func (this *GenerateV) genTemplateMethod(cursor, parent clang.Cursor, argClsCursor clang.Cursor) {
	clsName := argClsCursor.Spelling()
	elemClsName := clsName[:strings.LastIndexAny(clsName, "LHSM")]
	baseMthName := clsName + cursor.Spelling()
	midx := 0
	if midx_, ok := this.mthidxs[baseMthName]; ok {
		this.mthidxs[baseMthName] = midx_ + 1
		midx = midx_ + 1
	} else {
		this.mthidxs[baseMthName] = 0
	}

	rety := cursor.ResultType()
	isSelfRef := func(str string) bool {
		return strings.HasPrefix(str, parent.Spelling()+"<T>")
	}
	isElemRef := func(ty clang.Type) bool {
		log.Println(ty.Spelling(), ty.PointeeType().Spelling(), cursor.DisplayName(), parent.Spelling())
		return ty.Spelling() == "T" || ty.Spelling() == "const T" ||
			ty.PointeeType().Spelling() == "T" || ty.PointeeType().Spelling() == "const T"
	}

	retytxt := ""
	switch rety.Kind() {
	case clang.Type_Int:
		retytxt = "int"
	case clang.Type_Bool:
		retytxt = "bool"
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_Unexposed:
		if isSelfRef(rety.Spelling()) {
			retytxt = "" + clsName
		} else if isElemRef(rety) {
			retytxt = "" + elemClsName
		}
	default:
		log.Println(rety.Spelling(), rety.Kind().Spelling(), cursor.DisplayName())
	}

	validMethodName := rewriteOperatorMethodName(cursor.Spelling())
	cp := this.cp

	retstr := getTyDesc(cursor.ResultType(), AsVReturn, cursor)
	cp.APf("body", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())
	cp.APf("body", "type TC_%s_%s_%d = fn(voidptr) %s",
		clsName, validMethodName, midx, retstr)
	cp.APf("body", "pub fn (this %s) %s%d() %s {",
		clsName, validMethodName, midx, retytxt)
	cp.APf("body", "    // %s_%s_%d()", clsName, validMethodName, midx)
	cp.APf("body", "    // rv, err := qtrt.InvokeQtFunc6(\"C_%s_%s_%d\", qtrt.FFI_TYPE_POINTER, this.Cthis)", clsName, validMethodName, midx)
	cp.APf("body", "    // qtrt.ErrPrint(err, rv)")
	cp.APf("body", "    mut fnobj := TC_%s_%s_%d(0)", clsName, validMethodName, midx)
	cp.APf("body", "    fnobj = qtrt.sym_cfunc6(\"C_%s_%s_%d\")",
		clsName, validMethodName, midx)
	if rety.Kind() != clang.Type_Void {
		cp.APf("body", "    rv :=")
	}
	cp.APf("body", "    fnobj(this.cthis)")

	switch rety.Kind() {
	case clang.Type_Int:
		this.cp.APf("body", "    return 0")
	case clang.Type_Bool:
		this.cp.APf("body", "    return 0==0")
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_Unexposed:
		if isSelfRef(rety.Spelling()) {
			this.cp.APf("body", "    return this")
		} else if isElemRef(rety) {
			this.cp.APf("body", "    return %s{}", elemClsName)
		}
	}

	this.cp.APf("body", "}")

}

func (this *GenerateV) genTemplateInterface(tmplClsCursor, argClsCursor clang.Cursor) {
	if _, ok := tmplclsifgened[tmplClsCursor.Spelling()]; ok {
		// return
	}
	tmplclsifgened[tmplClsCursor.Spelling()] = 1

	log.Printf("%sITF\n", tmplClsCursor.Spelling())
	this.cp.APf("body", "type %sITF interface {", tmplClsCursor.Spelling())

	this.mthidxs = map[string]int{}
	tmplClsCursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.Mangling())
		case clang.Cursor_CXXMethod:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.NumTemplateArguments(), cursor.Mangling())
			// this.genTemplateMethod(cursor, parent, argClsCursor)
			this.genTemplateInterfaceSignature(cursor, cursor.SemanticParent(), argClsCursor)
		}
		return clang.ChildVisit_Continue
	})

	this.cp.APf("body", "}")
	this.cp.APf("body", "")
}

func (this *GenerateV) genTemplateInterfaceSignature(cursor, parent clang.Cursor, argClsCursor clang.Cursor) {
	clsName := argClsCursor.Spelling()
	elemClsName := clsName[:strings.LastIndexAny(clsName, "LHSM")]
	baseMthName := parent.Spelling() + cursor.Spelling() + "ITF"
	midx := 0
	if midx_, ok := this.mthidxs[baseMthName]; ok {
		this.mthidxs[baseMthName] = midx_ + 1
		midx = midx_ + 1
	} else {
		this.mthidxs[baseMthName] = 0
	}

	rety := cursor.ResultType()
	isSelfRef := func(str string) bool {
		return strings.HasPrefix(str, parent.Spelling()+"<T>")
	}
	isElemRef := func(ty clang.Type) bool {
		log.Println(ty.Spelling(), ty.PointeeType().Spelling(), cursor.DisplayName(), parent.Spelling())
		return ty.Spelling() == "T" || ty.Spelling() == "const T" ||
			ty.PointeeType().Spelling() == "T" || ty.PointeeType().Spelling() == "const T"
	}

	retytxt := ""
	switch rety.Kind() {
	case clang.Type_Int:
		retytxt = "int"
	case clang.Type_Bool:
		retytxt = "bool"
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_Unexposed:
		if isSelfRef(rety.Spelling()) {
			retytxt = "*" + clsName
		} else if isElemRef(rety) {
			retytxt = "*" + elemClsName
		}
	default:
		log.Println(rety.Spelling(), rety.Kind().Spelling(), cursor.DisplayName())
	}

	validMethodName := rewriteOperatorMethodName(cursor.Spelling())
	this.cp.APf("body", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())
	this.cp.APf("body", " %s%d() %s ", strings.Title(validMethodName), midx, retytxt)

}

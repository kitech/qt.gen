package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

func (this *GenerateGo) genTemplateSpecializedClasses() {
	log.Println("ddddddddd")
	reg := regexp.MustCompile(`^(Q[A-Z].*)([LSHM][ListSetHashMap]+)$`)
	for _, clsinst := range this.tmplclsspecs {
		mats := reg.FindAllStringSubmatch(clsinst.Spelling(), -1)
		// log.Println(clsinst.Spelling(), mats)
		if len(mats) == 0 {
			continue
		}
		var specls clang.Cursor
		// 查找物化类定义的模块
		clsinst.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			log.Println(cursor.Kind().String(), cursor.Spelling(), parent.Kind().String(), parent.Spelling())
			switch cursor.Kind() {
			case clang.Cursor_TypeRef:
				specls = cursor.Type().Declaration()
				mod := get_decl_mod(specls)
				log.Println(mod, specls.Spelling(), clsinst.Spelling())
				return clang.ChildVisit_Break
			}
			return clang.ChildVisit_Continue
		})
		tmplArgClsName := mats[0][1]
		tmplClsName := "Q" + mats[0][2]
		for _, tmplcls := range this.tmplclses {
			if tmplcls.Spelling() == tmplClsName {
				log.Println(tmplClsName, tmplArgClsName)

				this.cp = NewCodePager()
				this.genHeader(tmplcls, tmplcls.SemanticParent())
				this.genImports(tmplcls, tmplcls.SemanticParent())
				this.genTemplateInterface(tmplcls, clsinst)
				mod := get_decl_mod(tmplcls)
				if false {
					this.saveCodeToFile(mod, strings.ToLower(tmplcls.Spelling()))
				}

				this.cp = NewCodePager()
				this.genHeader(specls, specls.SemanticParent())
				this.genImports(specls, specls.SemanticParent())
				this.genTemplateInstant(tmplcls, clsinst)
				mod = get_decl_mod(specls)
				log.Println(mod)
				this.saveCodeToFile(mod, strings.ToLower(clsinst.Spelling()))
				// os.Exit(0)
			}
		}
	}
}

var mthidxs = map[string]int{}

func (this *GenerateGo) genTemplateInstant(tmplClsCursor, argClsCursor clang.Cursor) {
	// tmplArgClsName := argClsCursor.Spelling()
	// tmplClsName := tmplClsCursor.Spelling()

	this.cp.APf("body", "type %s struct {", argClsCursor.Spelling())
	this.cp.APf("body", "    *qtrt.CObject")
	this.cp.APf("body", "}")

	mthidxs = map[string]int{}
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

func (this *GenerateGo) genTemplateMethod(cursor, parent clang.Cursor, argClsCursor clang.Cursor) {
	clsName := argClsCursor.Spelling()
	elemClsName := clsName[:strings.LastIndexAny(clsName, "LHSM")]
	baseMthName := clsName + cursor.Spelling()
	midx := 0
	if midx_, ok := mthidxs[baseMthName]; ok {
		mthidxs[baseMthName] = midx_ + 1
		midx = midx_ + 1
	} else {
		mthidxs[baseMthName] = 0
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
	this.cp.APf("body", "func (this *%s) %s_%d() %s {",
		clsName, strings.Title(validMethodName), midx, retytxt)
	this.cp.APf("body", "    // %s_%s_%d()", clsName, validMethodName, midx)

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
			this.cp.APf("body", "    return &%s{}", elemClsName)
		}
	}

	this.cp.APf("body", "}")

}

var tmplclsifgened = map[string]int{}

func (this *GenerateGo) genTemplateInterface(tmplClsCursor, argClsCursor clang.Cursor) {
	if _, ok := tmplclsifgened[tmplClsCursor.Spelling()]; ok {
		// return
	}
	tmplclsifgened[tmplClsCursor.Spelling()] = 1

	log.Printf("%s_IF\n", tmplClsCursor.Spelling())
	this.cp.APf("body", "type %s_IF interface {", tmplClsCursor.Spelling())

	mthidxs = map[string]int{}
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
}

func (this *GenerateGo) genTemplateInterfaceSignature(cursor, parent clang.Cursor, argClsCursor clang.Cursor) {
	clsName := argClsCursor.Spelling()
	elemClsName := clsName[:strings.LastIndexAny(clsName, "LHSM")]
	baseMthName := parent.Spelling() + cursor.Spelling() + "_IF"
	midx := 0
	if midx_, ok := mthidxs[baseMthName]; ok {
		mthidxs[baseMthName] = midx_ + 1
		midx = midx_ + 1
	} else {
		mthidxs[baseMthName] = 0
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
	this.cp.APf("body", " %s_%d() %s ", strings.Title(validMethodName), midx, retytxt)

}

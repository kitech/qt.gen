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
		tmplArgClsName := mats[0][1]
		tmplClsName := "Q" + mats[0][2]
		for _, tmplcls := range this.tmplclses {
			if tmplcls.Spelling() == tmplClsName {
				log.Println(tmplClsName, tmplArgClsName)

				this.cp = NewCodePager()
				this.genHeader(clsinst, clsinst.SemanticParent())
				this.genImports(clsinst, clsinst.SemanticParent())
				this.genTemplateInstant(tmplcls, clsinst)
				mod := get_decl_mod(clsinst)
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
	baseMthName := clsName + cursor.Spelling()
	midx := 0
	if midx_, ok := mthidxs[baseMthName]; ok {
		mthidxs[baseMthName] = midx_ + 1
		midx = midx_ + 1
	} else {
		mthidxs[baseMthName] = 0
	}

	validMethodName := rewriteOperatorMethodName(cursor.Spelling())
	this.cp.APf("body", "// %s %s", cursor.ResultType().Spelling(), cursor.DisplayName())
	this.cp.APf("body", "func (this *%s) %s_%d() {", clsName, strings.Title(validMethodName), midx)
	this.cp.APf("body", "    // %s_%s_%d()", clsName, validMethodName, midx)
	this.cp.APf("body", "}")

}

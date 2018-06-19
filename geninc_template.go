package main

import (
	"fmt"
	"gopp"
	"log"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"
)

func (this *GenerateInline) genPlainTmplInstClses() {
	log.Println(len(this.plaintmplinstclses))
	this.cpcs = NewCodeFS()

	for i, icls := range this.plaintmplinstclses {
		log.Println(i, icls.Spelling(), icls.DisplayName())
		qtmod := get_decl_mod(icls)
		if qtmod == "stdglobal" {
			continue
		}

		this.genPlainTmplInstCls(icls)
	}

	// save
	dirs := this.cpcs.ListDirs()
	for _, dir := range dirs {
		files := this.cpcs.ListFiles(dir)
		for _, file := range files {
			cp := this.cpcs.GetFile(dir, file)
			log.Println("saving...", dir, file, cp.TotolLine(), cp.TotolLength())
			path := fmt.Sprintf("src/%s/%s.cxx", dir, file)
			err := cp.WriteFile(path)
			gopp.ErrPrint(err, path)
		}
	}
}

func (this *GenerateInline) genPlainTmplInstCls(ptInstCls clang.Cursor) {
	qtmod := get_decl_mod(ptInstCls)
	qtfile := "qplaintmplinstcls"
	cp := this.cpcs.GetFile(qtmod, qtfile)
	if len(cp.AllPoints()) == 0 {
		this.initBlocks(cp)
		cp.APf("header", "#include <%s>", this.getIncNameByMod(qtmod))
		cp.APf("header", "#include <stdint.h>")
		cp.APf("header", "#include <stdbool.h>")
	}

	ptInstCls.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			log.Println(qtmod, parent.DisplayName(), cursor.Spelling(), cursor.DisplayName(), cursor.Mangling(), cursor.Kind().String())

			this.genArgs(cursor, parent)
			argstr := strings.Join(this.argDesc, ", ")
			this.genParams(cursor, parent)
			prmstr := strings.Join(this.paramDesc, ", ")

			cp.APf("body", "// [%d] %s::%s", parent.Type().SizeOf(), parent.DisplayName(), cursor.DisplayName())
			cp.APf("body", "extern \"C\"")
			cp.APf("body", "void* %s(%s) {", this.mangler.convTo(cursor), argstr)
			cp.APf("body", "    return new %s(%s);", parent.Type().Spelling(), prmstr)
			cp.APf("body", "}")
			cp.APf("body", "")

		case clang.Cursor_Destructor:
			log.Println(qtmod, parent.DisplayName(), cursor.Spelling(), cursor.DisplayName(), cursor.Mangling(), cursor.Kind().String())
			if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
				break
			}
			cp.APf("body", "// [%d] %s::%s", parent.Type().SizeOf(), parent.DisplayName(), cursor.DisplayName())
			cp.APf("body", "extern \"C\"")
			cp.APf("body", "void %s(void* this_) {", this.mangler.convTo(cursor))
			cp.APf("body", "    delete((%s*)this_);", parent.Type().Spelling())
			cp.APf("body", "}")
			cp.APf("body", "")

		case clang.Cursor_CXXMethod:
			log.Println(qtmod, parent.Spelling(), parent.DisplayName(), parent.Mangling(), cursor.Spelling(), cursor.DisplayName(), cursor.Mangling(), cursor.Kind().String())
			if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
				break
			}
			if strings.HasPrefix(parent.Type().Spelling(), "QtPrivate::") {
				break
			}
			if strings.Contains(parent.Type().Spelling(), "QInputMethodQueryEvent::QueryPair") {
				break // it's private member
			}

			rety := cursor.ResultType()
			retstr := "void"
			retstmt := "/*return rv;*/"
			retassign := gopp.IfElseStr(rety.Kind() == clang.Type_Void, "/*auto rv = */", "auto rv = ")
			switch {
			case rety.Kind() == clang.Type_Bool || rety.Kind() == clang.Type_Int:
				retstr = getTyDesc(rety, AsCReturn, cursor)
				retstmt = "return rv;"
			case rety.Kind() == clang.Type_Record:
				retstr = getTyDesc(rety, AsCReturn, cursor)
				retstmt = fmt.Sprintf("return new %s(rv);", rety.Spelling())
			case rety.Kind() == clang.Type_LValueReference:
				if rety.PointeeType().Kind() == clang.Type_Record && is_qt_class(rety.PointeeType()) {
					retstr = getTyDesc(rety, AsCReturn, cursor)
					retstmt = fmt.Sprintf("return new %s(rv);", rety.PointeeType().Spelling())
				}
			case rety.Kind() == clang.Type_Pointer:
				retstr = getTyDesc(rety, AsCReturn, cursor)
				retstmt = fmt.Sprintf("return (%s)rv;", retstr)
			case rety.Kind() == clang.Type_Typedef:
				if rety.Declaration().TypedefDeclUnderlyingType().Kind() == clang.Type_Record &&
					is_qt_class(rety) {
					retstr = getTyDesc(rety, AsCReturn, cursor)
					retstmt = fmt.Sprintf("return new %s(rv);", rety.Spelling())
				}
			case rety.Kind() == clang.Type_Unexposed && rety.CanonicalType().Kind() == clang.Type_Record && is_qt_class(rety.CanonicalType()):
				retstr = getTyDesc(rety, AsCReturn, cursor)
				retstmt = fmt.Sprintf("return new %s(rv);", rety.CanonicalType().Spelling())
			default:
				if TypeIsCharPtr(rety) {
					retstr = getTyDesc(rety, AsCReturn, cursor)
					retstmt = fmt.Sprintf("return (%s)rv;", retstr)
				}
			}

			this.genArgs(cursor, parent)
			argstr := gopp.IfElseStr(len(this.argDesc) > 0, ", "+strings.Join(this.argDesc, ", "), "")
			this.genParams(cursor, parent)
			prmstr := strings.Join(this.paramDesc, ", ")
			cp.APf("body", "// [%d] %s %s::%s", cursor.ResultType().SizeOf(), cursor.ResultType().Spelling(), parent.DisplayName(), cursor.DisplayName())
			cp.APf("body", "extern \"C\"")
			cp.APf("body", "%s %s(void* this_ %s) {", retstr, this.mangler.convTo(cursor), argstr)
			cp.APf("body", "   %s ((%s*)this_)->%s(%s);",
				retassign, parent.Type().Spelling(), cursor.Spelling(), prmstr)
			cp.APf("body", "   %s", retstmt)

			cp.APf("body", "}")
			cp.APf("body", "")
		}
		return clang.ChildVisit_Continue
	})
}

// TODO 无法找到这种模板实例化类的定义，无法遍历其方法并得到方法的mangling name
// 也就是应该是在头文件里出现的这种typedef模板类声明实际上还没有实例化。
func (this *GenerateInline) genTydefTmplInstClses() {
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
		tmplElmClsName := mats[0][1]
		tmplClsName := "Q" + mats[0][2]
		for _, tmplcls := range this.tmplclses {
			log.Println(tmplcls.Spelling(), tmplClsName)
			if tmplcls.Spelling() != tmplClsName {
				continue
			}
			log.Println(tmplClsName, tmplElmClsName)

			this.cp = NewCodePager()
			clsctx := &GenClassContext{}
			clsctx.clscs = undcs
			this.genFileHeader(clsctx, undcs, undcs.SemanticParent())
			// this.genImports(tmplcls, tmplcls.SemanticParent())
			this.genTemplateInterface(tmplcls, clsinst)
			mod := get_decl_mod(tmplcls)
			if false {
				this.saveCodeToFile(mod, strings.ToLower(tmplcls.Spelling()))
			}

			this.cp = NewCodePager()
			this.genFileHeader(clsctx, undcs, undcs.SemanticParent())
			this.cp.APf("header", "#ifndef %sList", undcs.Spelling())
			this.cp.APf("header", "#ifndef %sLIST_H", strings.ToUpper(undcs.Spelling()))
			this.cp.APf("header", "typedef QList<%s> %sList;",
				undty.TemplateArgumentAsType(0).Spelling(), undcs.Spelling())
			this.cp.APf("header", "#endif")
			this.cp.APf("header", "#endif")
			// this.genImports(specls, specls.SemanticParent())
			this.genTemplateInstant(tmplcls, clsinst)
			mod = get_decl_mod(undcs)
			log.Println(mod, undcs.Spelling(), undty.Spelling())
			this.saveCodeToFile(mod, strings.ToLower(clsinst.Spelling()))
			// os.Exit(0)
		}
	}
}

func (this *GenerateInline) genTemplateInstant(tmplClsCursor, instClsCursor clang.Cursor) {

	this.mthidxs = map[string]int{}
	tmplClsCursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.Mangling())
			// this.genTemplateMethod(cursor, parent, argClsCursor)
		case clang.Cursor_CXXMethod:
			log.Println(cursor.Spelling(), cursor.DisplayName(), cursor.NumTemplateArguments(), cursor.Mangling(), cursor.AccessSpecifier().String(), parent.Spelling(), instClsCursor.DisplayName())
			// 忽略const iterator返回值的方法
			if cursor.AccessSpecifier() == clang.AccessSpecifier_Public &&
				!isTmplConstAnyIterRef(cursor.ResultType().Spelling(), parent) {
				if !funk.ContainsString([]string{"fromVector", "fromSet", "fromStdList", "fromList", "toSet"}, cursor.Spelling()) {
					// QJSValue缺少operator==导致完整的QList展开编译错误
					if instClsCursor.DisplayName() == "QJSValueList" {
						if !funk.ContainsString([]string{"startsWith", "endsWith", "count", "contains", "lastIndexOf", "indexOf", "removeAll", "removeOne", "operator==", "operator!="}, cursor.Spelling()) {
							this.genTemplateMethod(cursor, parent, instClsCursor)
						}
					} else {
						this.genTemplateMethod(cursor, parent, instClsCursor)
					}
				}
			}
		}
		return clang.ChildVisit_Continue
	})

}

func (this *GenerateInline) genTemplateMethod(cursor, parent clang.Cursor, instClsCursor clang.Cursor) {
	undty := instClsCursor.TypedefDeclUnderlyingType()
	clsName := instClsCursor.Spelling()
	elemClsName := undty.TemplateArgumentAsType(0).Spelling()
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

	retytxt := ""
	autoretxt := gopp.IfElseStr(rety.Kind() == clang.Type_Void, "", "auto rv = ")
	switch rety.Kind() {
	case clang.Type_Int:
		retytxt = "int"
	case clang.Type_Bool:
		retytxt = "bool"
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_Unexposed:
		if isSelfRef(rety.Spelling()) {
			retytxt = clsName + "*"
		} else if isTmplElemRef(rety) {
			if strings.Contains(instClsCursor.DisplayName(), "Hash") ||
				strings.Contains(instClsCursor.DisplayName(), "Map") {
				retytxt = "QVariant*"
			} else {
				retytxt = elemClsName + "*"
			}
		} else {
			retytxt = "void"
		}
	case clang.Type_Void:
		retytxt = "void"
	default:
		if isTmplAnyIterRef(rety.Spelling(), parent) {
			retytxt = fmt.Sprintf("%s::%s*", instClsCursor.DisplayName(), strings.Split(rety.Spelling(), "::")[1])
		} else {
			retytxt = "void"
		}
		log.Println(rety.Spelling(), rety.Kind().Spelling(), cursor.DisplayName())
	}

	argsDescs, prmsDescs := this.genTmplFuncArgs(cursor, parent, clsName, elemClsName)
	argsDesc := strings.Join(argsDescs, ", ")
	argsDesc = gopp.IfElseStr(len(argsDesc) == 0, argsDesc, ", "+argsDesc)
	prmsDesc := strings.Join(prmsDescs, ", ")

	validMethodName := rewriteOperatorMethodName(cursor.Spelling())
	this.cp.APf("body", "// [%d] %s %s", cursor.ResultType().SizeOf(), cursor.ResultType().Spelling(), cursor.DisplayName())
	this.cp.APf("body", "extern \"C\"")
	this.cp.APf("body", "%s C_%s_%s_%d(void* this_ %s) {",
		retytxt, clsName, validMethodName, midx, argsDesc)
	this.cp.APf("body", "    // %s_%s_%d()", clsName, validMethodName, midx)
	this.cp.APf("body", "    %s ((%s*)this_)->%s(%s);", autoretxt, clsName, cursor.Spelling(), prmsDesc)

	switch rety.Kind() {
	case clang.Type_Int:
		this.cp.APf("body", "    return rv;")
	case clang.Type_Bool:
		this.cp.APf("body", "    return rv;")
	case clang.Type_LValueReference:
		fallthrough
	case clang.Type_Unexposed:
		if isSelfRef(rety.Spelling()) {
			this.cp.APf("body", "    return (%s*)this_;", clsName)
		} else if isTmplElemRef(rety) {
			if strings.Contains(rety.Spelling(), "*") {
				this.cp.APf("body", "    return rv;")
			} else {
				this.cp.APf("body", "    return new decltype(rv)(rv);")
			}
		}
	default:
		if isTmplAnyIterRef(rety.Spelling(), parent) {
			this.cp.APf("body", "    return new decltype(rv)(rv);")
		} else {
		}
	}

	this.cp.APf("body", "}")
	this.cp.APf("body", "")
}

func (this *GenerateInline) genTmplFuncArgs(cursor, parent clang.Cursor, clsName, elemClsName string) ([]string, []string) {
	argstrs := []string{}
	prmstrs := []string{}
	for i := int32(0); i < cursor.NumArguments(); i++ {
		argstr, prmstr := this.genTmplFuncArg(cursor.Argument(uint32(i)), cursor, parent, clsName, elemClsName)
		argstrs = append(argstrs, argstr)
		prmstrs = append(prmstrs, prmstr)
	}
	return argstrs, prmstrs
}

func (this *GenerateInline) genTmplFuncArg(cursor, parent, pparent clang.Cursor, clsName, elemClsName string) (string, string) {
	var argstr string
	var prmstr string
	log.Println(cursor.Type().Spelling(), cursor.Type().Declaration().DisplayName(), cursor.Type().Declaration().NumTemplateArguments(), pparent.Spelling(), pparent.DisplayName())
	if isTmplSelfRef(cursor.Type().Spelling(), pparent) {
		argstr = fmt.Sprintf("%s* %s", clsName, cursor.Spelling())
		prmstr = fmt.Sprintf("*%s", cursor.Spelling())
	} else if isTmplElemRef(cursor.Type()) {
		if strings.HasSuffix(elemClsName, "*") {
			argstr = fmt.Sprintf("%s %s", elemClsName, cursor.Spelling())
			prmstr = fmt.Sprintf("%s", cursor.Spelling())
		} else {
			argstr = fmt.Sprintf("%s* %s", elemClsName, cursor.Spelling())
			prmstr = fmt.Sprintf("*%s", cursor.Spelling())
		}
	} else if isTmplAnyIterRef(cursor.Type().Spelling(), pparent) {
		argstr = fmt.Sprintf("%s::%s* %s", clsName, strings.Split(cursor.Type().Spelling(), "::")[1], cursor.Spelling())
		prmstr = fmt.Sprintf("*%s", cursor.Spelling())
	} else if isTmplNodeRef(cursor.Type().Spelling(), pparent) {
		argstr = fmt.Sprintf("%s::Node %s", clsName, cursor.Spelling())
		prmstr = fmt.Sprintf("%s", cursor.Spelling())
	} else if strings.ContainsAny(cursor.Type().Spelling(), "<T>") { // unknown
		argstr = fmt.Sprintf("void* %s", cursor.Spelling())
		prmstr = fmt.Sprintf("/*%s*/0", cursor.Spelling())
	} else if isTmplKeyRef(cursor.Type().Spelling(), pparent) {
		argstr = fmt.Sprintf("QString* %s", cursor.Spelling())
		prmstr = fmt.Sprintf("*%s", cursor.Spelling())
	} else {
		argstr = fmt.Sprintf("%s %s", cursor.Type().Spelling(), cursor.Spelling())
		prmstr = fmt.Sprintf("%s", cursor.Spelling())
	}
	return argstr, prmstr
}

var isTmplSelfRef = func(str string, parent clang.Cursor) bool {
	log.Println(parent.DisplayName(), str)
	reg := `.*(Q.+)<([KeyTKV, ]+)>*.`
	exp := regexp.MustCompile(reg)
	mats1 := exp.FindAllStringSubmatch(str, -1)
	mats2 := exp.FindAllStringSubmatch(parent.DisplayName(), -1)
	log.Println(mats1, mats2)
	return len(mats1) > 0 && len(mats2) > 0 && mats1[0][1] == mats2[0][1] &&
		len(strings.Split(mats1[0][2], ",")) == len(strings.Split(mats2[0][2], ","))
	// return strings.Contains(str, parent.DisplayName()) ||
	//	strings.Contains(parent.DisplayName(), str)
	// return strings.Contains(str, parent.Spelling()+"<T>")
}
var isTmplElemRef = func(ty clang.Type) bool {
	return ty.Spelling() == "T" || ty.Spelling() == "const T" ||
		ty.PointeeType().Spelling() == "T" || ty.PointeeType().Spelling() == "const T"
}
var isTmplIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::iterator")
}
var isTmplConstIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_iterator")
}
var isTmplRIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::reverse_iterator")
}
var isTmplConstRIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_reverse_iterator")
}
var isTmplKIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::key_iterator")
}
var isTmplConstKIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_key_iterator")
}
var isTmplKVIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::key_value_iterator")
}
var isTmplConstKVIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_key_value_iterator")
}
var isTmplAnyIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::") && strings.Contains(str, "iterator")
}
var isTmplConstAnyIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_") && strings.Contains(str, "iterator")
}
var isTmplNodeRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::Node")
}
var isTmplKeyRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, "Key")
}

func (this *GenerateInline) genTemplateInterface(tmplClsCursor, argClsCursor clang.Cursor) {
	if _, ok := tmplclsifgened[tmplClsCursor.Spelling()]; ok {
		// return
	}
	tmplclsifgened[tmplClsCursor.Spelling()] = 1

	log.Printf("%s_IF\n", tmplClsCursor.Spelling())
	this.cp.APf("body", "type %s_IF interface {", tmplClsCursor.Spelling())

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
}

func (this *GenerateInline) genTemplateInterfaceSignature(cursor, parent clang.Cursor, argClsCursor clang.Cursor) {
	clsName := argClsCursor.Spelling()
	elemClsName := clsName[:strings.LastIndexAny(clsName, "LHSM")]
	baseMthName := parent.Spelling() + cursor.Spelling() + "_IF"
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
	this.cp.APf("body", " %s_%d() %s ", strings.Title(validMethodName), midx, retytxt)

}

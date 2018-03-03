package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"

	"gopp"
)

var ast_file = "./qthdrsrc.ast"
var hdr_file = "./bsheaders/qthdrsrc.h"

// module depend table
var modDeps = modDepsAll               // auto generated
var skipClasses = make(map[string]int) // 全局过滤掉的class

func init() {
	if false {
		log.Println(123)
	}
}

type GenCtrl struct {
	tu       clang.TranslationUnit
	tuc      clang.Cursor
	cidx     clang.Index
	modules  []string
	args     []string
	save_ast bool

	filter     GenFilter
	genor      Generator
	qtenumgen  Generator // for generate global enums
	qtfuncgen  Generator // for generate global functions
	qttmplgen  Generator // for generate template
	qtconstgen Generator // for generate global macros
}

func NewGenCtrl() *GenCtrl {
	this := &GenCtrl{}

	return this
}

var genLang string = "" // c(c binding), go (go binding), rs (rust binding)

func (this *GenCtrl) main() {
	if len(os.Args) > 1 {
		genLang = os.Args[1]
	}
	if genLang == "" {
		log.Fatalln("must suply a lang to gen, usage: qt.gen <c|go|rs>")
	}

	this.setupLang()
	this.setupEnv()
	this.createTU()
	this.collectClasses()
	this.cleanupEnv()
}

func (this *GenCtrl) setupLang() {

	switch genLang {
	case "c":
		this.filter = &GenFilterInc{}
		this.genor = NewGenerateInline()
		this.qtenumgen = NewGenerateInline()
		this.qtfuncgen = NewGenerateInline()
		this.qttmplgen = NewGenerateInline()
		this.qtconstgen = NewGenerateInline()
	case "go":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateGo()
		this.qtenumgen = NewGenerateGo()
		this.qtfuncgen = NewGenerateGo()
		this.qttmplgen = NewGenerateGo()
		this.qtconstgen = NewGenerateGo()
	case "rs":
		fallthrough
	default:
		log.Fatalln("not supported or not impled:", genLang)
	}
}

func (this *GenCtrl) setupEnv() {
	cidx := clang.NewIndex(0, 1)
	// defer cidx.Dispose()

	// 这是要生成的模块表
	modules := []string{
		"QtCore", "QtGui", "QtWidgets",
		"QtNetwork", "QtQml", "QtQuick",
		"QtQuickTemplates2", "QtQuickControls2", "QtQuickWidgets",
		// for platform dependent modules, need copy headers if not exists
		"QtAndroidExtras", // TODO fatal error: 'jni.h' file not found
		"QtX11Extras",     // 这个包没生成出来什么代码
		"QtWinExtras",     // 缺少QtWinExtracsDepened头文件
		"QtMacExtras",     // 缺少QtMacExtracsDepened头文件
	}

	cmdlines := []string{
		"-x c++ -std=c++11 -D__CODE_GENERATOR__ -D_GLIBCXX_USE_CXX11ABI=1",
		"-DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++11 -Wall -W -D_REENTRANT -fPIC",
		"-I./bsheaders", "-I/usr/include/wine/windows/", // fix cross platform generate, win/mac
	}
	args := []string{}
	gopp.Domap(cmdlines, func(e interface{}) interface{} {
		args = append(args, strings.Split(e.(string), " ")...)
		return nil
	})

	qtdir := gopp.IfElseStr(os.Getenv("QT_DIR") == "", "/usr", "./qtheaders")
	if qtdir == "/usr" {
		args = append(args, fmt.Sprintf("-I%s/include/qt", qtdir))
	} else {
		args = append(args, fmt.Sprintf("-I./qtheaders/include"))
	}

	gopp.Domap(modules, func(e interface{}) interface{} {
		args = append(args, fmt.Sprintf("-DQT_%s_LIB", strings.ToUpper(e.(string)[2:])))
		args = append(args, fmt.Sprintf("-DGEN_GO_QT_%s_LIB", strings.ToUpper(e.(string)[2:])))
		if qtdir == "/usr" {
			args = append(args, fmt.Sprintf("-I/usr/include/qt/%s", e.(string)))
		} else {
			// args = append(args, fmt.Sprintf("-I%s/%s/gcc_64/include/%s", qtdir, qtver, e.(string)))
			args = append(args, fmt.Sprintf("-I./qtheaders/include/%s", e.(string)))
		}
		return nil
	})
	cmd := exec.Command("g++", "--print-file-name=include")
	out, err := cmd.Output()
	gopp.ErrPrint(err)
	args = append(args, fmt.Sprintf("-I%s", string(out[:len(out)-1])))
	args = append(args, fmt.Sprintf("-I%s-fixed", string(out[:len(out)-1])))
	cmd = exec.Command("g++", "-dumpversion")
	out, err = cmd.Output()
	gopp.ErrPrint(err)
	args = append(args, fmt.Sprintf("-I/usr/include/c++/%s", strings.TrimSpace(string(out))))
	log.Println(args)
	fullCmd := fmt.Sprintf("g++ %s -o qthdrsrc.o -c %s", strings.Join(args, " "), hdr_file)
	ioutil.WriteFile("bcmd.sh", []byte(fullCmd), 0755)
	// os.Exit(0)

	this.cidx = cidx
	this.args = args
	this.modules = modules
}

func (this *GenCtrl) createTU() {
	cidx := this.cidx
	args := this.args

	var tu clang.TranslationUnit
	save_ast := false
	if _, err := os.Stat(ast_file); err == nil {
		tu = cidx.TranslationUnit(ast_file)
	} else {
		save_ast = true
		opts := uint32(0)
		opts |= clang.TranslationUnit_DetailedPreprocessingRecord
		opts |= clang.TranslationUnit_IncludeBriefCommentsInCodeCompletion
		tu = cidx.ParseTranslationUnit(hdr_file, args, nil, opts)
	}
	if !tu.IsValid() {
		log.Panicln("wtf", "maybe cached qthdrsrc.ast file expired, delete and retry please.")
	}
	cursor := tu.TranslationUnitCursor()
	if false {
		log.Println(cursor)
	}

	this.tuc = cursor
	this.tu = tu
	this.save_ast = save_ast
}

func (this *GenCtrl) visfn(cursor, parent clang.Cursor) clang.ChildVisitResult {
	{
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName(), cursor.SpecializedCursorTemplate().DisplayName(), cursor.CanonicalCursor().Kind(), cursor.IsCursorDefinition())
	}

	switch cursor.Kind() {
	case clang.Cursor_ClassDecl:
		if !cursor.IsCursorDefinition() {
			break
		}
		if !is_qt_class(cursor.Type()) {
			break
		}
		clts.ClassCount += 1
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
		if !this.filter.skipClass(cursor, parent) {
			this.genor.genClass(cursor, parent)
		} else {
			clts.SkippedClassCount += 1
		}
		if cursor.Type().SizeOf() > clts.MaxClassSize {
			clts.MaxClassSize = cursor.Type().SizeOf()
			clts.MaxSizeClass = cursor.Type().Spelling()
		}
		clts.addClassSize(cursor.Type().SizeOf())
		// cursor.Visit(this.visfn)
	case clang.Cursor_FunctionDecl:
		clts.FunctionCount += 1
		clts.funcParents[parent.Spelling()] = 1
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName(), parent.Spelling(), cursor.Mangling(), len(clts.funcParents), cursor.IsCursorDefinition(), cursor.Definition().Spelling())
		this.qtfuncgen.putFunc(cursor)
	case clang.Cursor_StructDecl:
		if !this.filter.skipClass(cursor, parent) {
			this.genor.genClass(cursor, parent)
			// return clang.ChildVisit_Break
		}
	case clang.Cursor_CXXMethod:
		clts.MethodCount += 1
		if this.filter.skipMethod(cursor, parent) {
			clts.SkippedMethodCount += 1
		}
		if is_private_method(cursor) {
			clts.PrivateMethodCount += 1
		}
	case clang.Cursor_TypedefDecl:
		undty := cursor.TypedefDeclUnderlyingType()
		undcs := undty.Declaration()
		log.Println(cursor.Kind().String(), cursor.DisplayName(), undcs.Kind().String(), undcs.DisplayName(), undty.Spelling())
		if undcs.Kind() == clang.Cursor_NoDeclFound {

		}
		if funk.ContainsInt([]int{clang.Cursor_ClassDecl, clang.Cursor_StructDecl}, int(undcs.Kind())) &&
			strings.HasPrefix(cursor.Spelling(), "Q") {
			this.qttmplgen.putTydefTmplClsInst(cursor)
			clts.TydefTmplInstClsCount += 1
		}

		if is_qt_class(cursor.Type()) {
			log.Println("got", cursor.Spelling(), ",", cursor.Type().Spelling(), ",", cursor.Type().Kind(), cursor.Type().ClassType().Spelling(), ",", cursor.Type().ClassType().Kind(), ",", cursor.Type().CanonicalType().Spelling(), ",", cursor.Type().CanonicalType().Kind(), ",", cursor.Type().CanonicalType().Declaration().Kind(), ",", cursor.Type().CanonicalType().Declaration().Spelling(), ",", cursor.Type().CanonicalType().Declaration().Definition().Kind().String())
			log.Println(cursor.Type().CanonicalType().Declaration().Kind().String(), cursor.Type().CanonicalType().Declaration().Spelling(), cursor.Type().CanonicalType().Declaration().IsCursorDefinition(), cursor.Spelling())
			// this.qttmplgen.putPlainTmplClsInst(cursor)

			tplc := cursor.Type().Declaration()
			log.Println(tplc.Spelling(), "============", cursor.CanonicalCursor().Kind(), tplc.NumTemplateArguments())
			// cursor.Visit(this.visfn)
		}

	case clang.Cursor_TypeRef:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_TemplateRef:
		log.Println(cursor.Definition().Kind(), cursor.Definition().Spelling())
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_ClassTemplate:
		if !cursor.IsCursorDefinition() {
			break
		}
		this.qttmplgen.putTmplCls(cursor)
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
		log.Println(cursor.IsCursorDefinition(), cursor.NumTemplateArguments())
	case clang.Cursor_ClassTemplatePartialSpecialization:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_FunctionTemplate:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_Constructor:
		clts.MethodCount += 1
	case clang.Cursor_Destructor:
		clts.MethodCount += 1
	case clang.Cursor_ConversionFunction:
	case clang.Cursor_VarDecl:
	case clang.Cursor_EnumDecl:
		clts.EnumCount += 1
		this.qtenumgen.putEnum(cursor)
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_EnumConstantDecl:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_UnionDecl:
	case clang.Cursor_Namespace:
		cursor.Visit(this.visfn)
	case clang.Cursor_UsingDeclaration:
	case clang.Cursor_StaticAssert:
	case clang.Cursor_UnexposedDecl:
	case clang.Cursor_TypeAliasTemplateDecl:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
	case clang.Cursor_MacroDefinition:
		clts.ConstCount += 1
		this.qtconstgen.putConstant(cursor)
	case clang.Cursor_InvalidCode:
		fallthrough
	default:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName(), cursor.Type().Spelling())
	}

	// 查找直接的实例化过的模板类定义
	if funk.ContainsInt([]int{clang.Cursor_StructDecl, clang.Cursor_ClassDecl}, int(cursor.Kind())) {
		spcs := cursor.SpecializedCursorTemplate()
		if spcs.Kind() != clang.Cursor_InvalidFile {
			fi, lineno, _, _ := spcs.Location().FileLocation()
			fi2, lineno2, _, _ := cursor.Location().FileLocation()
			log.Println(spcs.Kind().String(), spcs.Spelling(), fi.Name(), lineno, spcs.DisplayName(), cursor.DisplayName(), cursor.Spelling(), cursor.Kind().String(), fi2.Name(), lineno2, cursor.TemplateArgumentType(0).Spelling(), spcs.TemplateArgumentType(0).Spelling())
			this.qttmplgen.putPlainTmplClsInst(cursor)
			clts.PlainTmplInstClsConut += 1
		}
	}

	clts.CursorCount += 1
	return clang.ChildVisit_Continue
	// return clang.ChildVisit_Recurse
}

func (this *GenCtrl) collectClasses() {
	cursor := this.tuc

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl:
			if !cursor.IsCursorDefinition() {
				break
			}
			if !this.filter.skipClass(cursor, parent) {
				skipClasses[cursor.Type().Spelling()] = 1
			}
			log.Println(cursor.DisplayName(), cursor.BriefCommentText(), cursor.RawCommentText())
		case clang.Cursor_StructDecl:
		case clang.Cursor_MacroDefinition:
			readSourceRange(cursor.Extent())
		}
		return clang.ChildVisit_Continue
	})

	cursor.Visit(this.visfn)

	this.qtfuncgen.genFunctions(cursor, cursor.SemanticParent())
	if genLang == "go" {
		var gg *GenerateGo = this.qtenumgen.(*GenerateGo)
		gg.cp.APf("header", "package qtcore")
		gg.genEnumsGlobal(cursor, cursor.SemanticParent())
		gg.saveCodeToFile("core", "qnamespace")

		gg = this.qtconstgen.(*GenerateGo)
		gg.cp.APf("header", "package qtcore")
		gg.genConstantsGlobal(cursor, cursor.SemanticParent())
		gg.saveCodeToFile("core", "qconstants")

		/*
			gg = this.genor.(*GenerateGo)
			for mod, cp := range gg.cpcs {
				gg.saveCodeToFileWithCode(mod, "qcallbacks", cp.ExportAll())
			}
		*/
	}
	this.qttmplgen.genPlainTmplInstClses()
	this.qttmplgen.genTydefTmplInstClses()

	log.Printf("%+v\n", clts)
}

func (this *GenCtrl) cleanupEnv() {
	if this.save_ast {
		this.tu.SaveTranslationUnit(ast_file, 0)
	}
}

type collects struct {
	CursorCount   int
	MaxClassSize  int64
	MaxSizeClass  string
	ClassSizeMap  map[int64]int // size => count
	ClassCount    int           // 查找到的所有的qt类
	MethodCount   int           // 查找到的所有的qt方法
	FunctionCount int           // 全局函数
	EnumCount     int
	ConstCount    int

	TydefTmplInstClsCount int
	PlainTmplInstClsConut int
	InnerClsCount         int
	SkippedClassCount     int
	SkippedMethodCount    int
	PrivateMethodCount    int

	funcParents map[string]int // got 11 elements
}

var clts = &collects{funcParents: map[string]int{}}

func init() { clts.ClassSizeMap = map[int64]int{} }
func (this *collects) addClassSize(sz int64) {
	if sz <= 256 {
		return
	}
	if _, ok := this.ClassSizeMap[sz]; ok {
		this.ClassSizeMap[sz] += 1
	} else {
		this.ClassSizeMap[sz] = 1
	}
}

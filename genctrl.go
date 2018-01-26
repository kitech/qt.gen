package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-clang/v3.9/clang"

	"gopp"
)

var ast_file = "./qthdrsrc.ast"
var hdr_file = "./headers/qthdrsrc.h"

// module depend table
// 也许可以用ldd自动推导出来
var modDeps = map[string][]string{
	"core":              []string{},
	"gui":               []string{"core"},
	"widgets":           []string{"core", "gui"},
	"network":           []string{"core"},
	"qml":               []string{"core", "network"},
	"quick":             []string{"core", "gui", "network", "qml"},
	"quicktemplate2":    []string{"core", "gui", "network", "qml", "quick"},
	"quickcontrols2":    []string{"core", "gui", "network", "qml", "quick", "qiucktemplate2"},
	"quickwidgets":      []string{"core", "gui", "network", "qml", "quick", "widgets"},
	"multimedia":        []string{"core", "gui", "network"},
	"multimediawidgets": []string{"core", "gui", "network", "widgets", "multimedia", "opengl"},
	"opengl":            []string{"core", "gui", "widgets"},
	"sql":               []string{"core"},
	"svg":               []string{"core", "gui", "widgets"},
}
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

	filter    GenFilter
	genor     Generator
	qtenumgen Generator // for generate global enums
	qtfuncgen Generator // for generate global functions
	qttmplgen Generator // for generate template
}

func NewGenCtrl() *GenCtrl {
	this := &GenCtrl{}

	return this
}

func (this *GenCtrl) main() {
	this.setupLang()
	this.setupEnv()
	this.createTU()
	this.collectClasses()
	this.cleanupEnv()
}

func (this *GenCtrl) setupLang() {
	this.filter = &GenFilterInc{}
	this.genor = NewGenerateInline()

	if true {
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateGo()
		this.qtenumgen = NewGenerateGo()
		this.qtfuncgen = NewGenerateGo()
		this.qttmplgen = NewGenerateGo()
	}
}

func (this *GenCtrl) setupEnv() {
	cidx := clang.NewIndex(0, 1)
	// defer cidx.Dispose()

	modules := []string{
		"QtCore", "QtGui", "QtWidgets",
		"QtNetwork", "QtQml", "QtQuick",
	}

	cmdlines := []string{
		"-x c++ -std=c++11 -D__CODE_GENERATOR__ -D_GLIBCXX_USE_CXX11ABI=1",
		"-I/usr/include/qt -DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++11 -Wall -W -D_REENTRANT -fPIC",
	}
	args := []string{}
	gopp.Domap(cmdlines, func(e interface{}) interface{} {
		args = append(args, strings.Split(e.(string), " ")...)
		return nil
	})
	gopp.Domap(modules, func(e interface{}) interface{} {
		args = append(args, fmt.Sprintf("-DQT_%s_LIB", strings.ToUpper(e.(string)[2:])))
		args = append(args, fmt.Sprintf("-I/usr/include/qt/%s", e.(string)))
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
		tu = cidx.ParseTranslationUnit(hdr_file, args, nil, 0)
	}
	if !tu.IsValid() {
		log.Panicln("wtf")
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
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName(),

			cursor.SpecializedCursorTemplate().Spelling(), cursor.CanonicalCursor().Kind())
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
		this.qtfuncgen.(*GenerateGo).funcs = append(this.qtfuncgen.(*GenerateGo).funcs, cursor)
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
		if is_qt_class(cursor.Type()) {
			log.Println("got", cursor.Spelling(), ",", cursor.Type().Spelling(), ",", cursor.Type().Kind(), cursor.Type().ClassType().Spelling(), ",", cursor.Type().ClassType().Kind(), ",", cursor.Type().CanonicalType().Spelling(), ",", cursor.Type().CanonicalType().Kind(), ",", cursor.Type().CanonicalType().Declaration().Kind(), ",", cursor.Type().CanonicalType().Declaration().Spelling(), ",", cursor.Type().CanonicalType().Declaration().Definition().Kind().String())
			log.Println(cursor.Type().CanonicalType().Declaration().Kind().String(), cursor.Type().CanonicalType().Declaration().Spelling(), cursor.Type().CanonicalType().Declaration().IsCursorDefinition(), cursor.Spelling())
			this.qttmplgen.(*GenerateGo).tmplclsspecs = append(this.qttmplgen.(*GenerateGo).tmplclsspecs, cursor)

			tplc := cursor.Type().Declaration()
			log.Println(tplc.Spelling(), "============", cursor.CanonicalCursor().Kind(), tplc.NumTemplateArguments())
			/*
				tplc.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
					log.Println(c1.Kind().String(), c1.Spelling(), p1.Spelling(), cursor.Spelling(), c1.NumTemplateArguments())
					return clang.ChildVisit_Recurse
					// return clang.ChildVisit_Continue
				})
			*/

			if cursor.Spelling() == "QWidgetList" || cursor.Spelling() == "QByteArrayList" {
				// os.Exit(0)
			}
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
		this.qttmplgen.(*GenerateGo).tmplclses = append(this.qttmplgen.(*GenerateGo).tmplclses, cursor)
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName())
		log.Println(cursor.IsCursorDefinition(), cursor.NumTemplateArguments())
		/*
			cursor.Visit(func(c1, p1 clang.Cursor) clang.ChildVisitResult {
				log.Println(c1.Kind().String(), c1.Spelling(), p1.Spelling(), cursor.Spelling(), ",", cursor.Mangling())
				// return clang.ChildVisit_Recurse
				return clang.ChildVisit_Continue
			})
		*/
		if cursor.Spelling() == "QList" {
			// os.Exit(0)
		}
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
		this.qtenumgen.(*GenerateGo).enums = append(this.qtenumgen.(*GenerateGo).enums, cursor)
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
	case clang.Cursor_InvalidCode:
		fallthrough
	default:
		log.Println(cursor.Spelling(), ",", cursor.Kind().String(), ",", cursor.DisplayName(), cursor.Type().Spelling())
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
		}
		return clang.ChildVisit_Continue
	})

	cursor.Visit(this.visfn)

	var gf *GenerateGo = this.qtfuncgen.(*GenerateGo)
	gf.genFunctions(cursor, cursor.SemanticParent())

	var gg *GenerateGo = this.qtenumgen.(*GenerateGo)
	gg.cp.APf("header", "package qtcore")
	gg.genEnumsGlobal(cursor, cursor.SemanticParent())
	gg.saveCodeToFile("core", "qnamespace")

	this.qttmplgen.(*GenerateGo).genTemplateSpecializedClasses()

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

	TemplateClassCount int
	SkippedClassCount  int
	SkippedMethodCount int
	PrivateMethodCount int

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

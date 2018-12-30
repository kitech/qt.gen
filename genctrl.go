package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-clang/v3.9/clang"
	funk "github.com/thoas/go-funk"

	"gopp"
)

var ast_file = "./qthdrsrc.ast"
var bshdr_file = "./bsheaders/qthdrsrc.h"

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
	modlstgen  Generator
}

func NewGenCtrl() *GenCtrl {
	this := &GenCtrl{}

	return this
}

var genLang string = ""  // c(c binding), go (go binding), rs (rust binding)
var genQtdir string = "" // format: /home/me/Qt5.10.1 or /usr
var genQtver string = "" // format: 5.10.1

func (this *GenCtrl) main() {
	if len(os.Args) > 1 {
		genLang = os.Args[len(os.Args)-1]
	}
	if genLang == "" {
		// log.Println("optional set QTDIR env")
		// sometimes need use ulimit -n 10240
		log.Fatalln("must suply a lang to gen, usage: qt.gen <c|go|rs>")
	}

	this.setupQtinfo()
	btime := time.Now()
	qdi.load(genQtdir, genQtver)
	log.Println(time.Now().Sub(btime))
	// log.Fatalln("test exit")

	this.setupLang()
	this.setupEnv()
	this.createTU()
	this.collectClasses()
	this.cleanupEnv()
}

func (this *GenCtrl) setupLang() {

	switch genLang {
	case "c":
		this.filter = NewGenFilterInc()
		this.genor = NewGenerateInline(genQtdir, genQtver)
		this.qtenumgen = NewGenerateInline(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateInline(genQtdir, genQtver)
		this.qttmplgen = NewGenerateInline(genQtdir, genQtver)
		this.qtconstgen = NewGenerateInline(genQtdir, genQtver)
	case "go":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateGo(genQtdir, genQtver)
		this.qtenumgen = NewGenerateGo(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateGo(genQtdir, genQtver)
		this.qttmplgen = NewGenerateGo(genQtdir, genQtver)
		this.qtconstgen = NewGenerateGo(genQtdir, genQtver)
	case "rs":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateRs(genQtdir, genQtver)
		this.qtenumgen = NewGenerateRs(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateRs(genQtdir, genQtver)
		this.qttmplgen = NewGenerateRs(genQtdir, genQtver)
		this.qtconstgen = NewGenerateRs(genQtdir, genQtver)
		this.modlstgen = NewGenerateRs(genQtdir, genQtver)
	case "jl":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateJl(genQtdir, genQtver)
		this.qtenumgen = NewGenerateJl(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateJl(genQtdir, genQtver)
		this.qttmplgen = NewGenerateJl(genQtdir, genQtver)
		this.qtconstgen = NewGenerateJl(genQtdir, genQtver)
		this.modlstgen = NewGenerateJl(genQtdir, genQtver)
	case "cr":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateCr(genQtdir, genQtver)
		this.qtenumgen = NewGenerateCr(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateCr(genQtdir, genQtver)
		this.qttmplgen = NewGenerateCr(genQtdir, genQtver)
		this.qtconstgen = NewGenerateCr(genQtdir, genQtver)
		this.modlstgen = NewGenerateCr(genQtdir, genQtver)
	case "dt":
		this.filter = &GenFilterGo{}
		this.genor = NewGenerateDt(genQtdir, genQtver)
		this.qtenumgen = NewGenerateDt(genQtdir, genQtver)
		this.qtfuncgen = NewGenerateDt(genQtdir, genQtver)
		this.qttmplgen = NewGenerateDt(genQtdir, genQtver)
		this.qtconstgen = NewGenerateDt(genQtdir, genQtver)
		this.modlstgen = NewGenerateDt(genQtdir, genQtver)
		// fallthrough
	default:
		log.Fatalln("not supported or not impled:", genLang, genQtdir, genQtver)
	}
}

func (this *GenCtrl) setupQtinfo() {
	qtdir := gopp.IfElseStr(os.Getenv("QTDIR") == "", "/usr", os.Getenv("QTDIR"))
	qtver := ""
	if qtdir == "/usr" {
	} else if strings.HasPrefix(qtdir, "qtheaders") {
	} else {
		log.Println(qtdir)
		reg := `Qt([0-9.]+)`
		exp := regexp.MustCompile(reg)
		mats := exp.FindAllStringSubmatch(qtdir, -1)
		log.Println(mats)
		qtver = mats[0][1]
		if gopp.FileExist2(fmt.Sprintf("%s/%s", qtdir, qtver)) {
			// for 5.12.0
		} else {
			qtver = gopp.IfElseStr(strings.HasSuffix(qtver, ".0"), qtver[:len(qtver)-2], qtver)
		}
	}
	genQtdir, genQtver = qtdir, qtver
	log.Println("qt info:", qtdir, qtver, os.Getenv("QTDIR"))

	rebuildModDepsAll(qtver)
}

func (this *GenCtrl) setupEnv() {

	// 预先处理头文件, cd gcc_64/include/ && ln -sv ../../android_x86/include/QtAndroidExtras
	// 预先处理头文件, cd gcc_64/include/ && ln -sv ../../Src/qtwinextras/include/QtWinExtras
	// 预先处理头文件, cd gcc_64/include/ && ln -sv ../../Src/qtmacextras/include/QtMacExtras
	// 这是要生成的模块表
	modules := []string{
		"QtCore", "QtGui", "QtWidgets",
		"QtNetwork", "QtQml", "QtQuick",
		"QtQuickTemplates2", "QtQuickControls2", "QtQuickWidgets",
		// for platform dependent modules, need copy headers if not exists
		"QtAndroidExtras", // fatal error: 'jni.h' file not found, link /opt/android-ndk/sysroot/usr/include/jni.h -> bsheaders/jni.h
		"QtX11Extras",     // 这个包没生成出来什么代码,
		"QtWinExtras",     // 缺少QtWinExtracsDepened头文件,link qt-opensource-linux.bin installs to gcc_64
		"QtMacExtras",     // 缺少QtMacExtracsDepened头文件
		// webengines
		"QtPositioning", "QtWebChannel", "QtWebEngineCore", "QtWebEngine", "QtWebEngineWidgets",
		// multimedia
		"QtSvg", "QtMultimedia",
	}
	// modules = []string{"QtCore", "QtGui", "QtWidgets"} // for test

	cmdlines := []string{
		"-x c++ -std=c++11 -D__CODE_GENERATOR__ -D_GLIBCXX_USE_CXX11ABI=1",
		"-DQT_NO_DEBUG -D_GNU_SOURCE -pipe -fno-exceptions -O2 -march=x86-64 -mtune=generic -O2 -pipe -fstack-protector-strong -std=c++11 -Wall -W -D_REENTRANT -fPIC",
		"-DQT_OPENGL_ES_2", "-DQT_OPENGL_ES_3",
		// "-DQ_CLANG_QDOC", // 开启QDOC，竟然会出错
		"-I./bsheaders", "-I/usr/include/wine/windows/", // fix cross platform generate, win/mac
	}

	args := []string{}
	gopp.Domap(cmdlines, func(e interface{}) interface{} {
		args = append(args, strings.Split(e.(string), " ")...)
		return nil
	})

	// this.setupQtinfo()
	qtdir, qtver := genQtdir, genQtver
	qtsysdir := qtdir
	if qtdir == "/usr" {
		args = append(args, fmt.Sprintf("-I%s/include/qt", qtdir))
	} else if strings.HasPrefix(qtdir, "qtheaders") {
		args = append(args, fmt.Sprintf("-I./qtheaders/include")) //depcreated self construct header tree
	} else {
		log.Println(qtdir)
		args = append(args, fmt.Sprintf("-I%s/%s/gcc_64/include", qtdir, qtver))
		qtsysdir += fmt.Sprintf("/%s/gcc_64", qtver)
	}
	log.Println("qt info:", qtdir, qtver, os.Getenv("QTDIR"))
	if !gopp.FileExist2(qtsysdir) {
		log.Fatalln("maybe QTDIR not exists error", qtdir, qtver, qtsysdir)
	}

	hdrdirok := true
	gopp.Domap(modules, func(e interface{}) interface{} {
		args = append(args, fmt.Sprintf("-DQT_%s_LIB", strings.ToUpper(e.(string)[2:])))
		args = append(args, fmt.Sprintf("-DGEN_GO_QT_%s_LIB", strings.ToUpper(e.(string)[2:])))
		if qtdir == "/usr" {
			args = append(args, fmt.Sprintf("-I/usr/include/qt/%s", e.(string)))
			_, err := os.Stat(fmt.Sprintf("/usr/include/qt/%s", e.(string)))
			gopp.ErrPrint(err)
			hdrdirok = gopp.IfElse(err == nil, hdrdirok, false).(bool)
		} else if strings.HasPrefix(qtdir, "qtheaders") {
			args = append(args, fmt.Sprintf("-I./qtheaders/include/%s", e.(string)))
		} else {
			args = append(args, fmt.Sprintf("-I%s/%s/gcc_64/include/%s", qtdir, qtver, e.(string)))
			_, err := os.Stat(fmt.Sprintf("%s/%s/gcc_64/include/%s", qtdir, qtver, e.(string)))
			gopp.ErrPrint(err)
			hdrdirok = gopp.IfElse(err == nil, hdrdirok, false).(bool)
		}
		return nil
	})
	if !hdrdirok {
		log.Fatalln("Some header dir(s) not exists")
	}

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
	fullCmd := fmt.Sprintf("g++ %s -o qthdrsrc.o -c %s", strings.Join(args, " "), bshdr_file)
	ioutil.WriteFile("bcmd.sh", []byte(fullCmd), 0755)
	// os.Exit(0)

	this.args = args
	this.modules = modules
}

func (this *GenCtrl) dryrunEnv() {
	log.Println("sh ./bcmd.sh ...")
	cmdo := exec.Command("sh", "./bcmd.sh")
	output, err := cmdo.CombinedOutput()
	gopp.ErrFatal(err, string(output))
}

func (this *GenCtrl) createTU() {
	cidx := clang.NewIndex(0, 1)
	// defer cidx.Dispose()
	this.cidx = cidx
	args := this.args

	var tu clang.TranslationUnit
	save_ast := false
	opts := uint32(0)
	if _, err := os.Stat(ast_file); err == nil {
		tu = cidx.TranslationUnit(ast_file)
	} else {
		save_ast = true
		opts |= clang.TranslationUnit_DetailedPreprocessingRecord
		opts |= clang.TranslationUnit_IncludeBriefCommentsInCodeCompletion
		opts |= clang.TranslationUnit_CreatePreambleOnFirstParse
		tu = cidx.ParseTranslationUnit(bshdr_file, args, nil, opts)
		// 需要正常编译能够通过
	}
	// log.Println(tu.Spelling(), tu.NumDiagnostics(), tu.DefaultReparseOptions(), opts)
	if !tu.IsValid() {
		log.Panicln("wtf", "maybe cached qthdrsrc.ast file expired, delete and retry please.", tu.NumDiagnostics())
	}
	cursor := tu.TranslationUnitCursor()
	if false {
		log.Println(cursor)
	}

	this.tuc = cursor
	this.tu = tu
	this.save_ast = save_ast

	if this.save_ast {
		this.tu.SaveTranslationUnit(ast_file, 0)
	}
}

func (this *GenCtrl) visfn(cursor, parent clang.Cursor) clang.ChildVisitResult {
	{
		log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName(), cursor.SpecializedCursorTemplate().DisplayName(), cursor.CanonicalCursor().Kind(), cursor.IsCursorDefinition(), cursor.Language(), cursor.Linkage(), cursor.HasAttrs(), cursor.Extent(), parent.Spelling(), get_decl_loc(parent), cursor.SemanticParent().Spelling())
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
			if genLang == "rs" {
				this.modlstgen.(*GenerateRs).genModLst(cursor)
			}
			if genLang == "dt" {
				this.modlstgen.(*GenerateDt).genModLst(cursor)
			}
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
			if genLang == "rs" {
				this.modlstgen.(*GenerateRs).genModLst(cursor)
			}
			if genLang == "dt" {
				this.modlstgen.(*GenerateDt).genModLst(cursor)
			}
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
		// clang reports 'extern "C" ...' as an unexposed decl, which we definitely
		// need to recurse into.
		cursor.Visit(this.visfn)
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

	/*
		cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
			file, lineno, _, _ := cursor.Location().SpellingLocation()
			log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName(), cursor.SpecializedCursorTemplate().DisplayName(), cursor.CanonicalCursor().Kind(), cursor.IsCursorDefinition(), cursor.Language(), cursor.Linkage(), cursor.HasAttrs(), cursor.Extent(), parent.Spelling(), get_decl_loc(parent), cursor.SemanticParent().Spelling(), file, file.Name(), lineno)
			switch cursor.Kind() {
			case clang.Cursor_CXXMethod, clang.Cursor_Constructor, clang.Cursor_FunctionDecl:
				return clang.ChildVisit_Continue
			case clang.Cursor_MacroExpansion:
				if funk.ContainsString([]string{"Q_SIGNALS", "Q_SLOTS"}, cursor.Spelling()) {
					clts.macroExpands[file.Name()] = append(clts.macroExpands[file.Name()], cursor)
				}
			case clang.Cursor_CXXAccessSpecifier:
				clts.macroExpands[file.Name()] = append(clts.macroExpands[file.Name()], cursor)
			}
			return clang.ChildVisit_Recurse
		})
	*/

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
		case clang.Cursor_MacroExpansion:
			if strings.Contains(cursor.Spelling(), "QT_REQUIRE_CONFIG") {
				// log.Println(cursor.Spelling(), get_decl_loc(cursor))
				srcfile := get_decl_loc(cursor)
				fbname := strings.Split(filepath.Base(srcfile), ":")[0]
				fbmod := filepath.Base(filepath.Dir(srcfile))
				fbpath := strings.ToLower(fmt.Sprintf("%s/%s", fbmod, fbname))

				csline := readCursorLine(cursor)
				featname := strings.Trim(csline[len(cursor.Spelling()):], "();")
				clts.qtreqcfgs[fbpath] = featname

			}
			// case clang.Cursor_MacroInstantiation:
		}
		return clang.ChildVisit_Continue
	})

	cursor.Visit(this.visfn)

	this.qtfuncgen.genFunctions(cursor, cursor.SemanticParent())
	if genLang == "go" {
		var gg *GenerateGo = this.qtenumgen.(*GenerateGo)
		gg.cp.APf("header", "package qtcore")
		gg.cp.APf("header", "import \"fmt\"")
		gg.genEnumsGlobal(cursor, cursor.SemanticParent())
		gg.cp.APf("keep", "func make_sure_usepkg_qnamespace(){if false{fmt.Println(123)}}")
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
	} else if genLang == "rs" {
		var gg *GenerateRs = this.modlstgen.(*GenerateRs)
		for modname, cp := range gg.cpcs {
			log.Println(modname, cp.TotolLine(), cp.TotolLength())
			gg.saveCodeToFileWithCode(modname, "lib", cp.ExportAll())
		}
	} else if genLang == "jl" {
	} else if genLang == "cr" {
	} else if genLang == "dt" {
		var gg *GenerateDt = this.modlstgen.(*GenerateDt)
		for modname, cp := range gg.cpcs {
			log.Println(modname, cp.TotolLine(), cp.TotolLength())
			gg.saveCodeToFileWithCode(modname, "qt"+modname, cp.ExportAll())
		}
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

	macroExpands map[string][]clang.Cursor // file name => macro explan, and clang.Cursor_CXXAccessSpecifier
	qtreqcfgs    map[string]string         // file name QtWidgets/qdialog.h => feature name
}

var clts = &collects{funcParents: map[string]int{}}

func init() {
	clts.ClassSizeMap = map[int64]int{}
	clts.macroExpands = map[string][]clang.Cursor{}
	clts.qtreqcfgs = map[string]string{}
}
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

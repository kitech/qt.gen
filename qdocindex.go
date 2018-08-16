package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-clang/v3.9/clang"
	"github.com/therecipe/qt/internal/binding/parser"
)

var qdi *QDocIndex = newQDocIndex()

type QDocIndex struct {
	loaded bool
}

func newQDocIndex() *QDocIndex {
	this := &QDocIndex{}
	return this
}

func (this *QDocIndex) load(qtdir, qtver string) {
	if this.loaded {
		return
	}
	// os.Setenv("QT_DIR", os.Getenv("HOME")+"/Qt5.10.1")
	// os.Setenv("QT_VERSION", "5.10.1")
	os.Setenv("QT_DIR", qtdir)
	os.Setenv("QT_VERSION", qtver)
	parser.LoadModules()
	this.loaded = true
}

func (this *QDocIndex) run(qtdir, qtver string) {
	this.load(qtdir, qtver)
	for i, lib := range parser.GetLibs() {
		log.Println(i, lib)
	}

	clsos := parser.SortedClassesForModule("QtCore", true)
	log.Println(len(clsos))
	for i, clso := range clsos {
		log.Println(i, clso.Name, clso.Status, clso.Module)
		for j, funco := range clso.Functions {
			log.Println(i, j, clso.Name, funco.Name, funco.IsSupported(), funco.Signature)
		}
		for j, enumo := range clso.Enums {
			log.Println(i, j, clso.Name, enumo.Name, enumo.Fullname, enumo.Values)
		}
		for j, propo := range clso.Properties {
			log.Println(i, j, clso.Name, propo.Name, propo.Fullname, propo.Getter)
		}
	}
	log.Println(len(clsos))
}

func (this *QDocIndex) findClass(class string) (*parser.Class, bool) {
	for _, module := range parser.GetLibs() {
		log.Println("module class count: ", module, len(parser.SortedClassNamesForModule(module, true)), len(parser.GetLibs()), len(parser.SortedClassesForModule(module, true)))
		module = "Qt" + module

		for _, clso := range parser.SortedClassesForModule(module, true) {
			if clso.Name == class {
				return clso, true
			}
		}
	}

	return nil, false
}

func (this *QDocIndex) findCoMethodDoc(cursor clang.Cursor) (funco *parser.Function) {
	return
}

func (this *QDocIndex) findCoMethodCursor(clscs clang.Cursor, funco *parser.Function) (mthcs clang.Cursor, found bool) {
	clscs.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_CXXMethod, clang.Cursor_Constructor, clang.Cursor_Destructor:
			if cursor.Spelling() != funco.Name {
				break
			}
			if cursor.NumArguments() != int32(len(funco.Parameters)) {
				break
			}
			if cursor.CXXMethod_IsConst() && !strings.HasSuffix(funco.Signature, " const") {
				break
			}
			if !cursor.CXXMethod_IsConst() && strings.HasSuffix(funco.Signature, " const") {
				break
			}
			// cursor.CXXMethod_IsConst()
			log.Println("...", cursor.NumArguments(), cursor.CXXMethod_IsConst(), funco.Name, funco.Signature, cursor.DisplayName(), cursor.Mangling())
			match := true
			for idx := int32(0); idx < cursor.NumArguments(); idx++ {
				argod := funco.Parameters[idx]
				argoc := cursor.Argument(uint32(idx))
				log.Printf("%s, %+v\n", argoc.Type().Spelling(), argod)
				if argod.Value != argoc.Type().Spelling() {
					match = false
					break
				}
			}
			if match {
				mthcs = cursor
				found = true
			}
		}
		return clang.ChildVisit_Continue
	})
	return
}

func (this *QDocIndex) findCoClassCursor(clso *parser.Class) (clscs clang.Cursor, found bool) {
	/*
		if cs, ok := allClasses[clso.Name]; ok {
			clscs = cs
			found = true
		}
	*/
	return
}

func (this *QDocIndex) findCoMethodObj(mthc clang.Cursor) (funco *parser.Function, found bool) {
	clsc := mthc.SemanticParent()
	clsName := clsc.Spelling()

	clso, found := this.findClass(clsName)
	if !found {
		return nil, false
	}

	funco, found = this.findMethod(clso, mthc)
	return
}

func (this *QDocIndex) findMethod(clso *parser.Class, mthc clang.Cursor) (funco *parser.Function, found bool) {
	funcName := mthc.Spelling()

	for _, funco_ := range clso.Functions {
		if funco_.Name != funcName {
			continue
		}
		if len(funco_.Parameters) != int(mthc.NumArguments()) {
			continue
		}

		if mthc.CXXMethod_IsConst() && !strings.HasSuffix(funco_.Signature, " const") {
			continue
		}
		if !mthc.CXXMethod_IsConst() && strings.HasSuffix(funco_.Signature, " const") {
			continue
		}

		match := true
		for idx := int32(0); idx < mthc.NumArguments(); idx++ {
			argod := funco_.Parameters[idx]
			argoc := mthc.Argument(uint32(idx))
			log.Printf("%s, %+v\n", argoc.Type().Spelling(), argod)
			if argod.Value != argoc.Type().Spelling() {
				match = false
				break
			}
		}
		if match {
			funco = funco_
			found = true
			return
		}
	}

	return
}

func (*QDocIndex) protoMatch(funco_ *parser.Function, mthc clang.Cursor) bool {
	funcName := mthc.Spelling()

	if funco_.Name != funcName {
		return false
	}
	if len(funco_.Parameters) != int(mthc.NumArguments()) {
		return false
	}

	if mthc.CXXMethod_IsConst() && !strings.HasSuffix(funco_.Signature, " const") {
		return false
	}
	if !mthc.CXXMethod_IsConst() && strings.HasSuffix(funco_.Signature, " const") {
		return false
	}

	match := true
	for idx := int32(0); idx < mthc.NumArguments(); idx++ {
		argod := funco_.Parameters[idx]
		argoc := mthc.Argument(uint32(idx))
		log.Printf("%s, %+v\n", argoc.Type().Spelling(), argod)
		if argod.Value != argoc.Type().Spelling() {
			match = false
			break
		}
	}
	if match {
		return true
	}
	return false
}

// TODO
func (this *QDocIndex) findCoFunction(fc clang.Cursor) (funco *parser.Function, found bool) {
	for _, lib := range parser.GetLibs() {
		_ = lib
	}
	return
}

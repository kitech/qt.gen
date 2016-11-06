package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/go-clang/v3.9/clang"
	"github.com/kitech/colog"

	"gopp"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Flags())

	colog.Register()
}

func main() {
	cidx := clang.NewIndex(0, 1)
	defer cidx.Dispose()

	modules := []string{
		"QtCore", "QtGui", "QtWidgets",
	}

	hdrsrc := "./headers/qthdrsrc.h"
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
	cmd := exec.Command("gcc", "--print-file-name=include")
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	args = append(args, fmt.Sprintf("-I%s", string(out[:len(out)-1])))
	args = append(args, fmt.Sprintf("-I%s-fixed", string(out[:len(out)-1])))
	args = append(args, "-I/usr/include/c++/6.2.1")
	log.Println(args)

	tu := cidx.ParseTranslationUnit(hdrsrc, args, nil, 0)
	log.Println(tu)
	cursor := tu.TranslationUnitCursor()
	if false {
		log.Println(cursor)
	}
	cnter := 0
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl:
		case clang.Cursor_FunctionDecl:
		case clang.Cursor_InvalidCode:
			fallthrough
		default:
			log.Println(cursor.Spelling(), cursor.Type().Kind().String())
		}

		cnter += 1
		return clang.ChildVisit_Continue
	})
	log.Println(cnter)
}

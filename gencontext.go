package main

import (
	"github.com/go-clang/v3.9/clang"
	"github.com/therecipe/qt/internal/binding/parser"
)

type GenContext struct {
}

type GenClassContext struct {
	clscs clang.Cursor
	clso  *parser.Class

	clsty    clang.Type
	clscanty clang.Type
	bclses   []clang.Cursor

	methods            []clang.Cursor
	funcs              []clang.Cursor
	tmplclses          []clang.Cursor
	plaintmplinstclses []clang.Cursor
	tydeftmplinstclses []clang.Cursor
	enums              []clang.Cursor
	constants          []clang.Cursor

	isPureVirtualClass  bool
	hasVirtualProtected bool
	isQObjectClass      bool
	isDeletedClass      bool
	hasProjectedDtor    bool
	hasExplictDtor      bool

	funcMangles map[string]int

	// method indexes, reset perclass
	mthidxs map[string]int
}

func (this *GenClassContext) walkMth(f func(ctx *GenMethodContext, cursor, parent clang.Cursor)) {

}

type GenMethodContext struct {
	isFunc bool // just a function

	isStatic    bool
	isConst     bool
	isPublic    bool
	isProtected bool
	isPrivate   bool

	hasImpl    bool
	isVirt     bool
	isPureVirt bool

	olidx int // overload index

	resty    clang.Type
	rescanty clang.Type

	argctxs []*GenArgumentContext

	_argDescs   []string
	_paramDescs []string
	_argtyDescs []string
}

func (this *GenMethodContext) walkArg(f func(ctx *GenArgumentContext, argcs, cursor, parent clang.Cursor)) {

}

type GenArgumentContext struct {
	idx      int
	argty    clang.Type
	argcanty clang.Type
}

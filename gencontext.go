package main

import (
	"github.com/go-clang/v3.9/clang"
	// "github.com/therecipe/qt/internal/binding/parser"
)

// should be Options???
type GenContext struct {
	qtdir      string
	qtver      string
	genlang    string
	// filter clip, default true
	// full code generate for
	noclip 		bool // use filter clip now
	specify_class string
	bsast_file string
	bshdr_file string

	debug int
	// generate C_ prefix wrap func, for avoid some func ROV cannot correct handled
	// default true
	// thus we can merge geninc.go and genincv0 in one
	cwrap bool
	// default false
	refmt_gened_code bool
}
func NewGenContext() *GenContext {
	rv := &GenContext{}
	rv.debug = 1
	rv.cwrap = true
	return rv
}

type GenClassContext struct {
	clscs clang.Cursor
	// clso  *parser.Class

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

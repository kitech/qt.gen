package main

import (
	"github.com/go-clang/v3.9/clang"
)

type TypeConvertor interface {
	toDest(clang.Type, clang.Cursor) string // dest language
	toBind(clang.Type, clang.Cursor) string // binding type
}

// ???
type ValueConvertor interface {
}

type TypeConvertGo struct {
}

func NewTypeConvertGo() *TypeConvertGo {
	this := &TypeConvertGo{}
	return this
}

func (this *TypeConvertGo) toDest(ty clang.Type, cursor clang.Cursor) string {
	return ""
}

func (this *TypeConvertGo) toBind(ty clang.Type, cursor clang.Cursor) string {
	return ""
}

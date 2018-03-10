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
}

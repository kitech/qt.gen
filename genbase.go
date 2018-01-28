package main

import (
	"log"

	"github.com/go-clang/v3.9/clang"
)

type Generator interface {
	// init(cursor, parent clang.Cursor)
	// genPassHeader(cursor, parent clang.Cursor)
	genClass(cursor, parent clang.Cursor)
}

func init() {
	if false {
		log.Println("hehre")
	}
}

type GenBase struct {
	tu *clang.TranslationUnit

	pureVirtual bool
}

// TODO is what?
func (this *GenBase) isSignal() bool {
	return false
}

func (this *GenBase) isSlot() bool {
	return false
}

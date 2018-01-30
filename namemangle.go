package main

import (
	"fmt"

	"github.com/go-clang/v3.9/clang"
)

type GenMangler interface {
	convTo(cursor clang.Cursor) string
	// convFrom(mname string) string
	origin(cursor clang.Cursor) string
}

type IncMangler struct {
}

func NewIncMangler() *IncMangler {
	this := &IncMangler{}
	return this
}

func (this *IncMangler) convTo(cursor clang.Cursor) string {
	if false {
		// C1/C2/C3 for case
		fmt.Println("what's the manglings:", cursor.Manglings().Strings())
	}
	return fmt.Sprintf("C%s", this.origin(cursor))
}

func (this *IncMangler) origin(cursor clang.Cursor) string {
	if false {
		// C1/C2/C3 for case
		fmt.Println("what's the manglings:", cursor.Manglings().Strings())
	}
	if cursor.Manglings().Count() > 1 {
		return cursor.Manglings().Strings()[0].String()
	}
	return cursor.Mangling()
}

type GoMangler struct {
}

func NewGoMangler() *GoMangler {
	this := &GoMangler{}
	return this
}

func (this *GoMangler) convTo(cursor clang.Cursor) string {
	if false {
		// C1/C2/C3 for case
		fmt.Println("what's the manglings:", cursor.Manglings().Strings())
	}
	return fmt.Sprintf("C%s", this.origin(cursor))
}

func (this *GoMangler) origin(cursor clang.Cursor) string {
	if false {
		// C1/C2/C3 for case
		fmt.Println("what's the manglings:", cursor.Manglings().Strings())
	}
	if cursor.Manglings().Count() > 1 {
		return cursor.Manglings().Strings()[0].String()
	}
	return cursor.Mangling()
}

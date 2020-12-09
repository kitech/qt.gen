package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

var mgpfx = "qm" // qm12345 as qt function

type GenMangler interface {
	convTo(cursor clang.Cursor) string
	// convFrom(mname string) string
	origin(cursor clang.Cursor) string
	crc32p(cursor clang.Cursor) string
	crc32(cursor clang.Cursor) string
	// crc32hex(cursor clang.Cursor) string
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

func (this *IncMangler) origin(cursor clang.Cursor) (defname string) {
	if false {
		// C1/C2/C3 for case
		fmt.Println("what's the manglings:", cursor.Manglings().Strings())
	}

	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		return strings.Replace(cursor.Mangling(), "C1E", "C2E", -1)
	case clang.Cursor_Destructor:
		return strings.Replace(cursor.Mangling(), "D1Ev", "D2Ev", -1)
	}
	return cursor.Mangling()
}
func (this *IncMangler) crc32p(cursor clang.Cursor) string {
	return fmt.Sprintf("%s%d", mgpfx, symcrc32(this.origin(cursor)))
}
func (this *IncMangler) crc32(cursor clang.Cursor) string {
	return fmt.Sprintf("%d", symcrc32(this.origin(cursor)))
}

// 不用hex格式，其他语言要处理格式麻烦
func symcrc32(s string) uint32 { return crc32.ChecksumIEEE([]byte(s)) }

type GoMangler struct {
}

func NewGoMangler() *GoMangler {
	this := &GoMangler{}
	return this
}

func (this *GoMangler) convTo(cursor clang.Cursor) string {
	if false {
		// C1/C2/C3 for case
		log.Println("what's the manglings:", cursor.Manglings().Strings())
	}
	return fmt.Sprintf("C%s", this.origin(cursor))
}

func (this *GoMangler) origin(cursor clang.Cursor) (defname string) {
	if false {
		// C1/C2/C3 for case
		log.Println("what's the manglings:", cursor.Manglings().Strings())
	}

	switch cursor.Kind() {
	case clang.Cursor_Constructor:
		return strings.Replace(cursor.Mangling(), "C1E", "C2E", -1)
	case clang.Cursor_Destructor:
		return strings.Replace(cursor.Mangling(), "D1Ev", "D2Ev", -1)
	}
	return cursor.Mangling()
}
func (this *GoMangler) crc32p(cursor clang.Cursor) string {
	return fmt.Sprintf("%s%d", mgpfx, symcrc32(this.origin(cursor)))
}
func (this *GoMangler) crc32(cursor clang.Cursor) string {
	return fmt.Sprintf("%d", symcrc32(this.origin(cursor)))
}

package main

import (
	"log"

	"github.com/go-clang/v3.9/clang"
)

type GenerateInline struct {
	filter GenFilter

	cursor, parent clang.Cursor
	methods        map[clang.Cursor]clang.Cursor
	cp             *CodePager
}

func NewGenerateInline() *GenerateInline {
	this := &GenerateInline{}
	this.filter = &GenFilterInc{}
	this.cp = NewCodePager()
	blocks := []string{"header", "main", "use", "ext", "body"}
	for _, block := range blocks {
		this.cp.AddPointer(block)
	}

	return this
}

func (this *GenerateInline) genClass(cursor, parent clang.Cursor) {
	this.cursor, this.parent = cursor, parent

	log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	file, line, col, _ := cursor.Location().FileLocation()
	log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())

	this.walkClass()
	this.genMethods()
	this.final()
}

func (this *GenerateInline) final() {
	log.Println(this.cp.ExportAll())
}

func (this *GenerateInline) walkClass() {
	cursor, _ := this.cursor, this.parent

	methods := make(map[clang.Cursor]clang.Cursor, 0)

	// pcursor := cursor
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_Constructor:
			fallthrough
		case clang.Cursor_Destructor:
			fallthrough
		case clang.Cursor_CXXMethod:
			if !this.filter.skipMethod(cursor, parent) {
				methods[cursor] = parent
			} else {
				// log.Println("filtered:", cursor.Spelling())
			}
		case clang.Cursor_UnexposedDecl:
			// log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			file, line, col, _ := cursor.Location().FileLocation()
			if false {
				log.Println(file.Name(), line, col, file.Time())
			}
		default:
			if false {
				log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
			}
		}
		return clang.ChildVisit_Continue
	})

	this.methods = methods
}

func (this *GenerateInline) genMethods() {
	log.Println("process class:", len(this.methods), this.cursor.Spelling())
	for cursor, _ := range this.methods {
		// log.Println(cursor.Kind().String(), cursor.DisplayName())
		if cursor.Spelling() == "QString" {
			log.Println("ctor", cursor.DisplayName())
		}
	}
}

func (this *GenerateInline) genCtor() {

}

func (this *GenerateInline) genDtor() {

}

func (this *GenerateInline) genNonStaticMethod() {

}

func (this *GenerateInline) genStaticMethod() {

}

package main

import (
	"log"

	"github.com/go-clang/v3.9/clang"
)

type GenerateInline struct {
}

func NewGenerateInline() *GenerateInline {
	this := &GenerateInline{}
	return this
}

func (this *GenerateInline) GenClass(cursor clang.Cursor) {
	// pcursor := cursor
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_CXXMethod:
			if cursor.AccessSpecifier() != clang.AccessSpecifier_Public {
				break
			}
			// 判断是否是类内声明，看来是不需要
			log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName(),
				cursor.AccessSpecifier().String(), cursor.IsFunctionInlined())
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
	log.Println(cursor.Spelling(), cursor.Kind().String(), cursor.DisplayName())
	file, line, col, _ := cursor.Location().FileLocation()
	log.Printf("%s:%d:%d @%s\n", file.Name(), line, col, file.Time().String())
}

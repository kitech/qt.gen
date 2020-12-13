package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-clang/v3.9/clang"
)

var isTmplSelfRef = func(str string, parent clang.Cursor) bool {
	log.Println(parent.DisplayName(), str)
	reg := `.*(Q.+)<([KeyTKV, ]+)>*.`
	exp := regexp.MustCompile(reg)
	mats1 := exp.FindAllStringSubmatch(str, -1)
	mats2 := exp.FindAllStringSubmatch(parent.DisplayName(), -1)
	log.Println(mats1, mats2)
	return len(mats1) > 0 && len(mats2) > 0 && mats1[0][1] == mats2[0][1] &&
		len(strings.Split(mats1[0][2], ",")) == len(strings.Split(mats2[0][2], ","))
	// return strings.Contains(str, parent.DisplayName()) ||
	//	strings.Contains(parent.DisplayName(), str)
	// return strings.Contains(str, parent.Spelling()+"<T>")
}
var isTmplElemRef = func(ty clang.Type) bool {
	return ty.Spelling() == "T" || ty.Spelling() == "const T" ||
		ty.PointeeType().Spelling() == "T" || ty.PointeeType().Spelling() == "const T"
}
var isTmplIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::iterator")
}
var isTmplConstIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_iterator")
}
var isTmplRIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::reverse_iterator")
}
var isTmplConstRIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_reverse_iterator")
}
var isTmplKIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::key_iterator")
}
var isTmplConstKIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_key_iterator")
}
var isTmplKVIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::key_value_iterator")
}
var isTmplConstKVIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_key_value_iterator")
}
var isTmplAnyIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::") && strings.Contains(str, "iterator")
}
var isTmplConstAnyIterRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::const_") && strings.Contains(str, "iterator")
}
var isTmplNodeRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, parent.Spelling()+"::Node")
}
var isTmplKeyRef = func(str string, parent clang.Cursor) bool {
	return strings.Contains(str, "Key")
}

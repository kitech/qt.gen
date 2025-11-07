package main

import (
	"github.com/go-clang/v3.9/clang"
)

type Typexx struct {
	clang.Type // origin

	tydef clang.Type
	canon clang.Type
	class clang.Type

	bare clang.Type
}

func TypexxNew(ty clang.Type) *Typexx {
	if ty.Kind() == clang.Type_Elaborated {
		ty = ty.Declaration().Type()
	}
	tyx := &Typexx{Type:ty}
	tyx.tydef = ty.CanonicalType()
	tyx.canon = ty.CanonicalType()
	tyx.class = ty.ClassType() // what about not class?
	tyx.bare = get_bare_type(ty).CanonicalType()
	return tyx
}
func (tyx *Typexx) IsTypedef() bool {
	// sometimes it's clang.Type_Typedef
	return tyx.Declaration().Kind() == clang.Type_Typedef
}
// func (tyx *Typexx) IsPrimitive() bool {
// }

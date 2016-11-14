package main

import (
	"log"

	"github.com/kitech/colog"
)

func init() {
	colog.Register()
	colog.SetFlags(log.LstdFlags | log.Lshortfile | log.Flags())
}

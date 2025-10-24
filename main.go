package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	ctrl := NewGenCtrl()
	ctrl.main()
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

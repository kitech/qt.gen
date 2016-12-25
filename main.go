package main

import (
	"flag"
)

func main() {
	flag.Parse()
	ctrl := NewGenCtrl()
	ctrl.main()
}

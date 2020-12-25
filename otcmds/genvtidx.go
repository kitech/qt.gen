package main

import (
	"fmt"
	"gopp"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

func main() {
	mods := []string{"Core", "Gui", "Widgets"}
	// mods = []string{"Widgets"}
	for _, mod := range mods {
		libpath := fmt.Sprintf("/usr/lib/libQt5%s.so", mod)

		cmdo := exec.Command("vtable-dumper", libpath)
		outputb, err := cmdo.CombinedOutput()
		gopp.ErrPrint(err)
		output := string(outputb)
		outputb = nil
		lines := strings.Split(output, "\n")
		output = ""

		log.Println(len(lines))
		// log.Println(output)

		code := ""
		code += fmt.Sprintf("package qt%s\n", strings.ToLower(mod))
		code += fmt.Sprintf("import \"github.com/kitech/qt.go/qtrt\"\n")
		code += fmt.Sprintf("func init() {\n")
		code += fmt.Sprintf("  qtrt.RegisterVtableIndexFunc(get_vtable_index)\n")
		code += fmt.Sprintf("}\n")
		code += fmt.Sprintf("func get_vtable_index(clsname, mtname string) int {\n")
		code += fmt.Sprintf("switch clsname {\n")
		curcls := ""
		seencls := map[string]int{}
		seenit := map[string]int{}
		_ = seenit
		for _, line := range lines {
			if strings.HasPrefix(line, "Vtable for") {
				// new class
				if curcls != "" {
					code += fmt.Sprintf("}\n")
				}
				curcls = strings.TrimSpace(line[10:])
				if _, ok := seencls[curcls]; ok {
					curcls = ""
					continue
				}
				if strings.Index(curcls, "::") != -1 {
					curcls = ""
					continue
				}
				seencls[curcls] = -1
				seenit = map[string]int{}
				code += fmt.Sprintf("case \"%s\": \n", curcls)
				code += fmt.Sprintf("  switch %s {\n", "mtname")
				continue
			} else if !strings.Contains(line, "::") {
				continue
			} else {
				//log.Println(curcls, line)
				fields := strings.Fields(line)
				// log.Println(len(fields), fields)
				offset, _ := strconv.Atoi(fields[0])
				offset = offset / 8
				pos1 := strings.Index(fields[3], "::")
				pos2 := strings.Index(fields[3], "(")
				item := fields[3][pos1+2 : pos2]
				// log.Println(curcls, item, offset)
				if _, ok := seenit[item]; ok {
					continue
				}
				if strings.HasPrefix(item, "~") {
					continue
				}
				if funk.Contains([]string{"metaObject", "qt_metacast", "qt_metacall", "devType", "metric", "redirected", "initPainter", "sharedPainter", "invalidate"}, item) {
					continue
				}
				if strings.Index(item, "::") != -1 {
					continue
				}
				seenit[item] = offset
				if curcls != "" {
					code += fmt.Sprintf("    case \"%s\": return %d\n", item, offset)
				}
			}
		}
		code += fmt.Sprintf("    } // endof mtname\n")
		if curcls != "" {
			code += fmt.Sprintf("  } // endof clsname\n")
		}
		code += fmt.Sprintf("  return -1\n")
		code += fmt.Sprintf("} // endof func\n")

		filename := fmt.Sprintf("../src/%s/vtable_index.go", strings.ToLower(mod))
		err = ioutil.WriteFile(filename, []byte(code), 0644)
		gopp.ErrPrint(err, filename)
		log.Printf("Wrote %s lines %d\n", filename, len(strings.Split(code, "\n")))
		cmdo = exec.Command("gofmt", "-w", filename)
		err = cmdo.Run()
		gopp.ErrPrint(err, filename)
	}
}

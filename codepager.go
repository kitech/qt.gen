package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"gopp"

	"github.com/emirpasic/gods/lists/arraylist"
)

/*
# 可以多点写入的代码编辑类
# 支持多点写入
# 支持前身写入
# 支持唯一写入
*/

type CodePager struct {
	code          string
	insert_points []string
	lines         map[string][]string
	export_times  int
	newline       string
}

func NewCodePager() *CodePager {
	this := &CodePager{}
	this.insert_points = make([]string, 0)
	this.lines = make(map[string][]string, 0)
	this.newline = "\n"
	return this
}

func (this *CodePager) AddPointer(name string) {
	if _, ok := this.lines[name]; !ok {
		this.lines[name] = make([]string, 0)
		this.insert_points = append(this.insert_points, name)
	}
}

func (this *CodePager) HasPointer(name string) bool {
	_, ok := this.lines[name]
	return ok
}

func (this *CodePager) AllPoints() []string {
	return this.insert_points
}

func (this *CodePager) Append(name, code string) {
	if !this.HasPointer(name) {
		this.AddPointer(name)
	}
	this.lines[name] = append(this.lines[name], code)
}

func (this *CodePager) Appendv(name string, args ...interface{}) {
	code := ""
	for _, argx := range args {
		code += fmt.Sprintf("%v, ", argx)
	}
	this.Append(name, code)
}

func (this *CodePager) Appendf(name, fmts string, args ...interface{}) {
	this.Append(name, fmt.Sprintf(fmts, args...))
}

func (this *CodePager) AppendUnique(name, code string) {
	if !this.HasPointer(name) {
		this.AddPointer(name)
	}

	exist := false
	gopp.Domap(this.insert_points, func(i interface{}) interface{} {
		kv := i.(string)
		if kv == code {
			exist = true
		}
		return nil
	})

	if !exist {
		this.lines[name] = append(this.lines[name], code)
	}
}

func (this *CodePager) Prepend(name, code string) {
	if !this.HasPointer(name) {
		this.AddPointer(name)
	}
	this.lines[name] = append([]string{code}, this.lines[name]...)
}

func (this *CodePager) PrependUnique(name, code string) {
	if !this.HasPointer(name) {
		this.AddPointer(name)
	}

	exist := false
	gopp.Domap(this.insert_points, func(i interface{}) interface{} {
		kv := i.(gopp.Pair)
		if kv.Val.(string) == code {
			exist = true
		}
		return nil
	})

	if !exist {
		this.lines[name] = append([]string{code}, this.lines[name]...)
	}
}

func (this *CodePager) AP(name, code string)  { this.Append(name, code) }
func (this *CodePager) APU(name, code string) { this.AppendUnique(name, code) }
func (this *CodePager) PP(name, code string)  { this.Prepend(name, code) }
func (this *CodePager) PPU(name, code string) { this.PrependUnique(name, code) }

func (this *CodePager) APf(name, format string, args ...interface{}) {
	this.Append(name, fmt.Sprintf(format, args...))
}
func (this *CodePager) APUf(name, format string, args ...interface{}) {
	this.AppendUnique(name, fmt.Sprintf(format, args...))
}
func (this *CodePager) PPf(name, format string, args ...interface{}) {
	this.Prepend(name, fmt.Sprintf(format, args...))
}
func (this *CodePager) PPUf(name, format string, args ...interface{}) {
	this.PrependUnique(name, fmt.Sprintf(format, args...))
}

func (this *CodePager) GetPoint(name string) string {
	return strings.Join(this.lines[name], this.newline)
}

func (this *CodePager) RemovePoint(name string) {
	if this.HasPointer(name) {
		delete(this.lines, name)
	}
}

// TODO
func (this *CodePager) RemoveLine(name, code string) {
	if this.HasPointer(name) {
		for _, line := range this.lines[name] {
			if line == code {
				break
			}
		}
		arr := arraylist.New()
		if false {
			log.Println(arr)
		}
	}

}

// 按照names给出的顺序合并并导出代码。
func (this *CodePager) ExportCode(names []string, comment ...string) string {
	comment_ := "// "
	if len(comment) > 0 {
		comment_ = comment[0]
	}
	blocks := gopp.Domap(names, func(name interface{}) interface{} {
		return fmt.Sprintf("%s %s block begin\n%s\n%s %s block end\n",
			comment_, name.(string),
			strings.Join(this.lines[name.(string)], this.newline),
			comment_, name.(string))
	})
	return strings.Join(gopp.IV2Strings(blocks), this.newline)
}

func (this *CodePager) ExportAll(comment ...string) string {
	return this.ExportCode(this.insert_points, comment...)
}

func (this *CodePager) TotolLength() int {
	return gopp.Doreduce(this.lines, 0, func(v, i interface{}) interface{} {
		kv := i.([]string)
		rv := gopp.Doreduce(kv, 0, func(v, i interface{}) interface{} {
			return v.(int) + len(i.(string))
		})
		return v.(int) + rv.(int)
	}).(int)
}

func (this *CodePager) TotolLine() int {
	return gopp.Doreduce(this.lines, 0, func(v, i interface{}) interface{} {
		kv := i.([]string)
		return v.(int) + len(kv)
	}).(int)
}

func (this *CodePager) Reset() {
	if this.export_times == 0 {
		log.Println("Warning code maybe not export")
	}
	this = NewCodePager()
}

func (this *CodePager) WriteFile(file string) error {
	scc := this.ExportAll()
	return ioutil.WriteFile(file, []byte(scc), 0644)
}

// 模拟写入源代码的文件系统
type CodeFS struct {
	// dir => file => page
	cfs map[string]map[string]*CodePager
	mu  sync.RWMutex
}

func NewCodeFS() *CodeFS {
	this := &CodeFS{}
	this.cfs = map[string]map[string]*CodePager{}
	return this
}

func (this *CodeFS) MkDir(dir string) {
	if _, ok := this.cfs[dir]; !ok {
		this.cfs[dir] = map[string]*CodePager{}
	}
}

func (this *CodeFS) MkFile(dir, file string) {
	this.MkDir(dir)
	if _, ok := this.cfs[dir][file]; !ok {
		this.cfs[dir][file] = NewCodePager()
	}
}

func (this *CodeFS) GetFile(dir, file string) *CodePager {
	this.MkFile(dir, file)
	return this.cfs[dir][file]
}

func (this *CodeFS) ListDirs() (dirs []string) {
	for dir, _ := range this.cfs {
		dirs = append(dirs, dir)
	}
	return
}

func (this *CodeFS) ListFiles(dir string) (files []string) {
	for file, _ := range this.cfs[dir] {
		files = append(files, file)
	}
	return
}

// 写入到磁盘文件系统
// bdir base directory
func (this *CodeFS) WriteToDiskFS(bdir string, ext string) error {
	log.Println("saving dirs:", len(this.cfs))
	fc := 0
	btime := time.Now()

	err := this._WriteToDiskFS(bdir, ext, &fc)

	log.Printf("saving all files done: dir: %d, files: %d, eclapse: %s\n",
		len(this.cfs), fc, time.Now().Sub(btime))
	return err
}
func (this *CodeFS) _WriteToDiskFS(bdir string, ext string, fcp *int) error {
	for dir, _ := range this.cfs {
		for file, cp := range this.cfs[dir] {
			*fcp += 1
			fpath := fmt.Sprintf("%s/%s/%s%s", bdir, dir, file, gopp.IfElseStr(ext == "", "", "."+ext))
			fdir := path.Dir(fpath)
			if !gopp.FileExist(fdir) {
				err := os.MkdirAll(fdir, 0755)
				gopp.ErrPrint(err, fdir)
			}
			err := cp.WriteFile(fpath)
			gopp.ErrPrint(err, fpath, cp.AllPoints(), fdir)
		}
	}
	return nil
}

/// C++ mangling

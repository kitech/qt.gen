
Yet another qt for go/golang binding with FFI/libffi.

speedup compile time and save compile memory usage.

### build

    cd @GOPATH/github.com/therecipe/qt
    git clone github.com/kitech/qt.gen
    cd qt.gen/
    go build

### run 

    QTDIR=$HOME/Qt5.9.1/ ./qt.gen c 2>&1|tee gen.log
    
Sometimes need `ulimit -n 10240`

### supported binding languages
* [x] C
* [x] golang
* [ ] rust
* [ ] ruby
* [ ] nim
* [ ] vlang

### TODOs
* [x] QString arguments as string
* [x] QString record/reference auto destroy
* [x] interface type for arguments passby
* [x] generate skipped class/method/function, but comment it.
* [x] default value process
* [x] 用c封装所有的函数，再用ffi调用
* [x] #define to const
* [x] global variable 全局变量获取

### depends
* therecipe/qt@a76e7081468b0d9d554349b66b4971929f036ce7


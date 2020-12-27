
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

### lang 参数说明
* [x] c 为所有Qt函数生成一个对应的封装函数，做参数返回值的简化，以便纯C语言中调用
* [x] c0 只生成Qt的内联函数/方法的符号生成，不做封装，减小包大小
* [ ] ch 生成Qt所有函数/方法的C原型头文件，不需要C++编译器
* [x] go 为所有Qt函数生成一个对应的Go封装函数
* [ ] gov2 为clipqt(Qt子集)中的函数/方法生成对应的Go封装，减小包大小

C symbol 生成的是所有能够支持的Qt函数，而不是Qt子集

### TODOs
* [x] QString arguments as string
* [x] QString record/reference auto destroy
* [x] interface type for arguments passby
* [x] generate skipped class/method/function, but comment it.
* [x] default value process
* [x] 用c封装所有的函数，再用ffi调用
* [x] #define to const
* [x] global variable 全局变量获取
* [x] 有些类不需要生成代理类

### go-clang TODO
* [x] 无法检查方法delete属性
* [x] 无法检查方法depcreated属性
* [x] 参数default value 的获取
* [ ] ifdef/ifndef块的检测咋用
* [x] sret 检测
* [ ] RECORD参数unpack直传 检测
* [ ] RECORD返回值unpack直传 检测
* [x] class是否抽像
* [ ] 获取comment
* [ ] isFunctionType/isFunctionPointerType/isMemberFunctionPointerType
* [x] isTemplateType
* [x] MSVC always passes 'sret' after 'this', unlike GCC
* [ ] 获取虚方法在vtable中的偏移
   * [ ] libclang实现
   * [x] 运行时实现 , qtrt/mthook3.c
* [ ] 是否有 complete dtor: \_ZN5QRectD2Ev

### depends
* therecipe/qt@a76e7081468b0d9d554349b66b4971929f036ce7
* extended go-clang https://github.com/kitech/go-clang-v3.9


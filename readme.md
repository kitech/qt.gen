
Yet another qt for go/golang binding with FFI/libffi.

speedup compile time and save compile memory usage.

### run 

    QTDIR=$HOME/Qt5.9.1/ ./qt.gen c 2>&1|tee gen.log
    
Sometimes need `ulimit -n 10240`


### TODOs
* [x] QString arguments as string
* [x] QString record/reference auto destroy
* [x] interface type for arguments passby
* [x] generate skipped class/method/function, but comment it.
* [x] default value process
* [x] 用c封装所有的函数，再用ffi调用
* [x] #define to const
* [x] global variable 全局变量获取


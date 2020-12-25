### qt.go v5.15 计划

* [x] 提高执行速度。
  1. 直接使用asmcgocall调用C函数,dl,ffi
  2. C部分直接按照 C++ ABI 调用，不再封装
  3. 参数中go string临时转 C char* 优化
  4. TODO 参数打包转换优化
  5. TODO 减少reflect包的使用
  6. [ ] 重用ffi\_cif，-30ns

* [x] 减小二进制程序
  1. 减小C部分代码量: 只生成Qt的内联函数符号表，不封装，所有函数/方法直接按照 C++ ABI 标准调用
  2. 减小Go部分代码量: 采用手动 clipqt(裁剪qt)的方式，维护qt的一个常用子集
  3. 使用Go 1.15 编译
  4. [ ] 使用Virtual Method Hook + ffi closure在binding语言中实现继承功能，不再需要显式的代理类封装。
  
* [ ] 生成符合C调用标准的Qt函数原型头文件，即显式看到C++对函数原型的改变

* [ ] 对C部分的库分包编译，base/quick/other，有利于程序的打包发布

希望较小的工具程序大小5M左右，中等程序10M左右。（不包括Qt库）

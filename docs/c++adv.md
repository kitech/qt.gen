### C++用法
* 在有const重载的时候，如何强制调用非const的那个

* rvalue this, https://stackoverflow.com/questions/8610571/what-is-rvalue-reference-for-this
  QString().toCaseFolded();

* 现在的inline方法都是 weak symbol, 使用 -O2 编译就优化没了。
  为啥现在测试了一下还有符号呢？
  什么时候会保留 inline 方法的symbol呢？似乎还是有些符号没了
  取inline方法/函数的地址，转换为整数做累加，就能保持符号不被优化掉。但对operator方法有点问题

* 引用类型初始化， T& v = *(T*)voidptrval;
* 右值引用初始化， T&& v = static_cast<T&&>(*(T*)voidptrval);

* 取非静态方法指针 auto x = (void(Class::**)(types)) &Classs::method
* 不能取构造函数指针

* 方法函数指针转 void*: (void*&)memfnptr;

* 关闭忽略返回值时的警告， (void)fn();

### C++ virtual方法hook
* https://www.codeproject.com/articles/1100579/polyhook-the-cplusplus-x-x-hooking-library
* https://github.com/stevemk14ebr/PolyHook_2_0



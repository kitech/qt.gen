
### golang封装生成时的考虑
* 使用结构体封装，
* 生成指针的方法，用于qt实例的GC回收
* Cthis => cthis字段 Cthis 无法简单引用，需要用方法，所以就小写了吧，成为私有字段
* GetCthis() 由于存在多继承，需要对多继承的类实现GetCthis
* SetCthis() 也许可以用 Fromptr替代。
* Fromptr() func (nilval *QObject) Fromptr()
* FromCthis() 功能与Fromptr相同，命名与GetCthis一致
* Addr() => GetAddr()

### vlang封装生成时的考虑
* 使用结构体封装， type QObject { cthis voidptr }
* 生成*非*指针的方法  fn (this QObject) someMethod() {}
* 即使改过的V编译器，函数/方法名也不能以大写字母开头
* qt实例的GC回收考虑使用bdwgc
* GetCthis()
* SetCthis()
* Fromptr()
* Addr()
* [ ] V embed 结构体实现还不够用
* [ ] V interface 实现还差太多

### ch 封装生成时的考虑
* 类名 typedef void* QObject;



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
* getCthis()
* setCthis()
* fromptr()
* freecpp()
* free()
* - getAddr() 如果用非指针方法，这个方法得到的结果是临时的
* - 在封装结构体中保存 cthis 的地址！！！
* [ ] V embed 结构体实现还不够用
* V embed 结构体不支持来自不同结构体的同名字段（作者说）。但是实际测试支持，与go行为一致
* [ ] V interface 实现还差太多
* 为每个类生成转换为所有基类的方法及自身的方法
* 类对应的interface名称为 toQClass(), 在go里是 QClass\_PTR()
  这样就能够在interface实现不完全的时候依旧能用
* 类对应的interface包含的方法集合，getCthis(),
* vlang 无法优雅的处理 typed integer 常量，占用编译的程序大小，占用CPU处理。不要用typed integer类型
* 用sumtype 替代 interface似乎可能，sumtype还可以放在sumtype中

### ch 封装生成时的考虑
* 类名 typedef void* QObject;
* 处理返回16字节Record的情况

### clipqt 组织时的事项
* 去掉 QT开头的宏
* 去掉内联实现函数体，只需要保留原型
* 去掉enum定义
* 去掉private部分
* 去掉deprecated部分原型
* 不需要析构函数原型
* 可以保留inline关键字


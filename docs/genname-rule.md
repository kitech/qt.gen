
### binding api rule for Vlang

include mut/& this, method name reformat or not.

and help method/function for qt class.

* module head for keep camel names, @[translated] module qtxxx
* add pub interface QClassNameITF {}
* add extra method to emulate virtual method,
    ```pub fn (this &QClassName) toQClassName() &QClassName```

* constructor 1, ```pub fn newClassName(...) &QClassName```
    original from Go binding, but V func/method cannot start with uppercase
* constructor 2, ```pub fn QClassName.new(...) &QClassName```
* constructor from C ptr 1, ```pub fn QClassNameFromptr(ptr voidptr) &QClassName```
* constructor from C ptr 2, ```pub fn QClassName.fromptr(ptr voidptr) &QClassName```
* constructor from C ptr 3, ```pub fn (_ &QClassName) newFromptr(ptr voidptr) &QClassName```
    thus not write code like, qtcore.QString.fromptr(ptr),
    just write as, s1.newFromptr(ptr),
    ofcause you need already have an old object.

* destructor 1, ```pub fn deleteQClassName(this &QClassName)```
    this func also used in set_finalizer
* destructor 2, ```pub fn (this &QClassName) dtor()```
    cannot use delete for name conflict in case.
* destructor 3, ```pub fn (this &QClassName) free()```
    for V autofree invoke, later when V fully implment.

* normal metohd, ```pub fn (this &Type) methodName(...)```
* static method:
  1. like normal method, but omit this:
    ```pub fn (_ &QClassName) methodName(...)```
  2. use V's static method:
    ```pub fn QClassName.methodName(...)```

* for inherit, ```pub fn (_ &QClassName) newForInherit_(...) &QClassName```

* \#define const, keep uppercase, trimed QT_, QT, Q_ prefix.
* class anonymus enum, ```QClassNameEnum.EnumName```
* class named enum, ```QClassNameEnumType.EnumName```
* global named enum, ```EnumType.EnumName```


### binding api rule for Vlang
  TODO

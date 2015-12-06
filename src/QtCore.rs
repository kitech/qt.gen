
use std::fmt::Debug;
use std::any::Any;

// 实现类似C++的overload方法
// 要为每个有重载的方法生成一个trait，再为每重载的方法生成不同参数的trait实现。
// 这样trait的个数就是C++类的唯一方法名个数。
// 实现trait的个数就是C++类的所有方法个数。
// 如果在做下优化，没有重载的方法不生成trait了。
// 有一个不好用的地方是，调用的使用都要使用参数列表的tuple形式，如arg((a1, a2))。
// 参数表有可能冲突，但最终实现是否会冲突呢？应该不会，类-方法trait-参数表tuple唯一。
//


pub struct QString {
    pub ival: i32
}

impl QString {
    pub fn arg<T: QString_arg>(&mut self, value: T) -> QString {
        let s = value.arg(self);
        println!("fff {}", s.ival);
        return s
        // return QString{ival: 4}
    }
}

pub trait QString_arg {
    fn arg(self, this:&mut QString) -> QString;
}

impl QString_arg for (QString, QString) {
    fn arg(self, this:&mut QString) -> QString {
        let args = self;
        println!("111,{},{}", "ieiiewr", this.ival);
        return QString{ival:1}
    }
}

impl QString_arg for (QString, QString, QString) {
    fn arg(self, this:&mut QString) -> QString {
        println!("222");
        return QString{ival:2}
    }
}

impl QString_arg for (i32) {
    fn arg(self, this:&mut QString) -> QString {
        println!("333");
        return QString{ival:3}
    }
}

/////
pub struct QByteArray {
    pub ival: i32
}

impl QByteArray {
    pub fn arg<T: QByteArray_arg>(&mut self, value: T) -> QByteArray {
        let s = value.arg(self);
        println!("fff {}", s.ival);
        return s
        // return QByteArray{ival: 4}
    }
}

pub trait QByteArray_arg {
    fn arg(self, this:&mut QByteArray) -> QByteArray;
}

impl QByteArray_arg for (QByteArray, QByteArray) {
    fn arg(self, this:&mut QByteArray) -> QByteArray {
        let args = self;
        println!("111,{},{}", "ieiiewr", this.ival);
        return QByteArray{ival:1}
    }
}

impl QByteArray_arg for (QByteArray, QByteArray, QByteArray) {
    fn arg(self, this:&mut QByteArray) -> QByteArray {
        println!("222");
        return QByteArray{ival:2}
    }
}

impl QByteArray_arg for (i32) {
    fn arg(self, this:&mut QByteArray) -> QByteArray {
        println!("333");
        return QByteArray{ival:3}
    }
}


extern crate libc;

use self::libc::c_void;
use std::ops::Deref;
use std::ops::DerefMut;
use std::fmt::Debug;
use std::any::Any;
use std::collections::HashMap;

// 实现类似C++的overload方法
// 要为每个有重载的方法生成一个trait，再为每重载的方法生成不同参数的trait实现。
// 这样trait的个数就是C++类的唯一方法名个数。
// 实现trait的个数就是C++类的所有方法个数。
// 如果在做下优化，没有重载的方法不生成trait了。
// 有一个不好用的地方是，调用的使用都要使用参数列表的tuple形式，如arg((a1, a2))。
// 参数表有可能冲突，但最终实现是否会冲突呢？应该不会，类-方法trait-参数表tuple唯一。
//


pub struct RString {
    pub ival: i32,
}

impl RString {
    pub fn arg<T: RString_arg>(&mut self, args: T) -> RString {
        let s = args.arg(self);
        println!("fff {}", s.ival);
        return s
        // return RString{ival: 4}
    }
}

pub trait RString_arg {
    fn arg(self, this:&mut RString) -> RString;
}

impl RString_arg for (RString, RString) {
    fn arg(self, this:&mut RString) -> RString {
        let args = self;
        println!("111,{},{}", "ieiiewr", this.ival);
        return RString{ival:1}
    }
}

impl RString_arg for (RString, RString, RString) {
    fn arg(self, this:&mut RString) -> RString {
        println!("222");
        let arg0 = self.0;
        let arg1 = self.1;
        let arg2 = self.2;
        let tmp = arg0.ival;
        return RString{ival:2}
    }
}

impl RString_arg for (i32) {
    fn arg(self, this:&mut RString) -> RString {
        println!("333");
        // let arg0 = self.0;
        let arg0 = self;
        return RString{ival:3}
    }
}

// 返回值重载，使用generic实现。
pub trait RString_read<RetType> {
    fn read(self, this:&mut RString) -> RetType;
}

// read(i32, i64) -> i32
impl RString {
    pub fn read<RetType, T: RString_read<RetType>>(&mut self, args: T) -> RetType {
        let res = args.read(self);
        return res;
    }
}

impl RString_read<i32> for (i32, i64) {
    fn read(self, this:&mut RString) -> i32 {
        return 1;
    }
}

impl RString_read<i64> for (i32, i64) {
    fn read(self, this:&mut RString) -> i64 {
        return 1;
    }

}

// 这个重载功能很强大哦！！！
impl RString_read<u64> for (i32, i64) {
    fn read(self, this:&mut RString) -> u64 {
        return 1;
    }
}

// 这种写法是错误的，
// error: wrong number of type arguments: expected 1, found 0 [E0243]
/*
impl RString_read<> for (i32, i64) {
    fn read(self, this:&mut RString) {
        return;
    }
}
*/

// 很强大，这样也行啊。()叫unit type。
// 是不是rust的void类型呢
/*
impl RString_read<()> for (i32, i64) {
    fn read(self, this:&mut RString) -> () {
        return (); // OK
        return; // OK
    }
}
 */

// 模板参数不能省，其他位置都可以省略。
impl RString_read<()> for (i32, i64) {
    fn read(self, this:&mut RString) {
        return (); // OK
        return; // OK
    }
}

impl RString_read<RString> for (RString) {
    fn read(self, this:&mut RString) -> RString {
        return RString{ival:0};
    }
}

// 继续generic???这就不支持了，不用担心继续generic的问题了
/*
impl RString_read<RetType> for (RString) {
    fn read(self, this:&mut RString) -> RString {
        return RString{ival:0};
    }
}
 */


// test for static method overload
impl RString {
    // static
    pub fn number<RetType, T: RString_number<RetType>>(args: T) -> RetType {
        return args.number();
        // return ();
    }
}

pub trait RString_number<RetType> {
    fn number(self) -> RetType;
}

impl RString_number<RString> for (i32) {
    fn number(self) -> RString {
        return RString{ival: 5};
    }
}

impl RString_number<RString> for (u32) {
    fn number(self) -> RString {
        return RString{ival: 5};
    }
}

impl RString_number<RString> for (i64) {
    fn number(self) -> RString {
        return RString{ival: 5};
    }
}

impl RString_number<RString> for (f64) {
    fn number(self) -> RString {
        return RString{ival: 5};
    }
}

// test for static method and non-static method overload
impl RString {
    // static，
    // 只好这样了，加个_s后缀
    pub fn compare_s<RetType, T:RString_compare_s<RetType>>(args: T) -> RetType {
        return args.compare_s();
    }
    // error: duplicate method [E0201]
    /*
    pub fn compare<RetType, T:RString_number<RetType>>(&mut self, args: T) -> RetType {
    }
     */
    // non-static
    // error: error: duplicate method [E0201]
    // 这种使用不同的trait也还是不行
    /*
    pub fn compare<RetType, T:RString_compare<RetType>>(&mut self, args: T) -> RetType {
        return args.compare(self);
    }
     */

    pub fn compare<RetType, T:RString_compare<RetType>>(&mut self, args: T) -> RetType {
        return args.compare(self);
    }
}

pub trait RString_compare_s<RetType> {
    fn compare_s(self) -> RetType;
}

pub trait RString_compare<RetType> {
    fn compare(self, &mut RString) -> RetType;
}




/////
pub struct RByteArray {
    pub ival: i32
}

impl RByteArray {
    pub fn arg<T: RByteArray_arg>(&mut self, value: T) -> RByteArray {
        let s = value.arg(self);
        println!("fff {}", s.ival);
        return s
        // return RByteArray{ival: 4}
    }
}

pub trait RByteArray_arg {
    fn arg(self, this:&mut RByteArray) -> RByteArray;
}

impl RByteArray_arg for (RByteArray, RByteArray) {
    fn arg(self, this:&mut RByteArray) -> RByteArray {
        let args = self;
        println!("111,{},{}", "ieiiewr", this.ival);
        return RByteArray{ival:1}
    }
}

impl RByteArray_arg for (RByteArray, RByteArray, RByteArray) {
    fn arg(self, this:&mut RByteArray) -> RByteArray {
        println!("222");
        return RByteArray{ival:2}
    }
}

impl RByteArray_arg for (i32) {
    fn arg(self, this:&mut RByteArray) -> RByteArray {
        println!("333");
        return RByteArray{ival:3}
    }
}

// 继承

pub struct NWidget {
    pub qclsinst: *mut c_void,
}

impl NWidget {
    pub fn hehe(&self) {
        println!("hehehehehhe");
    }
}

pub struct NButton {
    pub vbase: NWidget,
    // pub rbase: &'a NWidget,
    // pub mbase: &'a mut NWidget,
    pub qclsinst: *mut c_void,
}

impl NButton {
    pub fn new() -> NButton {
        let mut i = 1 as *mut c_void;
        let mut b = NWidget{qclsinst: i};
        return NButton{vbase: b, qclsinst: i};
    }
    pub fn haha(&self) {
        println!("hahahahhaha");
    }
}

impl AsRef<NWidget> for NButton {
    fn as_ref(&self) -> &NWidget{
        &self.vbase
    }
}

impl Deref for NButton {
    type Target = NWidget;

    fn deref(&self) -> &NWidget {
        &self.vbase
    }
}


pub fn test_inherint_1() {
    let btn = NButton::new();
    btn.haha();
    btn.hehe();
    btn.as_ref().hehe();
}

// generic模板
struct NList<T> {
    ival: T,
}

type NList1<T> = Vec<T>;
type NMap<K,V> = HashMap<K, V>;


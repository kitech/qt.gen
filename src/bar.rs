
#![feature(libc)]
extern crate libc;
use self::libc::c_void;
use self::libc::c_char;
use self::libc::int8_t;
use self::libc::strlen;

use super::foo::Foo;

pub struct Bar {
    pub qclsinst: *mut c_void,
}

fn test_refer_foo_member_var(a0: &mut Foo) {
    println!("fff %v");
    a0.qclsinst;
}


trait Bar_trait_test_lifetime {
    fn test1(self) -> i32;
}

/*
这种写法能重现错误提示
// error: missing lifetime specifier [E0106]
impl Bar_trait_test_lifetime for (&mut Foo) {

}
*/

/*
这种写法解决了E0106的错误提示
*/
impl<'a> Bar_trait_test_lifetime for (&'a mut Foo) {
    fn test1(self) -> i32 {
        return 1;
    }
}

// char *类型转换，OK
impl<'a> Bar_trait_test_lifetime for (&'a mut str) {
    fn test1(self) -> i32 {
        self.as_ptr() as *const c_void;
        self.as_ptr() as *const c_char;
        return 1;
    }
}

// bool & 类型转换，
fn test_boolstart<'a>(a0: &'a mut bool) {
    *a0 = true;
    {
        let mut swap: int8_t = 0;
        if swap == 1 {*a0 = true;}
    }
    *a0 = {let mut bv: int8_t = 0; if bv == 1 {true} else {false}}
}

fn test_i8start<'a>(a0: &'a mut i8) {
    // a0 as *mut c_void;  // error
    a0 as *mut int8_t;  // ok
}


// test trait for (&'a mut i32, &'a mut str, i32)
pub struct TestBar;
impl TestBar {
    pub fn newbar<T: BarTrait>(value: T) {
        value.newbar();
    }
}
pub trait BarTrait {
    fn newbar(self);
}

impl<'a> BarTrait for (&'a mut i32, &'a mut String, i32) {
    fn newbar(self) {
        let arg0 = self.0;
        println!("{}", arg0);
    }
}

// test Vec<String> for char **
pub fn test_vec_str(a0: Vec<String>) {

}

// test return multable reference
// error: error: `n` does not live long enough
// pub fn test_ret_mutable_ref<'a>() -> &'a mut i32 {
//     let n = 5;
//     let r = &mut n;
//     return r;
// }

// pub fn test_ret_mutable_ref<'a>() -> &'a Foo {
//     let n = Foo{qclsinst: 0 as *mut c_void};
//     let r = &n;
//     return r;
// }

// 看来只好这样了。
pub fn test_ret_mutable_ref() -> Foo {
    let n = Foo{qclsinst: 0 as *mut c_void};
    let r = n;
    return r;
}

// test *u8 to String
// error: non-scalar cast: `*mut u8` as `&mut collections::string::String`
pub fn test_star_u8_to_String() {
    let s1 = String::from("abcdefg");
    let s1u8 = s1.as_ptr();
    let s1u8mut = s1u8 as *mut u8;
    let s1u8mut2 = s1u8mut as *mut u8;
    let s2 = s1u8 as *mut c_char;
    unsafe {strlen(s2)};
    // unsafe {strlen(s1u8)};
    unsafe {String::from_raw_parts(s1u8mut2, 7, 9)};
}


// TODO
// QString => *mut c_uchar
// let arg1 = self.1.as_ptr()  as *mut c_uchar;
// Vec<?> => *mut *mut ?
// let arg1 = self.1.as_mut_ptr()  as *mut *mut c_float;
// return void * 
// return const float * 

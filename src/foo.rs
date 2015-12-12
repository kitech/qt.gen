
#![feature(libc)]
extern crate libc;
use self::libc::c_void;

pub struct Foo {
    pub qclsinst: *mut c_void,
}



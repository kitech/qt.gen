
#![feature(libc)]
#![feature(core)]
#![feature(collections)]

extern crate libc;
use self::libc::int32_t;
use self::libc::uint32_t;
use std::str;
use std::ffi::CStr;

#[link(name = "Qt5Core")]
extern {
    fn _Z5qrandv() -> int32_t;
    fn _Z6qsrandj(a: uint32_t) ;
    fn qVersion() -> *const libc::c_char;
}

pub fn qVersion_<'a>() -> &'a str {
    let ver = unsafe {qVersion()};
    let verstr: &CStr = unsafe {CStr::from_ptr(ver)};
    let verbuf: &[u8] = verstr.to_bytes();
    let vers: &str = unsafe {str::from_utf8_unchecked(verbuf)};
    return vers
} 

pub fn qrand() -> i32 {
    let x = unsafe {
        qVersion();
        _Z6qsrandj(568);
        _Z5qrandv()
    };
    let ver = unsafe {qVersion()};
    let verstr: &CStr = unsafe {CStr::from_ptr(ver)};
    let verbuf: &[u8] = verstr.to_bytes();
    let vers: &str = unsafe {str::from_utf8_unchecked(verbuf)};
    println!("verrrr {}", vers);
    return x;
}



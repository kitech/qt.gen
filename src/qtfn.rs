
#![feature(libc)]
#![feature(core)]
#![feature(collections)]

extern crate libc;
use self::libc::int32_t;
use self::libc::uint32_t;
use self::libc::c_void;
use self::libc::c_char;
use std::str;
use std::ffi;
use std::ffi::CStr;
use std::ffi::CString;
// use os;
use std::env;
use std::vec;


#[link(name = "Qt5Core")]
#[link(name = "Qt5Gui")]
#[link(name = "Qt5Widgets")]
extern {
    fn _Z5qrandv() -> int32_t;
    fn _Z6qsrandj(a: uint32_t) ;
    fn qVersion() -> *const libc::c_char;
    fn _ZN12QApplicationC1ERiPPci(this: *mut libc::c_void, argc: int32_t, argv: int32_t, a: int32_t);
    // fn _ZN12QApplicationC2ERiPPci(this: *mut libc::c_void, argc: int32_t, argv: int32_t, a: int32_t);
    fn _ZN12QApplicationC2ERiPPci(this: *mut libc::c_void, argc: *mut int32_t, argv: *mut *mut libc::c_char, a: int32_t);
    fn _ZN12QApplication4execEv(this: *mut libc::c_void) -> libc::int32_t;
}

pub fn NewClass() {
    let this: *mut libc::c_void = unsafe{libc::calloc(1, 200)};
    let argv: *const *const libc::c_char;
    let argv0: *mut libc::c_void = unsafe{libc::calloc(1, 20)};
    // unsafe{libc::memset(argv0, 0, 20)};
    // argv = argv0 as * const *const libc::c_char;
    argv = 0  as * const *const libc::c_char;

    env::args_os();
    

    // ok??? https://codeseekah.com/2015/01/25/rusts-osargs-to-cs-argv/
    // let argv2: Vec<ffi::CString> = std::vec::Vec.new(); // = env::args_os().into_iter().map(|arg| { ffi::CString::from_vec_unchecked(arg) } ).collect();
    // let args2:Vec<*const c_char> = argv2.into_iter().map(|arg| { arg.as_ptr() } ).collect();
    let mut argv2: vec::Vec<ffi::CString> = vec::Vec::new();
    argv2.push(ffi::CString::new("./abc").unwrap());
    let args2:Vec<*const c_char> = argv2.into_iter().map(|arg| { arg.as_ptr() } ).collect();

    println!("ret={}", "iop111");
    // unsafe {_ZN12QApplicationC1ERiPPci(this, 0, 0, 0)};
    unsafe {
        let mut argc: int32_t = 1;
        // _ZN12QApplicationC2ERiPPci(this, &mut argc as *mut int32_t, args2.as_ptr() as *mut *mut c_char, 0)}; // OK
        _ZN12QApplicationC2ERiPPci(this, &mut 0 as *mut int32_t, 0 as *mut *mut c_char, 0)}; // OK
    // println!("this={}", this);
    println!("ret={}", "iop");
    let ret = unsafe {_ZN12QApplication4execEv(this)};
    println!("ret={}", ret);
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



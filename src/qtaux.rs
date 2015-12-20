extern crate libc;
use self::libc::c_void;

pub type c_char16 = i16;
pub type c_uchar16 = u16;
pub type c_char32 = i32;
pub type c_uchar32 = u32;
pub type c_voidp = *mut c_void;
pub type c_cvoidp = *const c_void;

// static c_null: *mut c_void = 0 as *mut c_void;
// static c_cnull: *const c_void = 0 as *const c_void;

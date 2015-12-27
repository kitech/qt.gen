extern crate libc;
use self::libc::c_void;

pub type c_pointer = u64;
pub type c_char16 = i16;
pub type c_uchar16 = u16;
pub type c_char32 = i32;
pub type c_uchar32 = u32;
pub type c_enum = i32;
pub type c_voidp = *mut c_void;
pub type c_cvoidp = *const c_void;
pub type c_funcp = *mut c_void;

// static c_null: *mut c_void = 0 as *mut c_void;
// static c_cnull: *const c_void = 0 as *const c_void;

pub type QVector<T> = Vec<T>;
pub type QList<T> = Vec<T>;

//
pub trait AsCPtr {
    fn as_mut_ptr(self) -> *mut c_void;
    fn as_const_ptr(self) -> *const c_void;
}
impl AsCPtr for c_pointer {
    fn as_mut_ptr(self) -> *mut c_void {self as *mut c_void}
    fn as_const_ptr(self) -> *const c_void {self as *const c_void}
}

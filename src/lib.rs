
// #[warn(non_snake_case)]
// #[warn(non_camel_case_types)]
// #[warn(unused_mut)]
// #[warn(unused_attributes)]
// #[warn(unused_imports)]
#![allow(non_snake_case)]
#![allow(non_camel_case_types)]
#![allow(unused_mut)]
#![allow(unused_attributes)]
#![allow(unused_imports)]
// #![cfg_attr(lte_rustc_1_5, allow(raw_pointer_derive))]
// #![allow(custom_derive)]

#[link(name = "Qt5Core")]
#[link(name = "Qt5Gui")]
#[link(name = "Qt5Widgets")]
#[link(name = "QtInline")]
extern {}  // 这行还是需要的


pub mod QtCore;
pub mod qtfn;

pub mod foo;
pub mod bar;

mod qtaux;
pub mod core;
pub mod gui;
pub mod widgets;


// #[test]
// fn it_works() {
// }

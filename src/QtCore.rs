
use std::fmt::Debug;
use std::any::Any;

pub fn qrand() -> String {
    "123".to_string()
}

pub struct QString {
    pub x: i64,
}


trait QString_all_method {

}

trait QString_method_append_0 {
    fn append(&self) -> i32;
}
trait QString_method_append_1 {
    fn append<T0: Any>(&self, a0: &T0) -> i32;
}
// trait QString_method_append_2 {
//     fn append(a0: Any, a1: Any);
// }
// trait QString_method_append_3 {
//     fn append(&self, a0: Any, a1: Any, a2: Any) -> i32;
// }

pub trait QString_append<T> {
    fn append(&self);
}


impl<'a> QString_append<QString> for (&'a Any) {
    fn append(&self) {
        // let (host, port) = *self;
        let pa0 = *self;
        123;
    }
}

// impl QString_method_append_0 for QString {
//     fn append(&self) -> i32 {
//         println!("0 args");
//         123
//     }
// }
// impl QString_method_append_1 for QString {
//     fn append<T0: Any>(&self, a0: &T0) -> i32 {
//         println!("1 args");
//         let args = (a0);
//         // QString_method_append(self, args);
//         123
//     }
// }


pub struct QDateTime {
    pub x: i64,
    // o: *mut u64,
}

impl QDateTime {
    pub fn toTime_t(&mut self) -> u32 {
        // self.x = 123 as i64;
        120
    }
}


extern crate rustqt;

use rustqt::*;
use rustqt::QtCore::*;

fn main() {
    println!("Hello rust!!!");
    // 我都use了，为什么还要加个QtCore前缀呢
    let mut dt = QtCore::QDateTime{x:5};
    let qs = QtCore::QString{x:5};
    println!("qrand ret:{}, {}, {}, {}",
             QtCore::qrand(), 123, 456, dt.toTime_t());

    println!("qs ret:{:?}", qs.append(&0));
}

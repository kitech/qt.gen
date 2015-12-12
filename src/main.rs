
extern crate rustqt;

use rustqt::*;
// use rustqt::QtCore::*;
use rustqt::qtfn::*;
use rustqt::core::*;
use rustqt::bar::*;

fn test_QCoreApplication() {
    // QCoreApplication::NewQCoreApplication2(); // ok
    let mut a0 = 0;
    let mut a1  = String::from("aaaaa");
    let a2 = 0;
    let mut app = QCoreApplication::NewQCoreApplication((&mut a0, &mut a1, a2));
    (&mut app).exec(());
    // TestBar::newbar((&mut a0, &mut a1, a2));
}

fn main() {
    println!("Hello rustqt!!!");
    /*
    // 我都use了，为什么还要加个QtCore前缀呢
    let mut x = QString{ival:999};
    let s1 = QString{ival:111};
    let s2 = QString{ival:222};
    let s3 = QString{ival:333};
    let s4 = QString{ival:444};
    let s5 = QString{ival:555};
    let i1 = 111;

    let sr = x.arg((s1, s2));
    let sr2 = x.arg((i1));
    let sr3 = x.arg((s5, s4, s3));
    println!("sr={}, sr2={}, sr3={}", sr.ival, sr2.ival, sr3.ival);

    let rd1 = qrand();
    let ver: &str = qVersion_();
    println!("rd1={}, ver={}", rd1, ver);
    */
    // NewClass();

    test_QCoreApplication();
}

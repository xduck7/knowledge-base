mod stack_heap;
mod ownership_movement;
mod calculate_length;
mod mutable;
mod slice;
mod structs;
mod impl_example;

fn main() {
    stack_heap::run();

    ///////////////////////////////////////////////////////////////////

    // OWNERSHIP RULES
    // 1. Each value in Rust has a variable that's called its OWNER.
    // 2. There can only be one owner at a time.
    // 3. When the owner goes out of scope, the value will be dropped.

    { // s is not valid here. It's not yet declared
        let s = String::from("string"); // s is valid from this point
        // do smth
    } // scope is now over, and s is not valid

    ///////////////////////////////////////////////////////////////////

    let a = 5;
    let b = a; // Copied, because basic types always copy
    println!("main: a = {}, b = {}", a, b);

    let s1 = String::from("string");

    // let s2 = s1;
    // println!("{}, world!", s1); // ERROR: value was moved

    let s2 = s1.clone(); // OK
    println!("main: {}, world!", s2);

    ///////////////////////////////////////////////////////////////////

    ownership_movement::run();

    ///////////////////////////////////////////////////////////////////

    // SOME INFO FROM BOOK

    let st = String::from("hello");

    //  // STACK
    // st: String = st {               // HEAP
    //                 ptr  -----> |  index | value |
    //                 len         |    h | 0       |
    //                 cap         |    e | 1       |
    //               }             |    l | 2       |
    //                             |    l | 3       |
    //                             |    o | 4       |

    ///////////////////////////////////////////////////////////////////

    mutable::run();

    ///////////////////////////////////////////////////////////////////

    slice::run();

    //////////////////////////////////////////////////////////////////

    structs::run();

    //////////////////////////////////////////////////////////////////

    impl_example::run();
}
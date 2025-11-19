pub fn run() {
    // COPY
    let aa = 5;
    make_copy(aa);
    println!("ownership_movement: {}", aa);

    // GIVE OWNERSHIP
    let bb = String::from("hello");
    let cc = give_ownership();
    println!("ownership_movement: bb: {}, cc: {}", bb, cc);

    // TAKE AND GIVE BACK
    let ff = String::from("hello");
    let gg = takes_and_gives_back(ff);
    println!("ownership_movement: gg: {}", gg);
    // println!("ownership_movement: ff: {}", ff); ERROR: value was moved


    // TAKE OWNERSHIP
    let st = String::from("hello");
    takes_ownership(st);
    // println!("{}", st) // ERROR: value was moved
}                //  |
                 // | because the ownership of "hello" was received by some_string from st
                 // | after executing fn takes_ownership some_string dropped
                // \/
fn takes_ownership(some_string: String) {
    println!("ownership_movement: {}", some_string);
}

fn give_ownership() -> String {
    let some_string = String::from("hello");
    some_string
}

fn takes_and_gives_back(a_string: String) -> String {
    a_string
}

fn make_copy(some_integer: i32) {
    println!("ownership_movement: {}", some_integer);
}

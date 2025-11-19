pub fn run() {
    let mut s1 = String::from("hello");
    let s2 = String::from("world");
    let s3 = sum_strings(&mut s1, &s2);
    println!("mutable: s3 = {}", s3);

    ////////////////////////////////////////////

    let mut s = String::from("hello");
    let r1 = &mut s;     // LIMIT
    r1.push_str("added");
    // let r2 = &mut s; // second mutable borrow occurs here
    println!("mutable: {}", r1);

    ////////////////////////////////////////////

    let reff = correct_return_ref();
}

fn sum_strings(s1: &String, s2: &String) -> String {
    let mut s = String::from("");
    s.push_str(s1);
    s.push_str(" ");
    s.push_str(s2);
    s.push_str("!");
    s
}

// fn incorrect_return_ref() -> &String { // cant return ref to empty obj.
//     let s = String::from("hello");
//     &s
// }

fn correct_return_ref() -> String {
    let s = String::from("hello");
    s
}
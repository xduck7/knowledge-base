use unicode_segmentation::UnicodeSegmentation;

pub fn run() {
    // Strings stored as UTF-8 encoded bytes

    let s1 = "hello world".to_string();
    let s2 = "hello world";
    let s3 = String::from("hello world");
    let s4 = String::new();

    let mut sample_string = String::from("hello ");

    sample_string.push_str("world");
    sample_string.push('!');

    println!("The string is {}", sample_string);

    let a = String::from("hello ");
    let b = String::from("world");
    let c = format!("{}{}", a,b);
    println!("{}", c);

    ////////////// Iteration

    let iter_string = String::from("привет world");

    for b in iter_string.bytes() {
        print!("{} ", b);
    }

    println!();

    for c in iter_string.chars() {
        print!("{} ", c);
    }

    println!();

    for g in iter_string.graphemes(true) {
        print!("{} ", g);
    }

    println!();
}
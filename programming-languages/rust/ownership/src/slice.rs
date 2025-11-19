pub fn run() {
    let txt = String::from("abcdefghijk");
    let slice = &txt[0..3]; // abc
    let slice2 = &txt[3..6]; // def

    let slice3 = &txt[6..9];
    println!("slice: {}", slice3);
    let slice3 = "123";
    println!("slice: {}", slice3);

    println!("slice: {}", slice);
    println!("slice: {}", slice2);
    println!("slice: {}", txt);
}
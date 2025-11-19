pub fn run() {
    let s1 = String::from("hello");
    let (s2, len) = strange_calculate_length(s1);
    println!("calculate_length: The length of '{}' is {}.", s2, len);

    let s3 = String::from("hello");
    let len = normal_calculate_length(&s3);
    println!("calculate_length: The length of '{}' is {}.", s3, len);
}

fn strange_calculate_length(s: String) -> (String, usize) {
    let len = s.len();
    (s, len)
}

fn normal_calculate_length(s: &String) -> usize {
    s.len()
}
pub fn run() {
    my_func();
    recursion_sum(5, 10);
}

pub fn my_func() {
    println!("functions: Hello from my_func!");
}

pub fn recursion_sum(a: i32,b: i32) {
    if b == 0 {
        println!("functions: {}", a);
    } else {
        recursion_sum(a+1, b - 1);
    }
}
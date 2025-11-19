struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
}

pub fn run() {
    let rect = Rectangle { width: 30, height: 50 };
    println!("impl_example: The area of the rectangle is {} square pixels.", rect.area());
}
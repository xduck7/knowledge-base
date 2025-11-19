fn main() {

    // Heap: [0i32, 1i32, 2i32, 3i32, 4i32]  // 20 bytes
    // Stack: Vec { ptr, len, cap }            // 24 bytes

    let v = vec![0, 1, 2, 3, 4];
    println!("Vector is {} bytes", size_of_val(&v));

    let mut free_pointer = &v[0] as *const i32;
    let size_of_val = size_of_val(&v[0]);

    println!("v[0] is {size_of_val} bytes");
    for _ in 0..v.len() {
        unsafe {
            println!("val is {}", *free_pointer);
            free_pointer = free_pointer.add(1);
        }
    }
}
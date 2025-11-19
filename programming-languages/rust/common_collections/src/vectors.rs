pub fn run() {
    let a = [1,2,3];
    let mut v: Vec<i32> = Vec::new();
    v.push(1);
    v.push(2);
    v.push(3);
    v.push(4);
    v.push(5);

    let v_by_macros = vec![1,2,3,4,5];

    ////// v is equal v_by_macros

    // let second_in_array = a[10]; // ERROR: out of bound error right away
    let second_in_array = a[0];
    println!("second in array: {}", second_in_array);

    // Let second_in_vector = &v[10]; ERROR: out of bound error. In exec time
    let second_in_vector = v[1];
    println!("second in vector: {}", second_in_vector);

    //////// Error wrapping
    let elem_vector = v.get(20).unwrap_or(&0);
    println!("elem_vector: {}", elem_vector);


    ////// Iteration
    for val in &v {
        print!("{} ", val);
    }

    println!();

    for val in &mut v {
        *val += 50;
    }

    for val in &v {
        print!("{} ", val);
    }

    println!();
}

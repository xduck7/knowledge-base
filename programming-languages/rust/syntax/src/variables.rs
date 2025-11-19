pub fn run() {
    const MY_CONST: i32 = 10;
    println!("variables: MY_CONST: {}", MY_CONST);

    ////////////////////////////////////////////////////////////////////////

    let x = 5;
    println!("variables: x: {}", x);

    // x = 6; x is not mut

    ////////////////////////////////////////////////////////////////////////

    let mut y = 10;
    println!("variables: y before: {}", y);
    y = 101; // y is mut
    println!("variables: y after: {}", y);

    ////////////////////////////////////////////////////////////////////////

    // shadowing
    let z = 10;
    println!("variables: z before: {}", z);
    let z = 20;
    println!("variables: z after: {}", z);
}
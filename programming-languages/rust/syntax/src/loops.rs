pub fn run() {
    // LOOP
    let mut x = 10;
    'lp: loop {
        println!("loops: {}", x);
        x+=1;
        if x == 15 {
            break 'lp;
        }
    }

    // VAR BY LOOP
    let mut some_int = 11;
    let var_by_loop = loop {
        some_int+=1;
        if some_int %2 == 0 {
            break some_int;
        }
    };
    println!("loops: {}", var_by_loop);

    // WHILE
    let mut x = 7;
    while x < 10 {
        println!("loops: WHILE: {}", x);
        x+=1;
    }

    // ITERATOR
    let m = [0,1,2,3,4,5,6,7,8,9];
    for i in m.iter() {
        println!("loops: ITERATOR: {}", i);
    }

    for i in 0..10 { // [0..10)
        println!("loops: ITERATOR RANGE: {}", i);
    }
}
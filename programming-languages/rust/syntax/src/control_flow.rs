pub fn run() {
    let x = 10;
    if x == 10 {
        println!("control_flow: x is 10");
    } else if x == 11 {
        println!("control_flow: x is not 10, x is 11");
    }  else {
        println!("control_flow: x is not 10, x is not 11");
    }

    let val = if x == 10 { 5 } else { 21 };
    
    println!("control_flow: val is {}", val);
}
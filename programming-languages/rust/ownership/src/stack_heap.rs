pub fn run() {
    fn a() {
        let x = "string"; // STACK
        let y = 123; // STACK
        b();
    }

    fn b() {
        let s = String::from("string"); // HEAP
    }
}

//   STACK                HEAP
//  __________         (^^^^^^^^^^)
// [         ]         (          )
// [ b() s---]-------> ( "string" )
// [ a() x y ]         (          )
// [_________]         (__________)

// [func()a b c] = STACK FRAME. Stack frame size determines at compile time
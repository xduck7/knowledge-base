pub fn run() {
    // INTEGER
    let a: i8 = 1;
    let b: i128 = 1; // decimal
    let ab3: i8 = 0b11; // binary
    let ab2: i8 = 0o11; // octal
    let ab: i8 = 0x11; // hex

    // UNSIGN INTEGER
    let c: u8 = 1;
    let d: u128 = 1;
    let cd: u8 = b'A'; // byte

    // FLOAT
    let e: f32 = 1.0;
    let f: f64 = 1.0;

    // BOOL
    let g: bool = true;
    let h: bool = false;

    // CHAR
    let i: char = 'A';
    let j: char = ' '; // UNICODE

    // COMPOUND TYPE

    // TUPLE
    let tup = ("some text", 1, 1.0, true);
    let (aa,bb,cc,dd) = tup;
    println!("data_types: {} {} {} {}", aa, bb, cc, dd);
    let sub_tup = tup.0; // "some_text"

    // ARRAY
    let http_codes = [200,404,403,500];
    let ok = http_codes[0];
    //let index_out_of_bounds = http_codes[10];
    let byte = [0;8]; // [0,0,0,0,0,0,0,0]
}
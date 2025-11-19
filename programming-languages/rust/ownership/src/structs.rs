struct User {
    pub username: String,
    pub email: String,
    sign_in_count: u64,
    active: bool,
}

fn print_user(user: &User) {
    println!("structs: User: username {}, email {}", user.username, user.email);
}

fn build_user(email: String, username: String) -> User {
    User {
        email,
        username,
        active: true,
        sign_in_count: 1,
    }
}

pub fn run() {
    let user_1 = User {
        email: String::from("email_1"),
        username: String::from("user_1"),
        active: true,
        sign_in_count: 1,
    };

    print_user(&user_1);

    let builded_user = build_user(String::from("builded_email"), String::from("buuilded_name"));
    print_user(&builded_user);

    let user_3 = User {
        username: String::from("user_3"),
        ..builded_user
    };

    print_user(&user_3);

    let p_1 = Point(1, 2, 3);
    let c_1 = Color(1, 2, 3);

    print_coord((p_1.0, p_1.1, p_1.2));
    println!("hex: {}", rgb_to_hex(c_1));
    // println!("hex: {}", rgb_to_hex(p_1)); ERROR

    println!("{:#?}", Point(1, 2, 3));
    // display
    // Point(
    //     1,
    //     2,
    //     3,
    // )
}

#[derive(Debug)]
struct Point(i32, i32, i32);
struct Color(i32, i32, i32);

fn print_coord(coords: (i32,i32,i32)) {
    println!("coords: ({}, {}, {})", coords.0, coords.1, coords.2);
}

fn rgb_to_hex(c: Color) -> String {
    format!("#{:02x}{:02x}{:02x}", c.0, c.1, c.2)
}
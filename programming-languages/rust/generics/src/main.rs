fn main() {

    // FUNCS
    let v: Vec<i32> = vec![1, 2, 3, 4, 5,453453,0,1,123];
    println!("{}", get_largest(&v));

    let w: Vec<i64> = vec![10, 20, 30, 40, 50];
    // println!("w[0] = {}", get_largest(&w)); // ERROR: wrong type
    println!("{}", genetic_get_largest(v));
    println!("{}", genetic_get_largest(w));

    // STRUCTS
    let f_point = Point{
        x: 5.75, y: 5.123, color: "Red"
    };

    let i_point = Point{
        x: 5, y: 5, color: "Red"
    };

    // let bad_point = Point{
    //     x: 5.123, y: 5, color: "Red" // x and y must be same data type
    // };
}

// FUNCS

fn get_largest(list: &Vec<i32>) -> i32 {
    let mut largest = list[0];
    for item in list {
        if item > &largest {
            largest = *item;
        }
    }
    largest
}

fn genetic_get_largest<T: PartialOrd + Copy>(list: Vec<T>) -> T {
    let mut largest = list[0];
    for item in list {
        if item > largest {
            largest = item;
        }
    }
    largest
}

// STRUCTS

struct Point<T, U> {
    x: T,
    y: T,
    color: U,
}
///////////// Default enums
#[derive(Debug)]
enum Shape {
    Circle,
    Square,
    Rectangle,
    Triangle,
}

enum Color {
    Red(i8),
    Green(i8),
    Blue(i8),
}

impl Color {
    fn to_string(&self) -> String {
        match self {
            Color::Red(c) => String::from("Red"),
            Color::Green(c) => String::from("Green"),
            Color::Blue(c) =>  String::from("Blue"),
        }
    }
}

struct Figure {
    shape: Shape,
    color: Color,
}

impl Figure {
    fn new(shape: Shape, color: Color) -> Figure {
        Figure {
            shape,
            color,
        }
    }

    fn display(&self) {
        println!("{} {:?}", self.color.to_string(), self.shape)
    }
}

///////////////////////////////////

enum Score {
    Value(i8),
}

impl Score {
    fn to_mark(&self) -> i8 {
        let value = match self {
            Score::Value(v) => *v,
        };

        match value {
            85..=100 => 5,
            65..=84 => 4,
            51..=64 => 3,
            _ => 2,
        }
    }
}

////////////////////////////////// Option enum
// enum Option<T> { // can be null or some value
//     Some(T),
//     None,
// }

///////////////////////////////////////////////


fn main() {
    let circle = Figure::new(Shape::Circle, Color::Red(2));
    circle.display();

    let math_score = Score::Value(82);
    println!("{}", math_score.to_mark());


    // Option enum
    let some_num = Some(4);
    let some_string = Some("a string");
    let absent_num: Option<i32> = None;

    let normal_a = 5;
    let mut some_b: Option<i32> = Some(3);
    // let res_c = normal_a + some_b; // ERROR: cant sum normal and option
    let mut res_c = normal_a + some_b.unwrap_or(123); // if some_b = None then some_b = 123
    println!("{} + {:?} = {}", normal_a, some_b, res_c);

    res_c = 0;
    some_b = None;
    res_c += normal_a + some_b.unwrap_or(123);
    println!("{} + {:?} = {}", normal_a, some_b, res_c);
}


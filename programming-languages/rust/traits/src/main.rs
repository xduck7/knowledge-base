struct Human {
    name: String,
    age: i32,
}

impl Communication for Human {
    fn sound(&self) -> String {
        format!("Hello! My name is {}. I'm {} y.o", self.name, self.age)
    }
}

struct Animal {
    name: String,
    age: i32,
    danger: bool,
}

impl Communication for Animal {
    fn sound(&self) -> String {
        "shhhhh".to_string()
    }
}

pub trait Communication {
    fn sound(&self) -> String {
        String::from("....")
    }
}

struct Car {
    label: String,
    weight: i64,
}

impl Communication for Car {
}


fn main() {
    let me: Human = Human{
        name: "John".to_string(),
        age: 21,
    };
    println!("{}", me.sound()); // Hello! My name is John. I'm 21 y.o

    let cat: Animal = Animal{
        name: "Rizhik".to_string(),
        age: 5,
        danger: false,
    };
    println!("{}", cat.sound()); // shhhhh

    let car: Car = Car{label: "Rizhik".to_string(), weight: 7, };
    println!("{}", car.sound()); // ....
}

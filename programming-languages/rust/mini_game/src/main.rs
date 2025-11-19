use std::io;
use std::io::Write;
use rand::Rng;

const FROM: u32 = 1;
const TO: u32 = 6;

fn main() {
    println!("Enter your balance: ");
    let mut balance = String::new();
    io::stdin().read_line(&mut balance).unwrap();
    let mut balance: u32 = balance.trim().parse().unwrap_or_else(|_| 100u32);
    println!("Your start balance is {}", balance);

    let mut start_position = rand::rng().random_range(FROM..TO);
    let bullet_position = rand::rng().random_range(FROM..TO);

    'main: loop {
        println!("Push!");
        io::stdin().read_line(&mut String::new()).unwrap();

        if start_position.eq(&bullet_position) {
            println!("You lose!");
            balance = 0;
            break 'main;
        }

        balance = balance*2;
        println!("Your balance is {}", balance);

        rotate_drum(&mut start_position);

        if wanna_continue() {
            continue 'main;
        } else {
            break 'main;
        }
    }


    println!("Bye! Your final balance is {}", balance);
}

fn wanna_continue() -> bool {
    loop {
        print!("Do you want continue? (y/n): ");
        io::stdout().flush().expect("Failed TO flush stdout");

        let mut buffer = String::new();
        io::stdin().read_line(&mut buffer).expect("Failed TO read line");

        match buffer.trim() {
            "y" | "Y" => return true,
            "n" | "N" => return false,
            _ => println!("Please enter 'y' or 'n'"),
        }
    }
}

fn rotate_drum(pos: &mut u32) {
    *pos = (*pos + 1) % TO;
    if *pos == 0 {
        *pos = TO;
    }
}
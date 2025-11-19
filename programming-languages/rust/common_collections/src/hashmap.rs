use std::collections::HashMap;

#[derive(Eq, PartialEq, Hash, Debug)]
enum Team {
    Blue,
    Red,
    Green,
}

pub fn run() {
    let mut score = HashMap::new();
    score.insert(Team::Blue, 10);
    score.insert(Team::Red, 10);

    for (team, score) in &score {
        println!("Team: {:?} score: {}", team, score);
    }

    let blue_score = *score.entry(Team::Blue).or_insert(0);
    let red_score = *score.entry(Team::Red).or_insert(0);
    let green_score = *score.entry(Team::Green).or_insert(0);
    let total = blue_score + red_score + green_score;

    println!("Total score: {}", total);
}
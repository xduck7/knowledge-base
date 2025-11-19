/// `cargo bench`
/// `cargo test`

///////////////////////////////

pub fn recursive_sum(a: u64, b: u64) -> u64 {
    if b == 0 {
        return a;
    }

    return recursive_sum(a+1, b-1);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works_1() {
        let result = recursive_sum(1, 2);
        assert_eq!(result, 3);
    }

    #[test]
    fn it_works_2() {
        let result = recursive_sum(123, 123);
        assert_eq!(result, 246);
    }
}

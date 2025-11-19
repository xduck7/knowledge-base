use test_example::recursive_sum;


#[cfg(test)]
mod unit_tests {
    use super::*;

    #[test]
    fn test_recursive_sum() {
        assert_eq!(recursive_sum(2, 3), 5);
        assert_eq!(recursive_sum(0, 0), 0);
        assert_eq!(recursive_sum(10, 5), 15);
        assert_eq!(recursive_sum(100, 200), 300);
    }
}
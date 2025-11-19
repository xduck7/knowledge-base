//! # Crate example
//!
//! `crate_example` - my hobby

/// Adds two integers and returns the result.
///
/// # Arguments
///
/// * a - The first integer
/// * b - The second integer
///
/// # Returns
///
/// The tests of a and b as i32
///
/// # Examples
///
/// ```
/// let expected: i32 = 5;
/// let result = crate_example::sum(2, 3);
///
/// assert_eq!(result, expected);
/// ```
pub fn sum(a: i32, b: i32) -> i32 {
    a + b
}
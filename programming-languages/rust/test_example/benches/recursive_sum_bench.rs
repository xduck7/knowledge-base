use criterion::{black_box, criterion_group, criterion_main, Criterion};
use test_example::recursive_sum;

fn add_benchmark(c: &mut Criterion) {
    c.bench_function("add small numbers", |b| {
        b.iter(|| recursive_sum(black_box(2), black_box(3)))
    });

    c.bench_function("add large numbers", |b| {
        b.iter(|| recursive_sum(black_box(u64::MAX - 10), black_box(5)))
    });

    c.bench_function("add zero", |b| {
        b.iter(|| recursive_sum(black_box(0), black_box(0)))
    });
}

fn comparison_benchmark(c: &mut Criterion) {
    let mut group = c.benchmark_group("Add Comparison");

    group.bench_function("add function", |b| {
        b.iter(|| recursive_sum(black_box(100), black_box(200)))
    });

    group.bench_function("direct addition", |b| {
        b.iter(|| {
            let a = black_box(100);
            let b = black_box(200);
            a + b
        })
    });

    group.finish();
}

criterion_group!(benches, add_benchmark, comparison_benchmark);
criterion_main!(benches);
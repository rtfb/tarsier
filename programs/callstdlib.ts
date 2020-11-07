
let mapped = map([2, 4, 6], fn(n) {
    n + 1
});

let summed = sum([1, 2, 3, 4, 5]);

puts(mapped, summed);

unless(10 > 5, puts("not greater"), puts("greater"));

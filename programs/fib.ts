
let fib = fn(n) {
    if (n == 1) {
        return 1;
    }
    if (n == 2) {
        return 1;
    }
    return fib(n - 1) + fib(n - 2);
}

puts(fib(1));
puts(fib(2));
puts(fib(3));
puts(fib(4));
puts(fib(5));
puts(fib(6));
puts(fib(7));
puts(fib(8));
puts(fib(9));
puts(fib(10));

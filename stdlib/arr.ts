
let map = fn(arr, func) {
    let iter = fn(arr, accumulated) {
        if (len(arr) == 0) {
            accumulated
        } else {
            iter(rest(arr), push(accumulated, func(first(arr))));
        }
    };
    iter(arr, []);
};

let reduce = fn(arr, initial, func) {
    let iter = fn(arr, result) {
        if (len(arr) == 0) {
            result
        } else {
            iter(rest(arr), func(result, first(arr)));
        }
    };
    iter(arr, initial);
};

let sum = fn(arr) {
    reduce(arr, 0, fn(accum, el) {
        accum + el
    })
};

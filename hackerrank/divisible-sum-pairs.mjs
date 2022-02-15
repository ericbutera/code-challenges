// https://www.hackerrank.com/challenges/divisible-sum-pairs/problem?isFullScreen=true

import assert from 'assert/strict'

/**
 *
 * @param {int} n ar length
 * @param {int} k positive integer to divide with
 * @param {int[]} ar array of integers
 * @returns {int} number of pairs
 */
const divisibleSumPairs = (n, k, ar) => {
  let count = 0;

  for (let i = 0; i < ar.length; i++) {
    let next = i + 1; // prevent comparing same number

    for (let j = next; j < ar.length; j++) {
      let isDivisible = (ar[i] + ar[j]) % k === 0;
      if (isDivisible) {
        ++count;
      }
    }
  }

  return count;
}

assert.equal(3, divisibleSumPairs(6, 5, [1,2,3,4,5,6]))

let res = divisibleSumPairs(6, 3, [1, 3, 2, 6, 1, 2]);
assert.equal(5, res);

let data = '43 95 51 55 40 86 65 81 51 20 47 50 65 53 23 78 75 75 47 73 25 27 14 8 26 58 95 28 3 23 48 69 26 3 73 52 34 7 40 33 56 98 71 29 70 71 28 12 18 49 19 25 2 18 15 41 51 42 46 19 98 56 54 98 72 25 16 49 34 99 48 93 64 44 50 91 44 17 63 27 3 65 75 19 68 30 43 37 72 54 82 92 37 52 72 62 3 88 82 71'.split(' ');
assert.equal(248, divisibleSumPairs(100, 22, data));

/*

// array of integers
ar = [1,2,3,4,5,6]

// divide by
k = 5

find (i,j) pairs by taking permutations: ar[0],ar[1]; ar[0],ar[2]; ar[0],ar[3]; ar[0],ar[4]
determine the number of pairs
i < j and ar[i] + ar[j] divisible by k

for x=0
  for y=0  n^2 quadratic
    first = ar[x]
    second = ar[y]
    pair = [first,second]

    if first < second AND  <--- oops, it was saying if the index, not value. ;)
      first + second % k == 0
      pairs.push([first,second])
*/

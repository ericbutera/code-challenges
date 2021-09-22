'use strict';

// https://www.hackerrank.com/challenges/game-of-thrones/problem

const assert = require('assert');
const { isPalindrome, LABELS } = require('./palindrome');

const gameOfThrones = string => {
    return isPalindrome(string) ? LABELS.YES : LABELS.NO;
};

const tests = [
    {
        input: 'cdcdcdcdeeeef',
        output: LABELS.YES
    },
    {
        input: 'aaabbbb',
        output: LABELS.YES
    },
    {
        input: 'cdefghmnopqrstuvw',
        output: LABELS.NO
    },
];

for (var test of tests) {
    console.debug(`input: ${test.input} expected out: ${test.output}`);
    assert.strictEqual(gameOfThrones(test.input), test.output);
}

console.log(`Success! Test count: ${tests.length} `);
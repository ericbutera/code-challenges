'use strict';
const assert = require('assert');
const { miniMaxSum } = require('./lib');

const tests = [
    {
        input: [1, 2, 3, 4, 5],
        output: "10 14"
    },
    {
        input: [1, 3, 5, 7, 9],
        output: "16 24"
    },
    {
        input: [1, 4, 3, 5, 2],
        output: "10 14"
    },
    {
        input: [7, 69, 2, 221, 8974],
        output: "299 9271"
    }
];

for (var test of tests) {
    console.debug(`compare ${test.input} = ${test.output}`);
    assert.strictEqual(miniMaxSum(test.input), test.output);
}

console.log(`Success! Test count: ${tests.length}`)
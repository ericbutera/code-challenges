'use strict';
const assert = require('assert');
const { birthdayCakeCandles } = require('./lib');

const tests = [
    {
        input: [3, 2, 1, 3],
        output: 2
    },
];

for (var test of tests) {
    console.debug(`compare ${test.input} = ${test.output}`);
    assert.strictEqual(birthdayCakeCandles(test.input), test.output);
}

console.log(`Success! Test count: ${tests.length}`)
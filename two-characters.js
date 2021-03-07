'use strict';
const assert = require('assert');

/**
 * Create an array of unique 2 letter combinations from a string.
 * 
 * @param {string} string 
 * @returns array [ [a,b], [b,c] ]
 */
const pairs = string => {
    let unique = [];
    for (let x = 0; x < string.length; x++) {
        if (unique.indexOf(string[x]) === -1)
            unique.push(string[x]);
    }
    unique.sort();

    // input  = beabeefeab
    // unique = a,b,e,f
    // pairs  = [ [a,b], [a,e], [a,f], [b,e], [b,f], [e,f] ]
    const _pairs = [];
    for (let c = 0; c < unique.length; c++) {
        let char = unique[c];
        for (let x = (c + 1); x < unique.length; x++) {
            // (c + 1) skips aa
            _pairs.push([char, unique[x]]);
        }
    }

    return _pairs;
};

/**
 * Modify the given string to only contain letters in pair.
 * pair=[a,b], string=beabeefeab => b,a,b,a,b
 * @param {array} pair Pair [a,b]
 * @param {string} string String to reduce
 * @returns string
 */
const reduceToPair = (pair, string) => {
    let clean = [];
    for (let c = 0; c < string.length; c++) {
        if (pair.includes(string[c])) // [a,b] includes `b` (from beabeefeab)
            clean.push(string[c]);
    }

    return clean.join('');
};

/**
 * Is the given string a repeating alternate?
 * @param {string} string "babab"
 * @returns {boolean}
 */
const isAlternating = string => {
    /*
    let pair = [string[0], string[1]];
    let index = 0;
    for (let x = 2; x < string.length; x++) {
        // console.log(`${string[x]} == ${pair[index]}`)
        index = 1 - index;
        if (string[x] !== pair[index]) {
            return false;
        }
    }

    return true;
    */

    // might be slower but i think this reads nicer.
    // generate a matching repeating string and compare
    let base = string.substr(0, 2); // "ba"
    let compare = base
        .repeat(Math.ceil(string.length / 2))
        .substr(0, string.length); // "babab"

    if (compare === string)
        return true;

    return false;
};

/**
 * Two characters challenge
 * 
 * @param {string} string 
 * @returns {number}
 */
const alternate = string => {
    // reduce to 2 distinct characters beabeefeab -> babab
    // then characters must alternate `babab`=pass `bbaa`=fail

    let _pairs = pairs(string);
    if (!_pairs.length)
        return 0;

    const alternating = [];
    _pairs.forEach(pair => {
        let attempt = reduceToPair(pair, string);
        if (isAlternating(attempt))
            alternating.push(attempt);
    });

    if (!alternating.length)
        return 0;

    let counts = alternating.map(string => string.length);
    let total = Math.max(...counts);

    return total;
};

const tests = [
    {
        input: 'beabeefeab',
        output: 5
    },
    {
        input: 'asdcbsdcagfsdbgdfanfghbsfdab',
        output: 8
    },
    {
        input: 'asvkugfiugsalddlasguifgukvsa',
        output: 0
    },
];

let count = 0;
for (var test of tests) {
    console.log("running test %o ", count++);

    let result = alternate(test.input);
    assert.strictEqual(result, test.output);
}

console.log(`Success! Test count: ${tests.length} `);
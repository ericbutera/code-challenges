'use strict';
const readline = require('readline');
const { miniMaxSum } = require('./lib');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

rl.question('Enter a single line of five space-separated integers.\n', line => {
    const numbers = line.split(/\s+/)
        .map(number => parseInt(number));

    const result = miniMaxSum(numbers);
    console.log(`Result: ${result}`);

    rl.close();
});

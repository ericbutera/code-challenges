/**
 * Default find length
 */
const FIND_LENGTH = 4;

/**
 * Returns two space-separated long integers denoting the respective minimum and maximum 
 * values that canbe calculated by summing exactly four of the five integers. (The output 
 * can be greater than a 32 bit integer.
 * 
 * @param {array} numbers 
 * @param {number} length 
 * @returns string "low-sum high-sum"
 */
const miniMaxSum = (numbers, length) => {
    length = length || FIND_LENGTH;

    const sorted = numbers.sort((a, b) => a - b);

    let low = 0;
    for (let i = 0; i < FIND_LENGTH; i++) {
        low += sorted[i];
    }

    let high = 0;
    for (let i = sorted.length - FIND_LENGTH; i < sorted.length; i++) {
        high += sorted[i];
    }

    return `${low} ${high}`;
};

exports.miniMaxSum = miniMaxSum;
/**
 * @param {array} candles
 * @returns {number} Number of candles that are tallest
 */
const birthdayCakeCandles = candles => {
    let highest = Math.max(...candles);
    let count = 0;

    for (const height of candles) {
        if (height === highest)
            count++;
    }

    console.log(`Highest ${highest} Count: ${count}`);
    return count;
};

exports.birthdayCakeCandles = birthdayCakeCandles;

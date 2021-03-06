/**
 * Enum 
 */
const LABELS = {
    YES: 'YES',
    NO: 'NO'
};

/**
 * Is the given string a palindrome?
 * @param {string} string Input string
 * @returns {boolean}
 */
const isPalindrome = string => {
    // there can only be 1 odd. any more than that and false
    const letters = {};

    for (let letter of string) {
        if (!letters.hasOwnProperty(letter))
            letters[letter] = 0;

        letters[letter]++;
    }

    let usedOdd = false;

    for (let count of Object.values(letters)) {
        if (count % 2 !== 0) {
            if (usedOdd)
                return false;

            usedOdd = true;
        }
    }

    return true;
};

exports.LABELS = LABELS;
exports.isPalindrome = isPalindrome;
// https://leetcode.com/problems/reverse-words-in-a-string/

/*
Input: "the sky is blue"
Output: "blue is sky the"

Input: "  hello world!  "
Output: "world! hello"
Explanation: Your reversed string should not contain leading or trailing spaces.

Input: "a good   example"
Output: "example good a"
Explanation: You need to reduce multiple spaces between two words to a single space in the reversed string.

Note:
A word is defined as a sequence of non-space characters.
Input string may contain leading or trailing spaces. However, your reversed string should not contain leading or trailing spaces.
You need to reduce multiple spaces between two words to a single space in the reversed string.
*/

const assert = require('assert');

/**
 * @param {string} s
 * @return {string}
 */
var reverseWords = function(s) {
    let words = createWords(s);

    let reversed = [];
    while (words.length) 
        reversed.push( words.pop() );

    return reversed.join(' ');
};

const createWords = function(string) {
    let inword = false;
    let words = [];
    let word = '';

    let last = string.length - 1;

    // use indexOf(' ', start)?
    for (let x=0; x < string.length; x++) {
        let char = string[x];

        if (char == ' ') {
            if (inword) { // finish word & reset
                inword = false;
                words.push(word);
                word = '';
            }

            continue;
        } else {
            // letter
            inword = true;
            word += char;
        }

        // if at end make sure to ad word to words
        if (x == last)
            words.push(word);
    }

    return words;
};

assert.equal(reverseWords("a good   example"), "example good a");
assert.equal(reverseWords("  hello world!  "), "world! hello");
assert.equal(reverseWords("the sky is blue"), "blue is sky the");
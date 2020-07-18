/**
 * convert array of words into justified line 
 * @param {string[]} words 
 * @param {int} maxWidth 
 */
const justify = (words, length, maxWidth) => {
    let toFill = (maxWidth - length);
    //console.log("justify len %o max %o toFill %o words %o", length, maxWidth, toFill, words);

    let position = 0; // current word offset
    let maxoffset = words.length - 2; // 2 = max offset, ignore last word

    // TODO enhance 
    // there has to be a better way to distribute whitespace. 
    // figure out gaps. 
    // divide toFill against gaps. 
    // use remainder to determine if the last space should be truncated 
    // 'Science...is...what..we' - notice the two dots on the last gap

    // while toFill, cycle thru words appending whitespace
    // do not append whitespace to last word in line
    for (let x=0; x < toFill; x++) {
        if (position > maxoffset) {
            position = 0;
        }

        words[position] = words[position] + ' '; 

        position++; 
    }

    return words.join('');
};

/**
 * group disparate words into a singular line
 * @param {string[]} words 
 * @param {int} maxWidth 
 */
const wordsToLines = (words, maxWidth) => {
    let lines = [];
    let lastOffset = (words.length - 1);
    //console.log("last offset %o", lastOffset);

    let line = [];
    let lineLength = 0;  // total length of characters + whitespace
    let wordsLength = 0; // total length of characters. no whitespace

    for (let index = 0; index < words.length; index++) {
        let word = words[index];
        let wordLength = word.length;
        let attempt = wordLength + lineLength; // whitespace + wordlength + current
        //console.log("index %o word %o len %o attempt %o", index, word, wordLength, attempt);

        //if (word.length > maxWidth) throw Error('version2 will address word > maxWidth');

        if (attempt > maxWidth) {
            //console.log("overflow, create new line\n for %o", word);
            // overflow, make new line
            lines.push(justify(line, wordsLength, maxWidth)); 
            line = [];
            lineLength = 0;
            wordsLength=0;
        }

        line.push(word);
        wordsLength += wordLength;
        lineLength += wordLength + 1; 
        //console.log("wordsLength %o lineLength %o line %o", wordsLength, lineLength, line); // TODO merge words/line length

        // For the last line of text, it should be left justified and no extra space is inserted between words.
        if (index == lastOffset) 
            lines.push(justifyLast(line, maxWidth));
    }

    return lines;
};

var justifyLast = (line, maxWidth) => {
    //console.log("last line %o", line);
    let last = line.join(' ');
    let remain = maxWidth - last.length;

    if (remain) 
        last += ' '.repeat(remain);
    
    return last;
};

const fullJustify = (words, maxWidth) => {
    let lines = wordsToLines(words, maxWidth);

    console.log('-'.repeat(10));
    lines.forEach(function(line) {
        console.log("'" + line + "'");
    });
    console.log('-'.repeat(10));
    
    return lines;
};

fullJustify(["This", "is", "an", "example", "of", "text", "justification."], 16);
fullJustify(["Science", "is", "what", "we", "understand", "well", "enough", "to", "explain", "to", "a", "computer.", "Art", "is", "everything", "else", "we", "do"], 20);
fullJustify(["What","must","be","acknowledgment","shall","be"], 16);
fullJustify(["Listen","to","many,","speak","to","a","few."], 6);
//fullJustify(["a"], 1, ["a"]);
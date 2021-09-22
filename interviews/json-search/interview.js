const _ = require('lodash');

/**
 * Overload key to search all values
 */
const CRITERIA_ANY = '_any';

/**
 * Search criteria
 */
class Criteria {
    /**
     * Createa search criteria
     * @param {string} key 
     * @param {string} value
     */
    constructor(key, value) {
        /**
         * Collection key to search against.
         */
        this.key = key;

        /**
         * Search term
         */
        this.value = value;
    }
}

/**
 * es6, case sensitive, partial match
 * @param {string} haystack String to search
 * @param {string} needle Search term 
 * @returns {boolean}
 */
const containsString = (haystack, needle) => haystack.includes(needle);
//return row[criteria.key].includes(criteria.value); // 
//return row[criteria.key].indexOf(criteria.value) !== -1; // es5, case sensitive, partial match
//return row[criteria.key] === criteria.value; // strict case-sensitive exact-match

/**
 * es6 case insensitive compare. 
 * `String.localeCompare(input, 'en', {sensitivity:'base'})`
 * @param {string} haystack String to search
 * @param {string} needle Search term
 * @returns  {boolean}
 */
const containsStringInsensitive = (haystack, needle) =>
    haystack.localCompare(needle, 'en', { sensitivity: 'base' });

/**
 * Search a dataset using one or more `Criteria`
 * @param {array<object>} collection Dataset
 * @param {array<Criteria>|Criteria} criteria Collection of criterian or single criteria
 * @returns 
 */
const search = (collection, criteria) => {
    if (Array.isArray(criteria))
        return multiSearch(collection, criteria);

    return singleSearch(collection, criteria);
};

/**
 * Search a dataset using a single criteria.
 * @param {array<object>} collection City records
 * @param {Criteria} criteria Search
 * @returns {array<object>} Results matching search criteria
 */
const singleSearch = (collection, criteria) => {
    return collection.filter(row => {
        if (criteria.key === CRITERIA_ANY && wildcardSearch(row, criteria.value))
            return true;

        if (!row.hasOwnProperty(criteria.key))
            return false; // exclude rows missing key

        return containsString(row[criteria.key], criteria.value);
    });
};

/**
 * Search a dataset with multiple criteria. Works by progressively applying filters against the data set.
 * @param {array<object>} data Dataset
 * @param {array<Criteria>} criterian Criteria to filter by
 * @returns {array<object>} Filtered data set
 */
const multiSearch = (data, criterian) => {
    criterian.forEach(criteria => {
        data = singleSearch(data, criteria);
    });

    return data;
};

/**
 * Do any of the values in `row` match `value`? 
 * @param {object} row Object to search
 * @param {string} value Value to search for
 * @returns {boolean}
 */
const wildcardSearch = (row, value) => {
    const vals = Object.values(row);
    for (let x = 0; x < vals.length; x++) {
        if (containsString(vals[x], value))
            return true;
    }

    return false;
};

/**
 * Convert array of parameters into key=>value pairs
 * `[key1, val1, key2, val2] becomes [[key1,val1], [key2,val2]]
 * @param {array} argv input from argvToArgs
 * @returns {array}
 */
const argsToPairs = args => {
    const pairs = [];

    if (!args || args.length === 0)
        return pairs;

    if (args.length % 2 !== 0)
        args.push(null); // ensure even pairs

    let key, val = null;

    for (let x = 0; x < args.length; x += 2) {
        key = args[x];     // state
        val = args[x + 1]; // illinois
        pairs.push(new Criteria(key, val));
    }

    return pairs;
};

/**
 * Returns passed in parameters to process using process.argv.
 * `$program.js 1 2 3` becomes `[1,2,3]`
 * This isn't used in the demo. It would be used if this were to be wired up into a real CLI program
 * @param {array} args 
 * @returns {array}
 */
const argvToArgs = args => {
    /*
    $ node argv.js one two three

    argv [
    'C:\\Program Files\\nodejs\\node.exe',
    'C:\\Users\\ericb\\code\\code-challenges\\interview-search\\argv.js',
    'one',
    'two',
    'three',

     */
    return process.argv.slice(2);
};

exports.search = search;
exports.argsToPairs = argsToPairs;
//exports.multiSearch = multiSearch;
exports.Criteria = Criteria;
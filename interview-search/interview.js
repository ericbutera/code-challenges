const _ = require('lodash');

/**
 * Search criteria
 */
class Criteria {
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
 * 
 * @param {array<object>} collection City records
 * @param {Criteria} criteria Search
 * @returns {array<object>} Results matching search criteria
 */
const search = (collection, criteria) => {
    //if (criteria.key == '_any') // special wildcard search
    //wildcard();

    return collection.filter(row => {
        if (!row.hasOwnProperty(criteria.key))
            return false; // exclude rows missing key

        return row[criteria.key].includes(criteria.value);
    });
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

/**
 * Convert array of parameters into key=>value pairs
 * `[key1, val1, key2, val2] becomes [[key1,val1], [key2,val2]]
 * @param {array} argv input from argvToArgs
 * @returns {array}
 */
const argsToPairs = args => {
    if (args.length % 2 !== 0)
        args.push(null);

    return _.chunk(args, 2);
};

const multi = args => {
    let dataset = _.clone(data);

    // todo replace args with argvToPairs
    let parts = args.split(' ');
    let searches = _.chunk(parts, 2);

    console.log("args to search using apply...");
    searches.forEach(current => {
        console.log("search %o", current);
        dataset = search(current[0], current[1], dataset);
    });

    debugger;


    /*for (let x = 0; x < parts.length; x += 2) {
        let key = parts[x] // state
        let val = parts[x + 1] // ill

        dataset = search(key, val, dataset);
    }*/
    return dataset;
};

exports.search = search;
exports.argsToPairs = argsToPairs;
exports.multi = multi;
exports.Criteria = Criteria;
const test = require('ava');
const { search, argsToPairs, multi, Criteria } = require('./interview');

// command line:
// interview.js city Chicago
// input i ruby json - search.rb city Chicago
// output [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }]

// Object.is means [] !== [] https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is

test('argsToSearch city Chicago', t => {
    let pairs = argsToPairs(['city', 'Chicago']);
    t.deepEqual(pairs, [['city', 'Chicago']])
});

test('argsToSearch city Chicago landmark Tower', t => {
    let pairs = argsToPairs(['city', 'Chicago', 'landmark', 'Tower']);
    t.deepEqual(pairs, [['city', 'Chicago'], ['landmark', 'Tower']])
});

test('argsToSearch _any Tower', t => {
    let pairs = argsToPairs(['_any', 'Tower']);
    t.deepEqual(pairs, [['_any', 'Tower']])
});

test('argsToSearch one param only fails', t => {
    let pairs = argsToPairs(['any']);
    t.deepEqual(pairs, [['any', null]]);
})

test('search city Chicago', t => {
    let criteria = new Criteria('city', 'Chicago');
    let result = search(data, criteria);
    t.deepEqual(result, [data[0]]);
});

//test('integration city Chicago returns result', t => {
//});

const data = [
    { "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" },
    { "city": "Springfield", "state": "Illinois" },
    { "city": "New York City", "state": "New York", "landmark": "Empire State Building" },
    { "city": "Paris", "landmark": "Eiffel Tower" },
    { "city": "Paris Tower", "landmark": "Eiffel" },
];

//console.log(_.isEqual(result, expected));
//console.log('test2 %o', _.isEqual(search('landmark', 'Eiffel Tower'), [{ "city": "Paris", "landmark": "Eiffel Tower" }]));
//console.log('test 3 %o', _.isEqual(search('state', 'Illinois'), [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }, { "city": "Springfield", "state": "Illinois" }]));
//console.log('test 4 %o', _.isEqual(search('landmark', 'Tower'), [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }, { "city": "Paris", "landmark": "Eiffel Tower" }]));
//console.log('test %o', _.isEqual(argsToSearch('landmark Tower state Illinois'), [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }]));
//console.log('test %o', _.isEqual(search('_any', 'Tower'), [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }, { "city": "Paris", "landmark": "Eiffel Tower" }, { "city": "Paris Tower", "landmark": "Eiffel" }]));
//city Chicago
//_any Searchterm
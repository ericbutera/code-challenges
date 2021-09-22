const test = require('ava');
const { search, argsToPairs, Criteria } = require('./interview');

// command line:
// interview.js city Chicago
// input i ruby json - search.rb city Chicago
// output [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }]

// Object.is means [] !== [] https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is

test('argsToSearch city Chicago', t => {
    const pairs = argsToPairs(['city', 'Chicago']);
    const criteria = new Criteria('city', 'Chicago');
    t.deepEqual(pairs, [criteria]);
});

test('argsToSearch city Chicago landmark Tower', t => {
    const pairs = argsToPairs(['city', 'Chicago', 'landmark', 'Tower']);
    const cri1 = new Criteria('city', 'Chicago');
    const cri2 = new Criteria('landmark', 'Tower');
    t.deepEqual(pairs, [cri1, cri2]);
});

test('argsToSearch _any Tower', t => {
    const pairs = argsToPairs(['_any', 'Tower']);
    const criteria = new Criteria('_any', 'Tower');
    t.deepEqual(pairs, [criteria]);
});

test('argsToSearch one param only fails', t => {
    const pairs = argsToPairs(['any']);
    const criteria = new Criteria('any', null);
    t.deepEqual(pairs, [criteria]);
})

test('argsToSearch no params finds nothing', t => {
    t.deepEqual(argsToPairs(), []);
});

test('search city Chicago finds one record', t => {
    const criteria = new Criteria('city', 'Chicago');
    const result = search(data, criteria);
    t.deepEqual(result, [data[0]]);
});

test('search landmark Eiffel Tower finds one record', t => {
    // multi-word criteria isn't supported. the work around is to only search 1 term such as 'Tower' for 'Sears Tower'
    // to fix the command line parser should handle quoted strings. im confident there's an npm package for that...
    // known issue; ignored for now
    const criteria = new Criteria('landmark', 'Eiffel Tower');
    const result = search(data, criteria);
    t.deepEqual(result, [{ "city": "Paris", "landmark": "Eiffel Tower" }]);
});

test('search state Illiniois finds two records', t => {
    const criteria = new Criteria('state', 'Illinois');
    const actual = search(data, criteria);
    const expected = [
        { "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" },
        { "city": "Springfield", "state": "Illinois" }
    ];
    t.deepEqual(actual, expected);
});

test('search handle partial match', t => {
    const criteria = new Criteria('landmark', 'Tower');
    const actual = search(data, criteria);
    const expected = [
        { "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" },
        { "city": "Paris", "landmark": "Eiffel Tower" }
    ];
    t.deepEqual(actual, expected);
});

test('search multiple criteria against single record', t => {
    const criterian = [
        new Criteria('landmark', 'Tower'),
        new Criteria('state', 'Illinois'),
    ];
    const actual = search(data, criterian);
    const expected = [{ "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" }];
    t.deepEqual(actual, expected);
});

test('search using _any Tower', t => {
    const criteria = new Criteria('_any', 'Tower');
    const actual = search(data, criteria);
    const expected = [
        { "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" },
        { "city": "Paris", "landmark": "Eiffel Tower" },
        { "city": "Paris Tower", "landmark": "Eiffel" }
    ];
    t.deepEqual(actual, expected);
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
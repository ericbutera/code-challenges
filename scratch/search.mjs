import assert from 'assert/strict';
import { hasUncaughtExceptionCaptureCallback } from 'process';

const data = [
  { "city": "Chicago", "state": "Illinois", "landmark": "Sears Tower" },
  { "city": "Springfield", "state": "Illinois" },
  { "city": "New York City", "state": "New York", "landmark": "Empire State Building" },
  { "city": "Paris", "landmark": "Eiffel Tower" },
  { "city": "Paris Tower", "landmark": "Eiffel" },
];

// find city = 'chicago'
const findCity = (data, city) => {
  return data.filter(row => row['city'] == city);
}

// orig attempt before `_any`
const findPartial = (data, pairs) => {
  // find pairs in data
  // each row in data must contain each value in pairs
  return data.filter(row => {
    // let exists = true; // AND means each pair must contain a match; if any don't set to false
    // pairs.forEach(pair => {}) // don't use forEach as it cannot be exited early on result not found

    // BUG no matches returns all

    for (let x = 0; x < pairs.length; x++) {
      let pair = pairs[x];
      if (!row[pair.key]) continue;

      // data[pair.key] == pair.value         // case sensitive, exact match
      // !row[pair.key].includes(pair.value)  // case sensitive, partial match

      if (row[pair.key].indexOf(pair.value) === -1) // case insensitive, partial match
        return false; // exit immediately if a match isnt found
    }

    return true;
  });

}

const findPartialAny = (data, pairs) => {
  // find pairs in data
  // each row in data must contain each value in pairs

  // pairs
  // [{city:Chicago},{_any:Tower}]
  // [{landmark:Building},{_any:York}]
  return data.filter(row => {
    let hasOne = false;

    for (let x = 0; x < pairs.length; x++) {
      let pair = pairs[x];

      let match = false;

      if (pair.key == '_any') {
        match = searchValue(row, pair.value)
      } else {
        match = searchKeyValue(row, pair.key, pair.value)
      }

      if (!match)
        return false; // exit early

      hasOne = true; // record there was at least one match so an empty resultset doesnt return all
    }

    return hasOne;
  });
}

// search for a specific key value
const searchKeyValue = (row, key, value) => {
  if (!row[key]) return false;

  if (row[key].indexOf(value) !== -1)
    return true; // found match

  return false;
}

// search all keys for value
const searchValue = (row, value) => {
  for (let key in row) {
    if (row[key].indexOf(value) !== -1)
      return true; // found match
  }
  return false;
}

assert.deepStrictEqual(
  findCity(data, 'Chicago'),
  [
    { city: 'Chicago', state: 'Illinois', landmark: 'Sears Tower' }
  ]
)

/*
exposes bug with no found results returns all:
assert.deepStrictEqual(
  findPartial(data, [{ 'key': 'fake-key', 'value': 'fake-value' }]),
  []
)
*/

assert.deepStrictEqual(
  findPartialAny(data, [{ 'key': 'fake-key', 'value': 'fake-value' }]),
  []
)

assert.deepStrictEqual(
  findPartialAny(data, [
    { 'key': 'city', 'value': 'Chicago' }
  ]),
  [
    { city: 'Chicago', state: 'Illinois', landmark: 'Sears Tower' }
  ]
)

assert.deepStrictEqual(
  findPartialAny(data, [
    { 'key': 'landmark', 'value': 'Tower' }
  ]),
  [
    { city: 'Chicago', state: 'Illinois', landmark: 'Sears Tower' },
    { city: 'Paris', landmark: 'Eiffel Tower' }
  ]
)

assert.deepStrictEqual(
  findPartialAny(data, [
    { 'key': 'city', 'value': 'Paris' },
    { 'key': 'landmark', 'value': 'Eiffel' }
  ]),
  [
    { city: 'Paris', landmark: 'Eiffel Tower' },
    { city: 'Paris Tower', landmark: 'Eiffel' }
  ]
)

assert.deepStrictEqual(
  findPartialAny(data, [
    { 'key': 'city', 'value': 'Paris' },
    { 'key': 'landmark', 'value': 'Eiffel' }
  ]),
  [
    { city: 'Paris', landmark: 'Eiffel Tower' },
    { city: 'Paris Tower', landmark: 'Eiffel' }
  ]
)

assert.deepStrictEqual(
  findPartialAny(data, [{ 'key': '_any', 'value': 'Tower' }]),
  [
    { city: 'Chicago', landmark: 'Sears Tower', state: 'Illinois' },
    { city: 'Paris', landmark: 'Eiffel Tower' },
    { city: 'Paris Tower', landmark: 'Eiffel' }
  ]
)

assert.deepStrictEqual(
  findPartialAny(data, [
    { 'key': '_any', 'value': 'Tower' },
    { 'key': 'city', 'value': 'Paris' }
  ]),
  [
    { city: 'Paris', landmark: 'Eiffel Tower' },
    { city: 'Paris Tower', landmark: 'Eiffel' }
  ]
)

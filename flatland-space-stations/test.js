'use strict';
const assert = require('assert');

/*
 * Scenerio:
 * cities = 5
 * stations = [0, 4]
 *
 * city hasStation  distance-left   distance-right  nearest-station-city        max-distance
 * ---- ----------  -------------   --------------  --------------------        ------------
 * 0    1           none            4               0 (self)                    0km
 * 1    0           1               3               0 (left)                    1km
 * 2    0           2               2               0, 4 (left or right equal)  2km
 * 3    0           3               1               4 (right)                   1km
 * 4    1           4               none            4 (self)                    0km
 * 
 * Calculate max-distance: max(0, 1, 2, 1, 0) = 2
 * return 2
 */

/**
 * Calculate the maximum distance any city is from a space station
 * @param {number} cities Number of cities
 * @param {Array<number>} stations Cities with a space station, 1-based indexing
 * @returns {number} 
 */
const flatlandSpaceStations = (cities, stations) => {
    stations.sort((a, b) => a - b);

    let maximum = 0
        , leftStation = 0
        , rightStation = 0;

    for (let city = 0; city < cities; city++) {
        let distances = []
            , isStation = false
            , left = null
            , right = null;

        if (city == stations[leftStation] || city == stations[rightStation]) {
            isStation = true;
            distances.push(0);
        }

        // attempt to discover if there is a space station to the right
        if (city == stations[rightStation]) {
            if (stations.length > rightStation + 1)
                rightStation++;
        }

        if (stations[rightStation] > city) {
            right = stations[rightStation] - city;
            distances.push(right);
        }

        // attempt to discover if there is a closer "left-station" to the current city
        if (city > stations[leftStation + 1]) {
            leftStation++;
        }

        if (city > stations[leftStation]) {
            left = city - stations[leftStation];
            distances.push(left);
        }

        let distance = Math.min(...distances);
        if (distance > maximum)
            maximum = distance;

        //console.log(`city: ${city} isStation: ${isStation} left: ${left} right: ${right} left-station ${stations[leftStation]} right-station ${stations[rightStation]}`);
    }

    return maximum;
};

const tests = [
    {
        input: { cities: 100000, stations: [39572, 89524, 21749, 94613, 75569, 74800, 91713, 62107, 28574, 22617, 22271, 22624, 28116, 67573, 53717, 9358, 65220, 59894, 78686, 10945, 33641, 11708, 8851, 11860, 66780, 64697, 799, 47782, 41971, 54170, 8960, 81543, 60047, 47061, 92508, 51968, 38213, 84221, 14075, 66787, 23191, 52698, 5764, 51307, 20271, 59481, 77018, 1843, 19375, 55704, 12789, 53016, 83764, 37992, 64877, 50545, 19041, 82028, 98327, 61012, 52551, 7287, 42555, 12598, 70700, 51416, 80918, 8914, 35637, 11345, 75701, 58828, 80395, 97817, 26488, 17019, 57299, 3506, 18862, 93026, 75562, 48003, 62395, 59327, 85996, 27272, 9872, 5037, 25652, 8199, 82402, 78203, 31838, 41309, 7153, 18890, 92725, 88071, 27804, 28363, 99416, 19858, 3543, 79812, 17675, 30031, 96831, 91326, 49889, 15693, 84353, 25452, 80049, 46748, 84779, 66045, 90372, 94651, 87434, 16024, 19202, 69836, 94228, 67392, 27498, 1381, 86282, 20223, 5805, 14087, 48586, 5221, 50297, 68482, 85033, 67972, 98513, 98216, 59299, 48403, 30262, 60004, 73855, 10311, 6752, 74986, 92708, 13476, 85989, 96494, 29500, 5191, 82683, 40080, 88935, 10181, 57814, 75217, 30404, 63619, 5656, 95343, 68840, 55953, 63825, 70226, 23926, 62338, 68442, 99577, 27093, 15056, 59581, 17300, 25367, 82685, 92286, 34427, 96161, 78275, 30922, 25661, 99818, 13605, 82094, 88753, 23786, 39908, 80323, 54190, 3527, 85979, 65885, 72367, 41933, 29710, 58945, 82211, 8401, 43740] },
        output: 1504
    },
    {
        input: { cities: 5, stations: [1, 3, 4] },
        output: 1
    },
    {
        input: { cities: 100, stations: [93, 41, 91, 61, 30, 6, 25, 90, 97] },
        output: 14
    },
    {
        input: { cities: 5, stations: [0, 4] },
        output: 2
    },
    {
        input: { cities: 3, stations: [0] },
        output: 2
    },
    {
        input: { cities: 6, stations: [0, 1, 2, 4, 3, 5] },
        output: 0
    },
];

let count = 0;
for (var test of tests) {
    console.log("running test %o ", count++);
    let result = flatlandSpaceStations(test.input.cities, test.input.stations);
    assert.strictEqual(result, test.output);
}

console.log(`Success! Test count: ${tests.length} `);
/**
 * Simple bubble sort
 * @param {array} numbers 
 * @returns {array} Sorted in place
 */
const sort = (numbers) => {
    let swap = null;

    for (let i = 0; i < numbers.length; i++) {
        for (let current = 0; current < (numbers.length - i - 1); current++) {
            let next = current + 1;

            if (numbers[current] > numbers[next]) {
                swap = numbers[current];
                numbers[current] = numbers[next]; // overwrite current with next
                numbers[next] = swap; // update next with old current
            }
        }
    }

    return numbers;
};

export default sort;

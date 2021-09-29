# https://www.hackerrank.com/challenges/compare-the-triplets/problem

#
# Complete the 'compareTriplets' function below.
#
# The function is expected to return an INTEGER_ARRAY.
# The function accepts following parameters:
#  1. INTEGER_ARRAY a
#  2. INTEGER_ARRAY b
#

##
# @param a array of integers
# @param b array of integers
def compareTriplets(a, b)
    a_sum = 0
    b_sum = 0

    max_index = [a.size, b.size].min # _should_ always be 3
    max_range = max_index - 1 

    (0..max_range).each do |position|
        a_sum += 1 if a[position] > b[position]
        b_sum += 1 if b[position] > a[position]
    end

    [a_sum, b_sum]
end
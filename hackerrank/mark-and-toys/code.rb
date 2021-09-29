# https://www.hackerrank.com/challenges/mark-and-toys/problem

##
# Complete the 'maximumToys' function below.
#
# The function is expected to return an INTEGER.
# The function accepts following parameters:
#  1. INTEGER_ARRAY prices
#  2. INTEGER k
#
# @param prices the toy prices
# @param budget 
def maximumToys(prices, budget)
    # sort prices
    # iterate prices => price while total less than budget

    prices = prices.sort

    total = 0
    toy_count = 0

    prices.each do |price|
        attempt = total + price

        return toy_count if attempt > budget

        total = attempt
        toy_count += 1
    end

    toy_count
end

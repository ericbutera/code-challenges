##
# Convert array of numbers into string of numbers, separated by addition operator
# +numbers+ Array of integers
def format_as_addition(numbers)
    numbers.join(' + ')
end

##
# +numbers+ Array of integers
def sum(numbers)
    numbers.sum
end

##
# Pretty print an array of numbers as a math sum equation
# eg: '[1,2,3] = 1 + 2 + 3 = 6'
# +numbers+ Array of integers
def array_sum_as_string(numbers)
  string = format_as_addition(numbers)
  sum = sum(numbers)

  "#{string} = #{sum}"
end

##
# hackerrank method signature for this problem
def simpleArraySum(numbers)
  array_sum_as_string(numbers)
end
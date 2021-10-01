# https://www.hackerrank.com/challenges/ctci-array-left-rotation/problem

##
# @param array array to rotate 
# @param int number of rotations 
def rotLeft(array, rotations)
  rotate_left_using_stdlib(array, rotations)
end

def rotate_left_using_stdlib(array, rotations)
  # creates copies of array in memory, easy to read/understand
  for _ in (1..rotations) do
    array = array.push(array.shift)
  end

  array
end

def rotate_left_using_recursion(array, rotations)
  return array if rotations == 0
  array = rotate_left_using_stdlib(array, rotations)
  rotate_left_using_recursion(array, rotations - 1)
end

def rotate_left_in_place(array, rotations)
  for rotation in (0..rotations) do
    storage = nil
    max = array.length - 1
    for x in (0..max) do 
      storage = array[x] if x == 0 
      right = x + 1
      array[x] = array[right] 
    end

    array[max] = storage
  end

  array
end

RSpec.describe 'rotate left' do
  it 'rotates' do
    array = [1, 2, 3, 4, 5]
    rotations = 4
    rotated = [5, 1, 2, 3, 4]

    expect(rotate_left_using_stdlib(array, rotations)).to eq(rotated)
    expect(rotate_left_using_recursion(array, rotations)).to eq(rotated)
    expect(rotate_left_in_place(array, rotations)).to eq(rotated)
  end

  it 'solves input 00' do
    # input 00
    #5 4
    #output 00
    #5 1 2 3 4
    array = [1, 2, 3, 4, 5]
    expect(rotLeft(array, 4)).to eq([5, 1, 2, 3, 4])
  end

  it 'solves input 01' do
    # input 01
    #20 10
    #41 73 89 7 10 1 59 58 84 77 77 97 58 1 86 58 26 10 86 51
    # output 01
    #77 97 58 1 86 58 26 10 86 51 41 73 89 7 10 1 59 58 84 77
    array = [41, 73, 89, 7, 10, 1, 59, 58, 84, 77, 77, 97, 58, 1, 86, 58, 26, 10, 86, 51]
    result = [77, 97, 58, 1, 86, 58, 26, 10, 86, 51, 41, 73, 89, 7, 10, 1, 59, 58, 84, 77]
    expect(rotLeft(array, 10)).to eq(result)
  end

  it 'solves input 10' do
    array = [33, 47, 70, 37, 8, 53, 13, 93, 71, 72, 51, 100, 60, 87, 97]
    result = [87, 97, 33, 47, 70, 37, 8, 53, 13, 93, 71, 72, 51, 100, 60]
    expect(rotLeft(array, 13)).to eq(result)
  end
end



# input 10
#15 13
#33 47 70 37 8 53 13 93 71 72 51 100 60 87 97
# output 10
#87 97 33 47 70 37 8 53 13 93 71 72 51 100 60
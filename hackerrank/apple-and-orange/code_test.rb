# https://www.hackerrank.com/challenges/apple-and-orange/problem?isFullScreen=true

##
# @param {int} s house start point
# @param {int} t house end point
# @param {int} a apple tree position
# @param {int} b orange tree position
# @param {Integer[]} apples Apples
# @param {Integer[]} oranges
def countApplesAndOranges(s, t, a, b, apples, oranges)
  apples, oranges = count_apples_oranges(s, t, a, b, apples, oranges)
  puts "#{apples}\n#{oranges}"
end

##
# wrapper for countApplesAndOranges which returns the result
def count_apples_oranges(start, ending, apple_tree, orange_tree, apples, oranges)
  # Could also explore some modeling:
  # House { start:, end: }
  # Fallen { Tree, Fruits }
  # Distances { Tree, Fruits }
  apples = find(apples, apple_tree, start, ending)
  oranges = find(oranges, orange_tree, start, ending)
  return [apples.size, oranges.size]
end

##
# @param {Integer[]} fruits Fruit distances
# @param {int} tree fruit tree position
# @param {int} start house start position
# @param {int} end hosue end position
def find(fruits, tree, start, ending)
    resting = find_resting(fruits, tree)
    find_between(resting, start, ending)
end

##
# generate positions of fallen fruit from it's tree
# @param {Integer[]} fruits Distance the fruit traveled from tree
# @returns {Integer[]} Fruit ground positions
def find_resting(fruits, tree)
  fruits.map do |distance|
    distance + tree
  end
end

def find_between(positions, start, ending)
    positions.select do |position|
      position >= start && position <= ending
    end
end

RSpec.describe 'find_resting' do
  let(:resting) { find_resting(positions, tree) }

  context 'apple tree' do
    let(:tree) { 4 }
    let(:positions) { [2, 3, -4] }

    it 'example 0' do
      expect(resting).to eq([6, 7, 0])
    end
  end

  context 'orange tree' do
    let(:tree) { 12 }
    let(:positions) { [3, -2, -4] }

    it 'example 0' do
      expect(resting).to eq([15, 10, 8])
    end
  end
end

RSpec.describe 'find_between' do
  let(:trees) { [6, 7, 0] }
  let(:start) { 7 }
  let(:ending) { 10 }

  it 'handles example 0' do
    expect(find_between(trees, start, ending)).to eq([7])
  end
end

RSpec.describe 'countApplesAndOranges' do
  it 'example' do
    # input
    s = 7   # house start
    t = 10  # house end
    a = 4   # apple tree
    b = 12  # orange tree
    m = 3   # apple count
    n = 3   # orange count
    apples = [2, 3,-4]
    oranges = [3, -2, -4]
    # output
    # 1
    # 2
    res = count_apples_oranges(s, t, a, b, apples, oranges)
    expect(res).to eq([1, 2])

    res = with_captured_stdout { countApplesAndOranges(s, t, a, b, apples, oranges) }
    expect(res).to eq("1\n2\n")
  end

  it 'example integration' do
  end

  it 'sample input 0' do
    # input
    # 7 11    # s t
    # 5 15    # a b
    # 3 2     # m n
    # -2 2 1  # apples
    # 5 -6    # oranges
    #
    # output
    # 1       # apple count
    # 1       # orange count
  end
end

# https://stackoverflow.com/a/22777806/261272
def with_captured_stdout
  original_stdout = $stdout  # capture previous value of $stdout
  $stdout = StringIO.new     # assign a string buffer to $stdout
  yield                      # perform the body of the user code
  $stdout.string             # return the contents of the string buffer
ensure
  $stdout = original_stdout  # restore $stdout to its previous value
end

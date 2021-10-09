# https://www.hackerrank.com/challenges/counting-valleys/problem

RSpec.describe 'counting valleys' do
  it 'solves 0' do
    res = countingValleys(8, 'UDDDUDUU')
    expect(res).to eq(1)
  end

  it 'solves input 01' do
    res = countingValleys(12, 'DDUUDDUDUUUD')
    expect(res).to eq(2)
  end

  def countingValleys(steps, path)
    valley_count = 0
    height = 0

    path.each_char do |direction|
      case direction
      when 'U'
        valley_count += 1 if height == -1
        height += 1
      when 'D'
        height -= 1
      end
    end

    valley_count
  end
end

# begin tracking when 

# A valley is a sequence of consecutive steps below sea level, starting with a step down from sea level and ending with a step up to sea level.
# Example:
#  1
#  0 _/\      _     <-- sea level
# -1    \    /          < valley
# -2     \/\/           < valley

# ignoring above:
# 0, 0, -1 <-- fail cause step down from 1 up is still at sea-level
# -  U  D   D   D   U   D   U  U  -
# 0, 1, 0, -1, -2, -1, -2, -1, 0  0
#           ^                  ^
#           `- begin valley    |
#                              `- end valley   valley_count++
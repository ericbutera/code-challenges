# https://www.hackerrank.com/challenges/electronics-shop/problem

RSpec.describe 'electronics shop' do
  it 'solves example 0' do
    b = 10
    keyboards = [3, 1]
    drives = [5, 2, 8]
    expect(getMoneySpent(keyboards, drives, b)).to eq(9)
  end

  it 'solves example 1' do
    b = 5
    keyboards = [4]
    drives = [5]
    expect(getMoneySpent(keyboards, drives, b)).to eq(-1)
  end

  it 'solves example 00' do
    b = 10
    keyboards = [3, 1]
    drives = [5, 2, 8]
    expect(getMoneySpent(keyboards, drives, b)).to eq(9)
  end

  it 'solves input01' do
    b = 5
    keyboards = [4]
    drives = [5]
    expect(getMoneySpent(keyboards, drives, b)).to eq(-1)
  end

  def getMoneySpent(keyboards, drives, budget)
    # maximum int that can be spent, or -1 if not possible to buy both
    keyboards = keyboards.sort!.reverse!
    drives = drives.sort!.reverse!
    max = -1

    keyboards.each do |keyboard|
      drives.each do |drive| 
        attempt = keyboard + drive 
        max = attempt if budget >= attempt && attempt > max
      end
    end

    max
  end
end

# 3 1
# 8 5 2
# budget=10
# 0: key=3 drive=8 = 11, no
# 1: key 3 drive=5 = 8, yes

# maximum combination
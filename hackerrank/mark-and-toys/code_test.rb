require './code'

#  Given a list of toy prices and an amount to spend, determine the maximum number of gifts he can buy.
RSpec.describe 'mark and toys' do
  it 'finds 3 items' do
    prices = [1,2,3,4]
    budget = 7
    expect(maximumToys(prices, budget)).to eq(3)
  end

  it 'finds 4 gifts' do
    prices = [1, 12, 5, 111, 200, 1000, 10]
    budget = 50
    expect(maximumToys(prices, budget)).to eq(4)
  end

  it 'finds 3 gifts for 15 budget' do
    prices = [3, 7, 2, 9, 4]
    budget = 15
    expect(maximumToys(prices, budget)).to eq(3)
  end
end


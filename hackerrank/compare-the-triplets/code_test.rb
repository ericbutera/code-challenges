require './code'

RSpec.describe 'compare the triplets' do
    it 'compares' do
        a = [5, 6, 7]
        b = [3, 6, 10]
        result = compareTriplets(a, b)
        expect(result).to eq([1, 1])
    end

    it 'compares again' do
        a = [17, 28, 30]
        b = [99, 16, 8]
        result = compareTriplets(a, b)
        expect(result).to eq([2, 1])
    end
end
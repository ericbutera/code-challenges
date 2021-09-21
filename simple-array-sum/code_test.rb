require './code'

RSpec.describe 'simple array sum' do
    it 'sums 6 when given 1 2 3' do
      output = array_sum_as_string([1,2,3])
      expect(output).to eq('1 + 2 + 3 = 6')
    end

    it 'sums 31 when given 1 2 3 4 10 11' do
      output = array_sum_as_string([1, 2, 3, 4, 10, 11])
      expect(output).to eq('1 + 2 + 3 + 4 + 10 + 11 = 31')
    end

    it 'handles calls to simpleArraySum' do
      output = simpleArraySum([1,2])
      expect(output).to eq('1 + 2 = 3')
    end
end

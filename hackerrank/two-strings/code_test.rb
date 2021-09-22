require './code'

RSpec.describe 'two strings' do
    it 'compares hello world' do
        expect(twoStrings('hello', 'world')).to eq('YES')
    end

    it 'compares hi world' do
        expect(twoStrings('hi', 'world')).to eq('NO')
    end

    it 'compares wouldyoulikefrie abcabcabcabcabcabc' do
        expect(twoStrings('wouldyoulikefrie', 'abcabcabcabcabcabc')).to eq('NO')
    end

    it 'compares hackerrankcommunity cdecdecdecde' do
        expect(twoStrings('hackerrankcommunity','cdecdecdecde')).to eq('YES')
    end

    it 'compares jackandjill wentupthehill' do
        expect(twoStrings('jackandjill', 'wentupthehill')).to eq('YES')
    end

    it 'compares writetoyourparents fghmqzldbc' do
        expect(twoStrings('writetoyourparents', 'fghmqzldbc')).to eq('NO')
    end

    it 'compares aardvark apple' do
        expect(twoStrings('aardvark', 'apple')).to eq('YES')
    end

    it 'compares beetroot sandals' do
        expect(twoStrings('beetroot', 'sandals')).to eq('NO')
    end
end

RSpec.describe 'sock merchant' do
  def sockMerchant(n, ar)
      sock_merchant(n, ar)
  end

  def sock_merchant(sock_count, socks)
    matched = match_socks(socks)
    total_pairs = 0

    matched.each do |sock, count|
      pair = count / 2
      total_pairs += pair
    end

    total_pairs
  end

  def match_socks(socks)
    matched = {}

    socks.each do |sock|
      if matched.has_key?(sock)
        matched[sock] += 1
      else 
        matched[sock] = 1
      end
    end

    matched
  end

  it 'matches example' do
    socks = [1, 2, 1, 2, 1, 3, 2]
    res = sockMerchant(7, socks)
    expect(res).to eq(2)
  end

  it 'matches example 1' do
    socks = [10, 20, 20, 10, 10, 30, 50, 10, 20]
    res = sockMerchant(socks.size, socks)
    expect(res).to eq(3)
  end
end
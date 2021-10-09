# https://www.hackerrank.com/challenges/drawing-book/problem

RSpec.describe 'drawing book' do

  ##
  # @param <int> n number of pages in the book 
  # @param <int> p page number to turn to 
  # @return <int> minimum number of pages to turn
  def pageCount(n, p)
    find_page(n, p)
  end

  def find_page(pages, page)
    x = 0
    flipped = 0
    max = pages.even? ? pages + 1 : pages

    while x <= pages do
      # pairs: (0,1)(5,4) (2,3)(3,2)
      # front = [x, x + 1]              # (0,1)
      # back  = [max - x, max - x - 1]  # (5, 4)
      # return flipped if front.include?(page) 
      # return flipped if back.include?(page)

      return flipped if page == x \
        || page == x + 1 \
        || page == max - x \
        || page == max - x - 1

      x += 2 
      flipped += 1
    end
  end

  it 'solves sample 0' do
    res = find_page(5, 3)
    expect(res).to eq(1)
  end
  
  it 'solves example 00' do
    res = find_page(6, 2)
    expect(res).to eq(1)
  end

  it 'solves example 01' do
    res = find_page(5, 4)
    expect(res).to eq(0)
  end

  it 'solves example 26' do
    res = find_page(6, 5)
    expect(res).to eq(1)
  end

  # size 6
  # x1 23 45 6x

  # size 7
  # x1 23 45 67

  # size 8
  # x1 23 45 67 8x

end
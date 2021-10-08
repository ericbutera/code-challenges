# https://www.hackerrank.com/challenges/bon-appetit/problem

def bonAppetit(bill, k, b)
  puts bon_appetit(bill, k, b)
end

##
# bonAppetit has the following parameter(s):
# @param bill an array of integers representing the cost of each item ordered
# @param k an integer representing the zero-based index of the item Anna doesn't eat
# @param b the amount of money that Anna contributed to the bill
#
# return Bon Appetit if the bill is fairly split. Otherwise, return integer amount of money that Brian owes Anna. 
def bon_appetit(bill, item_skipped, amount_paid)
  total = bill.sum
  shared_item_total = total - bill[item_skipped]
  cost_per_person = shared_item_total / 2

  if amount_paid == cost_per_person
    return 'Bon Appetit'
  elsif
    refund = amount_paid - cost_per_person
    return refund
  end
end


# The first line contains two space-separated integers n and k, the number of items ordered and the 0-based index of the item that Anna did not eat. 
# The second line contains n space-separated integers bill[i] where 0 < i < n.
# The third line contains an integer b, the amount of money Brian charged Anna for her share of the bill.

# If Brian did not overcharge Anna, print Bon Appetit on a new line; otherwise, print the difference (i.e., b charged - b actual) that Brian must refund to Anna. This will always be an integer. 


RSpec.describe 'bon appetit' do
  it 'sample input 0' do
    # 4 1         # <-- anna didnt eat item at index 1
    # 3 10 2 9    #     bill[1] = $10
    # 12          # <-- anna contributed $12
    k = 1
    b = 12
    bill = [3, 10, 2, 9]
    res = bonAppetit(bill, k, b)
    expect(res).to eq(5)
  end

  it 'sample input 1' do
    k = 1
    b = 7
    bill = [3, 10, 2, 9]
    res = bonAppetit(bill, k, b)
    expect(res).to eq('Bon Appetit')
  end

  it 'input00' do
    # 4 1
    # 3 10 2 9
    # 12
    #
    # 5
    k = 1
    b = 12
    bill = [3, 10, 2, 9]
    res = bonAppetit(bill, k, b)
    expect(res).to eq(5)
  end

  it 'input06' do
    # 4 1
    # 3 10 2 9
    # 7
    #
    # Bon Appetit
    k = 1
    b = 7
    bill = [3, 10, 2, 9]
    res = bonAppetit(bill, k, b)
    expect(res).to eq('Bon Appetit')
  end
end
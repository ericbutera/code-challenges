
RSpec.describe 'cats and a mouse' do
  it 'solves example' do
    expect(catAndMouse(2, 5, 4)).to eq('Cat B')
  end

  it 'solves sample input 0' do
    expect(catAndMouse(1, 2, 3)).to eq('Cat B')
    expect(catAndMouse(1, 3, 2)).to eq('Mouse C')
  end

  def catAndMouse(cat_a, cat_b, mouse)
    cat_a_distance = (cat_a - mouse).abs
    cat_b_distance = (cat_b - mouse).abs

    if cat_a_distance < cat_b_distance
      'Cat A'
    elsif cat_b_distance < cat_a_distance
      'Cat B'
    else
      'Mouse C'
    end
  end
end
# frozen_string_literal: true

require "test_helper"

class GameTest < Minitest::Test
  def test_default_equal
    stats1 = Stats.new
    stats2 = Stats.new
    assert_equal stats1, stats2
  end

  def test_init
    stats = Stats.new(abort_count: 1, explode_count: 2)
    assert_equal(2, stats.explode_count)
    assert_equal(1, stats.abort_count)
  end

  def test_init_values
    stats = Stats.new
    assert_equal 0, stats.abort_count
    assert_equal 0, stats.explode_count
    assert_equal 0, stats.retry_count
    assert_equal 0, stats.distance_traveled_km
    assert_equal 0, stats.fuel_burned_liters
    assert_equal 0, stats.flight_time_seconds
  end
end

# frozen_string_literal: true

require "test_helper"

class GameTest < Minitest::Test
  def test_game
    stats = Stats.new
    stats.distance_traveled_km = 1
    stats.abort_count = 1
    stats.explode_count = 1
    stats.fuel_burned_liters = 1
    stats.flight_time_seconds = 1

    game = Game.new
    game.record(stats)

    assert_equal 1, game.stats.distance_traveled_km
    assert_equal 1, game.stats.abort_count
    assert_equal 1, game.stats.explode_count
    assert_equal 1, game.stats.fuel_burned_liters
    assert_equal 1, game.stats.flight_time_seconds
  end
end

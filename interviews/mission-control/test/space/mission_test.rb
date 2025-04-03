# frozen_string_literal: true

require "test_helper"

class MissionTest < Minitest::Test
  def test_distance_traveled
    mission = mission_test
    # mission.start(Time.now - 60)
    # mission.finish
    mission.ticks(60)
    assert_in_delta 25, mission.distance_traveled
  end

  def test_abort_afterburner
    mission = mission_test

    mission.afterburner
    assert_equal :afterburner, mission.state_sym

    mission.abort_launch
    assert_equal :abort_launch, mission.state_sym

    mission.pending
    assert_equal :pending, mission.state_sym
  end

  def test_explode
    mission = mission_test
    mission.pending
    mission.afterburner
    mission.explode
    assert_equal :explode, mission.state_sym
  end

  def test_explode_stage
    mission = Mission.new(random_abort: false)
    mission.pending
    mission.explode_stage = :afterburner
    mission.afterburner
    assert_equal :explode, mission.state_sym
  end

  def test_abort_retry_count # rubocop:disable Metrics/AbcSize
    mission = mission_test
    mission.pending
    assert_equal 0, mission.stats.abort_count
    assert_equal 0, mission.stats.retry_count

    mission.abort_launch
    assert_equal 1, mission.stats.abort_count
    assert_equal 0, mission.stats.retry_count

    mission.pending
    assert_equal 1, mission.stats.abort_count
    assert_equal 1, mission.stats.retry_count
  end

  def test_burn
    # NOTE: 6 min burn is 1_009_440 liters of fuel
    mission = mission_test
    mission.tick(60 * 6)
    stats = mission.calculate_stats
    assert_equal(1_009_440, stats.fuel_burned_liters)
  end
end

def mission_test
  Mission.new(Mission::DEFAULT_SEED, random_abort: false, random_explode: false)
end

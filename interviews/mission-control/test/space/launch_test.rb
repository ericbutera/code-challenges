# frozen_string_literal: true

require "test_helper"

class LaunchTest < Minitest::Test
  def test_invalid_transition
    launch = Launch.new
    assert_equal false, launch.launch
    assert_equal :pending, launch.state.to_sym
  end

  def test_default_state
    launch = Launch.new
    assert_equal :pending, launch.state.to_sym
  end

  def test_complete
    Launch::COMPLETION_STATES.each do |state|
      launch = Launch.new
      assert_equal false, launch.complete?
      launch.state = state
      assert_equal true, launch.complete?
    end
  end

  def test_launch_sequence # rubocop:disable Metrics/AbcSize
    launch = Launch.new

    assert_equal true, launch.afterburner
    assert_equal :afterburner, launch.state.to_sym

    launch.disengage_release_structure
    assert_equal :disengage_release_structure, launch.state.to_sym

    launch.cross_checks
    assert_equal :cross_checks, launch.state.to_sym

    launch.launch
    assert_equal :launch, launch.state.to_sym
  end

  def test_afterburner_abort
    launch = Launch.new
    launch.afterburner
    launch.abort_launch
    assert_equal :pending, launch.state.to_sym
  end

  # TODO: test explosion states
end

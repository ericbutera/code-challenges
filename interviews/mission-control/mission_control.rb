# frozen_string_literal: true

require_relative "autoload"
require "cli/ui"
require_relative "config/logger"
require "active_support/all"
require "action_view" # Ensure ActionView is loaded.
require "action_view/helpers" # Ensure action view helpers are loaded.

class MissionControl
  include ActiveSupport::NumberHelper
  include ActionView::Helpers::DateHelper

  def initialize
    @game = Game.new
    @counter = 1
  end

  def run
    CLI::UI::StdoutRouter.enable

    puts "Welcome to Mission Control!"

    loop do
      main_menu
      @counter += 1
    end
  end

  private

  def main_menu
    CLI::UI::Prompt.ask("Select an option:") do |handler|
      handler.option("New mission") do
        CLI::UI::Frame.open("Mission #{@counter}") do
          mission
        end
      end
      handler.option("Flight statistics") do
        display_flight_stats
      end
      handler.option("Exit") do
        display_flight_stats
        exit
      end
    end
  end

  def display_flight_stats
    display_stats "Cumulative flight summary", @game.stats
  end

  def mission
    mission = Mission.new(@counter)

    display_plan mission
    display_launch_sequence mission
    display_flight mission

    stats = mission.calculate_stats # TODO: move into mission.stats
    @game.record stats # TODO: this should be a hidden concern inside game
    display_stats "Mission summary", stats
  end

  # ask for confirmation
  def prompt?(message)
    CLI::UI::Prompt.confirm(message)
  end

  # acknowledge will proceed to next step (cannot be cancelled)
  def acknowledge?(message)
    CLI::UI::Prompt.confirm(message)
    true
  end

  def launch_sequence(mission) # rubocop:disable Metrics/AbcSize,Metrics/MethodLength
    puts "What is the name of this mission? <Minerva>" # TODO: capture name
    return unless prompt? "Would you like to proceed?" # puts "Would you like to proceed? (Y/n) Y"

    puts "Mission name: Minerva" # TODO: capture name

    loop do
      next unless prompt? "Engage afterburner?" # puts "Engage afterburner? (Y/n) Y"

      # puts "Launch aborted! Would you like to retry? (Y/n) Y" if mission.state_sym == :abort_launch
      if mission.state_sym == :abort_launch
        res = prompt?("Retry launch?")
        return unless res
      end

      mission.pending
      break if mission.afterburner
      # TODO: ensure launch exits if afterburner fails (this block isn't like the other steps)
    end

    puts "Afterburner engaged!"

    # Any step past this point can result in an explosion

    # puts "Release support structures? (Y/n) Y"
    acknowledge? "Release support structures?" # puts "Release support structures?"
    return unless mission.disengage_release_structure

    puts "Support structures released!"

    acknowledge? "Perform cross-checks?" # puts "Perform cross-checks? (Y/n) Y"
    return unless mission.cross_checks

    puts "Cross-checks performed!"

    acknowledge? "Launch?" # puts "Launch? (Y/n) Y"
    return unless mission.launch

    puts "Launched!"

    true
  end

  def display_launch_sequence(mission)
    CLI::UI::Frame.open("\u{1F680} Control") do
      if launch_sequence mission
        puts "Mission launch sequence complete!"
      else
        # TODO: make this error friendly
        puts "ERROR: launch issue! State: `#{mission.state_sym}`"
      end
    end
  end

  def display_flight(mission)
    return unless mission.state_sym == :launch

    res = true # TODO: tidy flow control
    CLI::UI::Frame.open("\u{1F680} Launch") do
      CLI::UI::Spinner.spin("Mission Status: Loading...") do |spinner|
        until mission.completed?
          success = mission.tick
          unless success
            res = false
            break
          end

          spinner.update_title(mission_status(mission))
          sleep(0.005)
        end

        mission.complete if res
      end
    end
  end

  def mission_status(mission)
    # TODO: fuel remaining (percentage)
    # TODO: estimated time remaining
    [
      "Distance: #{number_to_rounded(mission.distance_traveled_km,
                                     significant: true)}/#{Mission::TARGET_DISTANCE_KM} km",
      "Fuel: #{mission.fuel_burned_liters}/#{Mission::FUEL_LITERS} liters",
      "Elapsed: #{seconds_to_words(mission.duration_seconds)}"
    ].join(" | ")
  end

  def display_plan(mission)
    # Travel distance:  160.0 km
    # Payload capacity: 50,000 kg
    # Fuel capacity:    1,514,100 liters
    # Burn rate:        168,240 liters/min
    # Average speed:    1,500 km/h
    # Random seed:      12
    CLI::UI::Frame.open("\u{1F680} Mission plan") do
      puts "Travel distance:  #{Mission::TARGET_DISTANCE_KM} km"
      puts "Payload capacity: #{number_to_delimited(Mission::PAYLOAD_CAPACITY_KG)} kg"
      puts "Fuel capacity:    #{number_to_delimited(Mission::FUEL_LITERS)} liters"
      puts "Burn rate:        #{number_to_delimited(Mission::BURN_RATE_LITERS_PER_MINUTE)} liters/min"
      puts "Average speed:    #{number_to_delimited(Mission::SPEED_KPH)} km/h"
      puts "Random seed:      #{mission.seed}"
    end
  end

  def display_stats(header, stats)
    # 1. Total distance traveled (for all missions combined)
    # 2. Number of abort and retries (for all missions combined)
    # 3. Number of explosions (for all missions combined)
    # 4. Total fuel burned (for all missions combined)
    # 5. Total flight time (for all missions combined)
    CLI::UI::Frame.open("\u{1F4CA} #{header}") do
      puts "Total distance traveled: #{number_to_delimited(stats.distance_traveled_km)} km"
      puts "Number of abort and retries: #{stats.abort_count}/#{stats.retry_count}"
      puts "Number of explosions: #{stats.explode_count}"
      puts "Total fuel burned: #{number_to_delimited(stats.fuel_burned_liters)} liters"
      puts "Flight time: #{seconds_to_words(stats.flight_time_seconds)}"
      puts "Total missions: #{stats.mission_count}" if stats.mission_count > 1
    end
  end

  def seconds_to_words(seconds)
    distance_of_time_in_words(seconds) # TODO: include_seconds
  end
end

# frozen_string_literal: true

require "semantic_logger"

# TODO: mission failed out_of_fuel state

# Mission tracks individual launches
class Mission
  include SemanticLogger::Loggable

  # TODO: revisit if these need to be externally accessible
  attr_accessor :ticks, :explode_stage
  attr_reader :seed, :stats

  DEFAULT_SEED = 12
  BURN_RATE_LITERS_PER_MINUTE = 168_240
  FUEL_LITERS = 1_514_100
  PAYLOAD_CAPACITY_KG = 50_000
  # Speed in kilometers per hour
  SPEED_KPH = 1_500
  # Flight path distance in kilometers
  TARGET_DISTANCE_KM = 160

  def initialize(seed = nil, random_abort: true, random_explode: true)
    # @name = name
    @state = Launch.new
    @stats = Stats.new
    @ticks = 0
    @explode_tick = nil

    @random_abort = random_abort
    if random_explode
      @random_explode = random_explode
      @explode_stage = generate_explode_stage
    end

    @seed = seed || srand # TODO: use for deterministic 'randomness' :)
    srand(@seed)

    logger.debug("initialize",
                 seed: @seed,
                 explode_stage: @explode_stage,
                 random_abort: @random_abort, random_explode: @random_explode)
  end

  def state_sym
    @state.state.to_sym
  end

  def tick(seconds = 1)
    if @ticks == @explode_tick
      explode
      return false
    end

    @ticks += seconds
    calculate_stats
    logger.trace("tick", seconds: seconds, ticks: @ticks, explode_tick: @explode_tick)

    true
  end

  def completed?
    state_complete = @state.complete?
    destination_reached = distance_traveled_km >= TARGET_DISTANCE_KM
    out_of_fuel = fuel_burned_liters >= FUEL_LITERS
    logger.trace("completed?",
                 state: state_sym,
                 state_complete: state_complete,
                 destination_reached: destination_reached,
                 out_of_fuel: out_of_fuel)

    state_complete || destination_reached || out_of_fuel
  end

  def calculate_stats
    @stats.flight_time_seconds = duration_seconds
    @stats.distance_traveled_km = distance_traveled_km
    @stats.fuel_burned_liters = fuel_burned_liters

    logger.trace("calculate_stats", ticks: @ticks, stats: @stats)

    @stats
  end

  # Flight time in seconds
  def duration_seconds
    @ticks
  end

  def speed_km_per_minute
    # 1,500 km per hour / 60 mins = 25 km per minute
    SPEED_KPH / 60.0
  end

  def speed_km_per_second
    speed_km_per_minute / 60.0 # km per second
  end

  def distance_traveled_km
    # 1,500 km per hour / 60 mins = 25 km per minute
    # 25 km per minute / 60 seconds = 0.41 km per second
    speed_km_per_second * duration_seconds
  end

  def fuel_burned_liters
    burn_rate_seconds = BURN_RATE_LITERS_PER_MINUTE / 60.0
    burn_rate_seconds * duration_seconds
  end

  # TODO:
  # estimated time of arrival in seconds
  # def eta_seconds
  #   remaining_time_minutes = distance_remaining_km / speed_km_per_minute
  #   remaining_time_minutes * 60
  # end
  # def distance_remaining_km
  #   remain = TARGET_DISTANCE_KM - distance_traveled_km
  #   remain.negative? ? 0 : remain
  # end
  # def fuel_remaining_liters
  #   FUEL_LITERS - fuel_burned_liters
  # end

  # One in every 5th launch will explode
  def generate_explode_stage
    must_explode = rand <= 0.20 # 1 in 5
    return unless must_explode

    Launch::EXPLODABLE_STATES.sample # TODO: weight to launch stage
  end

  def explode_stage?(stage)
    return false unless @random_explode

    match = stage == @explode_stage
    logger.debug("explode_stage?",
                 stage: stage,
                 explode_stage: @explode_stage,
                 match: match)
    explode if match
  end

  # One in every 3rd launch will require an abort and retry after stage 1
  def random_abort?
    return false unless @random_abort
    return unless rand <= 0.33 # 1 in 3

    abort_launch
    false
  end

  def pending
    @stats.retry_count += 1 if state_sym == :abort_launch
    @state.pending
  end

  def afterburner
    return false if explode_stage? :afterburner
    return false if random_abort?

    @state.afterburner
  end

  def disengage_release_structure
    return false if explode_stage? :disengage_release_structure

    @state.disengage_release_structure
  end

  def cross_checks
    return false if explode_stage? :cross_checks

    @state.cross_checks
  end

  def launch
    # NOTE: explosion check is inside tick
    @explode_tick = rand(1..360) if explode_stage == :launch # precompute explosion time

    @state.launch
  end

  def abort_launch
    if state_sym != :abort_launch
      @stats.abort_count += 1
      @ticks = 0
      @state.abort_launch
    else
      true
    end
  end

  def explode
    @stats.explode_count = 1
    @state.explode
  end

  def complete
    @state.complete
  end
end

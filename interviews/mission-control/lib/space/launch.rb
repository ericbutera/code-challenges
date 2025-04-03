# frozen_string_literal: true

require "state_machines"
require "semantic_logger"

# The rocket launch system is comprised of 4 stages, which must happen in this precise order:
# 1. Manually transition between launch stages in the expected order
# 2. Mission control should be able to safely abort launch after stage 1 and retry
# 3. One in every 3rd launch will require an abort and retry after stage 1, randomize when it actually happens
# 4. One in every 5th launch will explode, randomize when it actually happens
#
# Active mission controls:
# 1. Enable stage 1 afterburner
# 2. Disengaging release structure
# 3. Cross-checks
# 4. Launch

# Launch defines rules to ensure launch sequences happen in allowed order.
class Launch
  include SemanticLogger::Loggable

  ABORTABLE_STATES = %i[pending afterburner].freeze
  EXPLODABLE_STATES = %i[afterburner disengage_release_structure cross_checks launch].freeze
  COMPLETION_STATES = %i[complete explode].freeze

  logger = SemanticLogger["Launch"]

  def complete?
    Launch::COMPLETION_STATES.include?(state.to_sym)
  end

  state_machine :state, initial: :pending do # rubocop:disable Metrics/BlockLength
    event :pending do
      transition %i[afterburner abort_launch] => :pending, pending: same
    end

    event :afterburner do
      transition pending: :afterburner
    end

    event :disengage_release_structure do
      transition afterburner: :disengage_release_structure
    end

    event :cross_checks do
      transition disengage_release_structure: :cross_checks
    end

    event :launch do
      transition cross_checks: :launch
    end

    event :explode do
      transition EXPLODABLE_STATES => :explode
    end

    event :complete do
      transition launch: :complete
    end

    event :abort_launch do
      transition ABORTABLE_STATES => :abort_launch
    end

    before_transition do |mission, transition|
      logger.debug("before_transition",
                   mission: mission.state,
                   transition: transition.to_name)
    end
  end
end

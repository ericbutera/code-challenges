# frozen_string_literal: true

require "semantic_logger"

# Game engine
class Game
  include SemanticLogger::Loggable

  attr_accessor :abort_count, :explode_count
  attr_reader :stats

  def initialize
    @stats = Stats.new
  end

  def record(stats)
    logger.debug("merging stats", stats: stats)
    stats.mission_count = 1
    @stats.merge stats
  end
end

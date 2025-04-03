# frozen_string_literal: true

require_relative "config/logger"
logger = SemanticLogger["app"]
logger.info("Starting application")

require_relative "mission_control"
game = MissionControl.new
game.run

logger.info("Exit")
